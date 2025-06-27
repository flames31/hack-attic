package main

import (
	"encoding/base64"
	"fmt"

	"github.com/flames31/hack-attic/tools"
)

const PROBLEM_NAME = "password_hashing"

func main() {
	accessToken, err := tools.GetAccessToken()
	if err != nil {
		errExit("failed to get access token", err)
	}

	problem, err := tools.FetchProblem(PROBLEM_NAME, accessToken)
	if err != nil {
		errExit("failed to get prolem", err)
	}

	saltEncoded := problem["salt"].(string)
	password := problem["password"].(string)
	pbkdf2Data := problem["pbkdf2"].(map[string]interface{})
	scryptData := problem["scrypt"].(map[string]interface{})

	salt, err := base64.StdEncoding.DecodeString(saltEncoded)
	if err != nil {
		errExit("error decoding salt", err)
	}

	resBody := map[string]interface{}{
		"sha256": sha256Hasher(password),
		"hmac":   hmac256Hasher(password, salt),
		"pbkdf2": pbkdf2Hasher(password, salt, pbkdf2Data),
		"scrypt": scryptHasher(password, salt, scryptData),
	}

	response, err := tools.SendReponse(resBody, accessToken, PROBLEM_NAME)
	if err != nil {
		errExit("error with response :", err)
	}

	fmt.Println(response)

}
