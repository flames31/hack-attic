package main

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/flames31/hack-attic/tools"
)

const PROBLEM_NAME = "collision_course"

func main() {
	accessToken, err := tools.GetAccessToken()
	if err != nil {
		errExit("error getting token", err)
	}

	problem, err := tools.FetchProblem(PROBLEM_NAME, accessToken)
	if err != nil {
		errExit("error fetching problem", err)
	}

	includeString := problem["include"].(string)

	file1, _ := os.ReadFile("collision1.gif")
	file2, _ := os.ReadFile("collision2.gif")

	file1 = append(file1, []byte(includeString)...)
	file2 = append(file2, []byte(includeString)...)

	//hash1 := hashFile(file1)
	//hash2 := hashFile(file2)

	resBody := map[string]interface{}{
		"files": []string{base64.StdEncoding.EncodeToString(file1),
			base64.StdEncoding.EncodeToString(file2)},
	}

	response, err := tools.SendReponse(resBody, accessToken, PROBLEM_NAME)
	if err != nil {
		errExit("error sending response", err)
	}

	fmt.Println(response)
}
