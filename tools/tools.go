package tools

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func FetchProblem(problemName string) (map[string]interface{}, error) {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("error loading .env : %v\n", err)
		os.Exit(1)
	}

	accessToken := os.Getenv("ACCESS_TOKEN")

	resp, err := http.Get("https://hackattic.com/challenges/" + problemName + "/problem?access_token=" + accessToken)
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
