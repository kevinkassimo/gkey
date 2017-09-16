package entry

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kevinkassimo/gokey/src/confirm"
	"golang.org/x/crypto/bcrypt"
	"io"
	"io/ioutil"
	"log"
)

type UserEntry struct {
	FileName string
	Name     string
	Password []byte
	Data     []byte
	Entry    DataEntry
}

func (e *UserEntry) ReadData() bool {
	data, err := ioutil.ReadFile(e.FileName)
	if err != nil {
		log.Fatalf("Cannot read file %s: %s\n", e.FileName, err)
	}
	e.Data = data

	if err := json.Unmarshal(data, &(e.Entry)); err != nil {
		log.Fatalf("User data unmarshalling failed: %s\n", err)
	}

	if err := bcrypt.CompareHashAndPassword(e.Entry.Hash, e.Password); err != nil {
		return false // Password validation failed
	}

	for i, entry := range e.Entry.Entries {
		entry.Decode(e.Password)
		// Must do this! entry is NOT a pointer, BUT a copy!
		e.Entry.Entries[i] = entry
	}

	return true
}

func (e *UserEntry) WriteData() {
	// We are assuming the Entry is having some data

	hash, err := bcrypt.GenerateFromPassword(e.Password, bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Bcrypt GenerateFromPassword failed: %s\n", err)
	}

	e.Entry.Hash = hash

	var newDataEntry DataEntry
	// make is necessary! When just trying to copy, it will NOT have enough capacity...
	newDataEntry.Hash = make([]byte, len(e.Entry.Hash))
	copy(newDataEntry.Hash, e.Entry.Hash)

	for _, ent := range e.Entry.Entries {
		var newEntry PasswordEntry
		newEntry.Name = make([]byte, len(ent.Name))
		copy(newEntry.Name, ent.Name)
		newEntry.Desc = make([]byte, len(ent.Desc))
		copy(newEntry.Desc, ent.Desc)
		newEntry.Password = make([]byte, len(ent.Password))
		copy(newEntry.Password, ent.Password)

		newEntry.Encode(e.Password)
		newDataEntry.Entries = append(newDataEntry.Entries, newEntry)
	}

	data, err := json.Marshal(newDataEntry)
	if err != nil {
		log.Fatalf("User data marshalling failed: %s\n", err)
	}

	if err := ioutil.WriteFile(e.FileName, data, 0770); err != nil {
		log.Fatalf("Write user data failed: %s\n", err)
	}
}

type DataEntry struct {
	Hash    []byte
	Entries []PasswordEntry `json: entries`
}

func (e *DataEntry) AddEntry(pe PasswordEntry) {
	findDup := false
	shouldOverwrite := false
	overwriteIndex := -1
	for i, entry := range e.Entries {
		if bytes.Equal(entry.Name, pe.Name) {
			findDup = true
			shouldOverwrite = confirm.Ask("Past record found, overwrite?")
			overwriteIndex = i
			break
		}
	}

	if !findDup {
		e.Entries = append(e.Entries, pe)
	} else {
		if shouldOverwrite {
			e.Entries[overwriteIndex] = pe
		}
	}
}

func (e *DataEntry) RemoveEntry(name []byte) {
	removeIndex := -1
	for i, entry := range e.Entries {
		if bytes.Equal(entry.Name, name) {
			removeIndex = i
			break
		}
	}

	if removeIndex < 0 {
		fmt.Printf("No entry of %s found, nothing to remove.\n", name)
	} else {
		// Non-in-seq removal
		if confirm.Ask("Are you sure?") {
			e.Entries[removeIndex] = e.Entries[len(e.Entries)-1]
			e.Entries = e.Entries[:len(e.Entries)-1]
		}
	}
}

type PasswordEntry struct {
	Name     []byte `json: name`
	Desc     []byte `json: description`
	Password []byte `json: password`
}

//Adapted from https://astaxie.gitbooks.io/build-web-application-with-golang/en/09.6.html

func (e *PasswordEntry) Decode(p []byte) {
	c, err := aes.NewCipher(p)
	if err != nil {
		log.Fatalf("Generate cipher error: %s\n", err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		log.Fatalf("Generate GCM error: %s\n", err)
	}
	nonceSize := gcm.NonceSize()

	if len(e.Name) < nonceSize {
		log.Fatalf("Error: Cipher text too short\n")
	}

	nonce0, cipherText0 := e.Name[:nonceSize], e.Name[nonceSize:]
	result0, err0 := gcm.Open(nil, nonce0, cipherText0, nil)
	if err0 != nil {
		log.Fatalf("Decoding failed: %s", err0)
	}
	e.Name = result0

	if len(e.Desc) < nonceSize {
		log.Fatalf("Error: Cipher text too short\n")
	}
	nonce1, cipherText1 := e.Desc[:nonceSize], e.Desc[nonceSize:]
	result1, err1 := gcm.Open(nil, nonce1, cipherText1, nil)
	if err1 != nil {
		log.Fatalf("Decoding failed: %s", err1)
	}
	e.Desc = result1

	if len(e.Password) < nonceSize {
		log.Fatalf("Error: Cipher text too short\n")
	}
	nonce2, cipherText2 := e.Password[:nonceSize], e.Password[nonceSize:]
	result2, err2 := gcm.Open(nil, nonce2, cipherText2, nil)
	if err2 != nil {
		log.Fatalf("Decoding failed: %s", err2)
	}
	e.Password = result2
}
func (e *PasswordEntry) Encode(p []byte) {
	c, err := aes.NewCipher(p)
	if err != nil {
		log.Fatalf("Generate cipher error: %s\n", err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		log.Fatalf("Generate GCM error: %s\n", err)
	}

	nonce := make([]byte, 12)

	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		log.Fatalf("Error: %s\n", err)
	}

	e.Name = gcm.Seal(nonce, nonce, e.Name, nil)
	e.Desc = gcm.Seal(nonce, nonce, e.Desc, nil)
	e.Password = gcm.Seal(nonce, nonce, e.Password, nil)
}

func decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonceSize := 12
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	return gcm.Open(nil, nonce, ciphertext, nil)
}
