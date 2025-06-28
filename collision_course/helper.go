package main

import (
	"crypto/md5"
	"fmt"
)

func hashFile(data []byte) []byte {
	hash := md5.Sum(data)
	return hash[:]
}

func errExit(msg string, err error) error {
	fmt.Printf(msg+": %v\n", err)
	return err
}
