package main

import (
	"fmt"
	"os"

	"github.com/flames31/hack-attic/tools"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

const PROBLEM_NAME = "backup_restore"

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("error loading .env : %v", err)
		os.Exit(1)
	}

	accessToken := os.Getenv("ACCESS_TOKEN")
	dbUrl := os.Getenv("DB_URL")
	if accessToken == "" || dbUrl == "" {
		fmt.Printf("access token or db_url not set : %v", err)
		os.Exit(1)
	}

	problem, err := tools.FetchProblem(PROBLEM_NAME, accessToken)
	if err != nil {
		fmt.Printf("error fetching problem : %v", err)
		os.Exit(1)
	}

	dumpString := problem["dump"].(string)

	decompDump, err := decodeDumpData(dumpString)
	if err != nil {
		fmt.Printf("err decoding dump : %v", err)
		os.Exit(1)
	}

	err = restoreDump(decompDump, dbUrl)
	if err != nil {
		fmt.Printf("err restoring dump to db : %v", err)
		os.Exit(1)
	}

	ssnList, err := getSSNFromDB(dbUrl)
	if err != nil {
		fmt.Printf("error getting ssn from DB : %v", err)
		os.Exit(1)
	}

	resBody := map[string]interface{}{
		"alive_ssns": ssnList,
	}

	response, err := tools.SendReponse(resBody, accessToken, PROBLEM_NAME)
	if err != nil {
		fmt.Printf("error sending response : %v", err)
		os.Exit(1)
	}

	fmt.Println(response)
}
