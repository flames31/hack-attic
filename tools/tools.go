package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const BaseURL = "https://hackattic.com/challenges/"

func FetchProblem(problemName, accessToken string) (map[string]interface{}, error) {

	resp, err := http.Get(BaseURL + problemName + "/problem?access_token=" + accessToken)
	if err != nil {
		return nil, fmt.Errorf("err in response : %w", err)
	}

	var problem map[string]interface{}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&problem); err != nil {
		return nil, fmt.Errorf("err unmarshalling into problem: %w", err)
	}

	return problem, nil
}

func SendReponse(resBody map[string]interface{}, accessToken, problemName string) (string, error) {

	respData, err := json.Marshal(resBody)
	if err != nil {
		return "", fmt.Errorf("err in marshalling to json : %w", err)
	}

	finalResp, err := http.Post(BaseURL+problemName+"/solve?access_token="+accessToken, "application/json", bytes.NewBuffer(respData))
	if err != nil {
		return "", fmt.Errorf("err is sending response : %w", err)
	}

	defer finalResp.Body.Close()
	finalData, err := io.ReadAll(finalResp.Body)
	if err != nil {
		return "", fmt.Errorf("err in reading response body : %w", err)
	}

	return string(finalData), nil
}
