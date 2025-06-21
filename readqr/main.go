package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"net/http"
	"os"

	"github.com/caiguanhao/readqr"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("error loading .env : %v\n", err)
		os.Exit(1)
	}

	accessToken := os.Getenv("ACCESS_TOKEN")
	resp, err := http.Get("https://hackattic.com/challenges/reading_qr/problem?access_token=" + accessToken)
	if err != nil {
		fmt.Printf("err is response : %v\n", err)
		os.Exit(1)
	}

	var imgResp struct {
		ImageURL string `json:"image_url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&imgResp); err != nil {
		fmt.Printf("err unmarshalling : %v\n", err)
		os.Exit(1)
	}

	qrResp, err := http.Get(imgResp.ImageURL)
	if err != nil {
		fmt.Printf("err fetching qr : %v\n", err)
		os.Exit(1)
	}

	defer qrResp.Body.Close()
	qrImg, _, err := image.Decode(qrResp.Body)
	if err != nil {
		fmt.Printf("err is decoding to img : %v\n", err)
		os.Exit(1)
	}
	decodedQr, err := readqr.DecodeImage(qrImg)
	if err != nil {
		fmt.Printf("err is decoding qr img : %v\n", err)
		os.Exit(1)
	}

	var ansResp struct {
		Code string `json:"code"`
	}

	ansResp.Code = decodedQr

	respData, err := json.Marshal(ansResp)
	if err != nil {
		fmt.Printf("err is decoding qr img : %v\n", err)
		os.Exit(1)
	}

	finalResp, err := http.Post("https://hackattic.com/challenges/reading_qr/solve?access_token="+accessToken, "application/json", bytes.NewBuffer(respData))

	if err != nil {
		fmt.Printf("err is decoding qr img : %v\n", err)
		os.Exit(1)
	}

	fmt.Println(finalResp.StatusCode)
}
