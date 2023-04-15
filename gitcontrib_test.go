package gitcontrib

import (
	"io/ioutil"
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

func Test_MapAuthorCommits(t *testing.T) {
	gitOutput := `    42  Author One
     3  Author Two
`
	m, err := mapAuthorCommits(gitOutput)
	if err != nil {
		t.Fatalf("error mapping author commits: %s", err)
	}

	t.Logf("map: %v", m)

	if got := m["Author One"]; got != 42 {
		t.Errorf("Expected 42 commits for author one, got: %d", got)
	}
	if got := m["Author Two"]; got != 3 {
		t.Errorf("Expected 3 commits for author two, got: %d", got)
	}

}

func Test_MapLineChanges(t *testing.T) {
	knownSums := map[string]int{
		"Christopher Frantz":  3897,
		"Mariusz Nowostawski": 33,
		"siamak":              8,
		"Svein-Kåre Bjørnsen": 7,
		"Jon Gunnar Fossum":   2,
	}

	// read in test data file as string
	buf, err := ioutil.ReadFile("testdata/numstat-example")
	if err != nil {
		t.Fatalf("unable to read file: %s", err)
	}
	output := string(buf)

	authorMap, err := parseLineChanges(output)
	if err != nil {
		t.Fatalf("unable to parse output: %s", err)
	}

	for k, v := range authorMap {
		if knownSums[k] != v.Additions+v.Deletions {
			t.Errorf(
				"Author %q changes mismatch. Exp: %d, a: %d, d: %d",
				k, knownSums[k], v.Additions, v.Deletions,
			)
		}
	}
}
