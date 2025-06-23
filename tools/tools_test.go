package tools

import (
	"fmt"
	"os"
	"testing"
)

func TestFetchProblem(t *testing.T) {
	t.Run("fetching problem from hackattic", func(t *testing.T) {
		t.Helper()
		token := os.Getenv("HACKATTIC_TOKEN")
		got, err := FetchProblem("help_me_unpack", token)
		_, ok := got["bytes"]
		if err != nil || !ok {
			t.Errorf("got unexpected result : %v", err)
		}

		fmt.Println(got)
	})
}
