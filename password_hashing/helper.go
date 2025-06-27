package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"

	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/scrypt"
)

func errExit(msg string, err error) {
	fmt.Printf(msg+" : %v", err)
	os.Exit(1)
}

func sha256Hasher(password string) string {

	hash := sha256.Sum256([]byte(password))

	return hex.EncodeToString(hash[:])
}

func hmac256Hasher(password string, secret []byte) string {
	hash := hmac.New(sha256.New, secret)

	hash.Write([]byte(password))

	return hex.EncodeToString(hash.Sum(nil))
}

func pbkdf2Hasher(password string, secret []byte, data map[string]interface{}) string {
	iter := data["rounds"].(float64)

	hash := pbkdf2.Key([]byte(password), secret, int(iter), sha256.Size, sha256.New)

	return hex.EncodeToString(hash)
}

func scryptHasher(password string, secret []byte, data map[string]interface{}) string {
	N := data["N"].(float64)
	r := data["r"].(float64)
	p := data["p"].(float64)
	keyLen := data["buflen"].(float64)
	hash, _ := scrypt.Key([]byte(password), secret, int(N), int(r), int(p), int(keyLen))

	return hex.EncodeToString(hash)
}
