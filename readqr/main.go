package main

import (
	"fmt"
	"image"
	"net/http"
	"os"

	"github.com/caiguanhao/readqr"
	"github.com/flames31/hack-attic/tools"
	"github.com/joho/godotenv"
)

const PROBLEM_NAME = "reading_qr"
const REQ_PARAM = "image_url"

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("error loading .env : %v\n", err)
		os.Exit(1)
	}

	accessToken := os.Getenv("ACCESS_TOKEN")

	problem, err := tools.FetchProblem(PROBLEM_NAME, accessToken)
	if err != nil {
		fmt.Printf("err is decoding qr img : %v\n", err)
		os.Exit(1)
	}

	imgURL := problem[REQ_PARAM].(string)

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

	resBody := map[string]interface{}{"code": decodedQr}

	finalResp, err := tools.SendReponse(resBody, accessToken, PROBLEM_NAME)
	if err != nil {
		fmt.Printf("err in sending response : %v\n", err)
		os.Exit(1)
	}

	fmt.Println(finalResp)
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
