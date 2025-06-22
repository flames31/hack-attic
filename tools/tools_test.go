package tools

import (
	"fmt"
	"testing"
)

func TestFetchProblem(t *testing.T) {
	t.Run("fetching problem from hackattic", func(t *testing.T) {
		t.Helper()
		got, err := FetchProblem("help_me_unpack")
		_, ok := got["bytes"]
		if err != nil || !ok {
			t.Errorf("got unexpected result : %v", err)
		}

		fmt.Println(got)
	})
}
