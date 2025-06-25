package main

import (
	"fmt"

	"github.com/flames31/hack-attic/tools"
)

const PROBLEM_NAME = "mini_miner"

type Block struct {
	Data  []interface{} `json:"data"`
	Nonce int           `json:"nonce"`
}

var difficulty int

func main() {
	accessToken, err := tools.GetAccessToken()
	if err != nil {
		errorExit("could not fetch access token", err)
	}

	problem, err := tools.FetchProblem(PROBLEM_NAME, accessToken)
	if err != nil {
		errorExit("could not fetch problem", err)
	}

	difficulty = int(problem["difficulty"].(float64))
	if err != nil {
		errorExit("could not parse difficulty", err)
	}

	blockMap := problem["block"].(map[string]interface{})
	data := blockMap["data"].([]interface{})

	block := Block{
		Data: data,
	}

	resBody := map[string]interface{}{
		"nonce": findNonce(block),
	}

	response, err := tools.SendReponse(resBody, accessToken, PROBLEM_NAME)
	if err != nil {
		errorExit("error sending response", err)
	}

	fmt.Println(response)
}
