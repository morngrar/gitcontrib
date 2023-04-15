package gitcontrib

import (
	"testing"
)

func Test_ExtractCheckedOutBranch(t *testing.T) {
	gbOutput := `  asdasd
* main
  bottombranch
`
	branch, err := extractCheckedOutBranch(gbOutput)
	if err != nil {
		t.Fatalf("encountered error extracting branch: %s", err)
	}

	if branch != "main" {
		t.Errorf("Expected 'main' branch, got %q", branch)
	}
}
