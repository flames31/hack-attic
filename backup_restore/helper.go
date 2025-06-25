package main

import (
	"bytes"
	"compress/gzip"
	"database/sql"
	"encoding/base64"
	"fmt"
	"io"
	"os/exec"
)

func restoreDump(dumpBytes []byte, dbURL string) error {
	cmd := exec.Command("psql", dbURL)
	cmd.Stdin = bytes.NewReader(dumpBytes)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("pg_restore failed: %w\nout: %v", err, string(out))
	}
	return nil
}

func decompressGzip(data []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("gzip reader error: %w", err)
	}
	defer reader.Close()

	var out bytes.Buffer
	if _, err := io.Copy(&out, reader); err != nil {
		return nil, fmt.Errorf("gzip decompression failed: %w", err)
	}

	return out.Bytes(), nil
}

func decodeDumpData(dumpString string) ([]byte, error) {
	dumpData, err := base64.StdEncoding.DecodeString(dumpString)
	if err != nil {
		return []byte{}, fmt.Errorf("error decoding base64 : %v", err)
	}

	decompDump, err := decompressGzip(dumpData)
	if err != nil {
		return []byte{}, fmt.Errorf("error decompressing dump bytes : %v", err)
	}

	return decompDump, nil
}

func getSSNFromDB(dbUrl string) ([]string, error) {
	db, err := sql.Open("postgres", dbUrl+"?sslmode=disable")
	if err != nil {
		return []string{}, fmt.Errorf("error opening DB: %v", err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT ssn FROM criminal_records WHERE status = 'alive'")
	if err != nil {
		return []string{}, fmt.Errorf("error querying DB: %v", err)
	}
	defer rows.Close()

	ssnList := make([]string, 0)
	for rows.Next() {
		var ssn string
		if err := rows.Scan(&ssn); err != nil {
			return []string{}, fmt.Errorf("error saving ssn from row: %v", err)
		}
		ssnList = append(ssnList, ssn)
	}

	return ssnList, nil
}
