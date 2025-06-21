package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"io"
	"net/http"
	"os"

	"github.com/caiguanhao/readqr"
	"github.com/joho/godotenv"
)

var imgResp struct {
	ImageURL string `json:"image_url"`
}

var ansResp struct {
	Code string `json:"code"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("error loading .env : %v\n", err)
		os.Exit(1)
	}

	accessToken := os.Getenv("ACCESS_TOKEN")

	imgURL, err := fetchImgUrl(accessToken)
	if err != nil {
		fmt.Printf("err in fetching qr : %v\n", err)
		os.Exit(1)
	}

	qrImg, err := decodeImgQR(imgURL)
	if err != nil {
		fmt.Printf("err is decoding qr img : %v\n", err)
		os.Exit(1)
	}

	decodedQr, err := readqr.DecodeImage(qrImg)
	if err != nil {
		fmt.Printf("err is decoding qr img : %v\n", err)
		os.Exit(1)
	}

	finalResp, err := sendResponse(decodedQr, accessToken)
	if err != nil {
		fmt.Printf("err in sending response : %v\n", err)
		os.Exit(1)
	}

	fmt.Println(finalResp)
}

func fetchImgUrl(accessToken string) (string, error) {
	resp, err := http.Get("https://hackattic.com/challenges/reading_qr/problem?access_token=" + accessToken)
	if err != nil {
		return "", fmt.Errorf("err is response : %w", err)
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&imgResp); err != nil {
		return "", fmt.Errorf("err unmarshalling : %w", err)
	}

	return imgResp.ImageURL, nil
}

func decodeImgQR(imgURL string) (image.Image, error) {
	qrResp, err := http.Get(imgURL)
	if err != nil {
		return nil, fmt.Errorf("err fetching img from url : %w", err)
	}

	defer qrResp.Body.Close()
	qrImg, _, err := image.Decode(qrResp.Body)
	if err != nil {
		return nil, fmt.Errorf("err decoding to img : %w", err)
	}

	return qrImg, nil
}

func sendResponse(decodedQR, accessToken string) (string, error) {
	ansResp.Code = decodedQR

	respData, err := json.Marshal(ansResp)
	if err != nil {
		return "", fmt.Errorf("err in marshalling to json : %w", err)
	}

	finalResp, err := http.Post("https://hackattic.com/challenges/reading_qr/solve?access_token="+accessToken, "application/json", bytes.NewBuffer(respData))
	if err != nil {
		return "", fmt.Errorf("err is decoding qr img : %w", err)
	}

	finalData, err := io.ReadAll(finalResp.Body)
	if err != nil {
		return "", fmt.Errorf("err in reading resp body : %w", err)
	}
	return string(finalData), nil
}
