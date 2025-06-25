package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
)

func findNonce(block Block) float64 {
	numWorkers := 10
	noncesPerWorker := 100000
	nonceSol := make(chan int)
	for i := 0; i < numWorkers; i++ {
		start := i * noncesPerWorker
		end := start + noncesPerWorker
		go worker(start, end, nonceSol, block)
	}

	answerNonce := <-nonceSol
	close(nonceSol)

	return float64(answerNonce)
}

func worker(start, end int, nonceSol chan int, block Block) {
	for nonce := start; nonce < end; nonce++ {
		select {
		case <-nonceSol:
			return
		default:
			if isValid(hashWithNonce(nonce, block)) {
				nonceSol <- nonce
				return
			}
		}
	}
}

func hashWithNonce(nonce int, block Block) [32]byte {
	block.Nonce = nonce
	bytes, _ := json.Marshal(block)
	return sha256.Sum256(bytes)
}

func isValid(hash [32]byte) bool {
	for i := 0; i < difficulty/8; i++ {
		if hash[i] != 0 {
			return false
		}
	}

	remBits := difficulty % 8
	if remBits > 0 {
		mask := byte(0xFF << (8 - remBits))
		return hash[difficulty/8]&mask == 0
	}
	return true
}

func errorExit(msg string, err error) {
	fmt.Printf(msg+": %v", err)
	os.Exit(1)
}
