package main

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"os"

	"github.com/flames31/hack-attic/tools"
	"github.com/joho/godotenv"
)

const PROBLEM_NAME = "help_me_unpack"

func main() {
	accessToken, err := getAccessToken()
	if err != nil {
		fmt.Printf("access token not set : %v", err)
		os.Exit(1)
	}
	problem, err := tools.FetchProblem(PROBLEM_NAME, accessToken)
	if err != nil {
		fmt.Printf("error fetching problem : %v", err)
		os.Exit(1)
	}

	bytes := problem["bytes"].(string)

	byteData, err := base64.StdEncoding.DecodeString(bytes)
	if err != nil {
		fmt.Printf("error decoding base64 : %v", err)
		os.Exit(1)
	}

	resBody := map[string]interface{}{
		"int":               int32(binary.LittleEndian.Uint32(byteData[:4])),
		"short":             int16(binary.LittleEndian.Uint16(byteData[8:10])),
		"double":            math.Float64frombits(binary.LittleEndian.Uint64(byteData[16:24])),
		"float":             math.Float32frombits(binary.LittleEndian.Uint32(byteData[12:16])),
		"uint":              binary.LittleEndian.Uint32(byteData[4:8]),
		"big_endian_double": math.Float64frombits(binary.BigEndian.Uint64((byteData[24:32]))),
	}
	response, err := tools.SendReponse(resBody, accessToken, PROBLEM_NAME)
	if err != nil {
		fmt.Printf("error decoding base64 : %v", err)
		os.Exit(1)
	}

	fmt.Println(response)

}

func getAccessToken() (string, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return "", fmt.Errorf("error loading .env : %w", err)
	}

	accessToken := os.Getenv("ACCESS_TOKEN")
	if accessToken == "" {
		return "", errors.New("no access token set")
	}

	return accessToken, nil

}
