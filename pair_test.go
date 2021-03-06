package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
)

func TestNamesForUsernames(t *testing.T) {
	names, err := namesForUsernames([]string{}, map[string]string{})
	if names != "" {
		t.Fatalf("expected empty string for empty list of usernames, got %s", names)
	}
	if err != nil {
		t.Fatalf("expected no error for empty list of usernames, got %v", err)
	}

	names, err = namesForUsernames([]string{"mb"}, map[string]string{"mb": "Michael Bluth"})
	if names != "Michael Bluth" {
		t.Fatalf("expected 'Michael Bluth' for single username 'mb', got %s", names)
	}
	if err != nil {
		t.Fatalf("expected no error for single existing username, got %v", err)
	}

	names, err = namesForUsernames([]string{"lb", "mb"}, map[string]string{"mb": "Michael Bluth", "lb": "Lindsay Bluth"})
	if names != "Lindsay Bluth and Michael Bluth" {
		t.Fatalf("expected 'Lindsay Bluth and Michael Bluth', got %s", names)
	}
	if err != nil {
		t.Fatalf("expected no error for two existing usernames, got %v", err)
	}

	names, err = namesForUsernames([]string{"lb"}, map[string]string{"mb": "Michael Bluth"})
	if err == nil {
		t.Fatalf("expected error for a missing username, got nil")
	}
}

func ExampleEmailAddressForUsernames() {
	email, _ := emailAddressForUsernames("git@example.com", []string{})
	fmt.Println(email)
	email, _ = emailAddressForUsernames("git@example.com", []string{"mb"})
	fmt.Println(email)
	email, _ = emailAddressForUsernames("git@example.com", []string{"lb", "mb"})
	fmt.Println(email)

	// Output:
	// git@example.com
	// mb@example.com
	// git+lb+mb@example.com
}

func TestReadAuthorsByUsername(t *testing.T) {
	authorMap, err := readAuthorsByUsername(strings.NewReader(""))
	if len(authorMap) != 0 {
		t.Fatalf("expected reading an empty file to get zero authors, got %d", len(authorMap))
	}
	if err != nil {
		t.Fatalf("expected no error for empty authors file, got %v", err)
	}

	authorMap, err = readAuthorsByUsername(strings.NewReader("---\nmb: Michael Bluth"))
	if len(authorMap) != 1 || authorMap["mb"] != "Michael Bluth" {
		t.Fatalf("expected reading a single author as YAML to return one entry, got %v", authorMap)
	}
	if err != nil {
		t.Fatalf("expected reading a single author as YAML to have no errors, got %v", err)
	}

	authorMap, err = readAuthorsByUsername(strings.NewReader("---\nlb: Lindsay Bluth\nmb: Michael Bluth"))
	if len(authorMap) != 2 {
		t.Fatalf("expected reading multiple authors as YAML to return multiple entries, got %v", authorMap)
	}
	if err != nil {
		t.Fatalf("expected reading multiple authors as YAML to have no errors, got %v", err)
	}
}

func ExampleSplitEmail() {
	user, host, err := SplitEmail("a@b.com")
	fmt.Printf("error=%v user=%s host=%s\n", err, user, host)

	user, host, err = SplitEmail("")
	fmt.Printf("error=%v user=%s host=%s\n", err, user, host)

	user, host, err = SplitEmail("a@b@c")
	fmt.Printf("error=%v user=%s host=%s\n", err, user, host)

	// Output:
	// error=<nil> user=a host=b.com
	// error=invalid email address:  user= host=
	// error=invalid email address: a@b@c user= host=
}

func TestGitConfig(t *testing.T) {
	tempGitConfigFile, err := ioutil.TempFile(os.TempDir(), "pair-git-config")
	if err != nil {
		t.Fatal("unable to create temporary git config")
	}
	tempGitConfigPath := tempGitConfigFile.Name()

	err = setGitConfig(tempGitConfigPath, "user.name", "Michael Bluth")
	if err != nil {
		t.Fatalf("expected no error when setting git config, got %v", err)
	}

	value, err := gitConfig(tempGitConfigPath, "user.name")
	if err != nil {
		t.Fatalf("expected no error when getting git config, got %v", err)
	}
	if value != "Michael Bluth" {
		t.Fatalf("expected getting previously-set `user.name` to have the correct value, got %s", value)
	}
}

func ExamplePrintCurrentPairedUsers() {
	tempGitConfigFile, err := ioutil.TempFile(os.TempDir(), "pair-git-config")
	if err != nil {
		log.Fatal("unable to create temporary git config")
	}
	tempGitConfigPath := tempGitConfigFile.Name()

	err = setGitConfig(tempGitConfigPath, "user.name", "Michael Bluth")
	if err != nil {
		log.Fatalf("expected no error when setting git config, got %v", err)
	}

	err = setGitConfig(tempGitConfigPath, "user.email", "mb@example.com")
	if err != nil {
		log.Fatalf("expected no error when setting git config, got %v", err)
	}

	printCurrentPairedUsers(tempGitConfigPath)

	// Output:
	// Michael Bluth <mb@example.com>
}

func ExampleSetAndPrintNewPairedUsers() {
	tempPairsFile, err := ioutil.TempFile(os.TempDir(), "pair-pairs")
	if err != nil {
		log.Fatal("unable to create temporary pairs file")
	}
	io.WriteString(tempPairsFile, "---\nmb: Michael Bluth")
	tempPairsFile.Close()

	tempGitConfigFile, err := ioutil.TempFile(os.TempDir(), "pair-git-config")
	if err != nil {
		log.Fatal("unable to create temporary git config")
	}

	setAndPrintNewPairedUsers(tempPairsFile.Name(), tempGitConfigFile.Name(), "git@example.com", []string{"mb"})

	value, err := gitConfig(tempGitConfigFile.Name(), "user.name")
	if err != nil {
		log.Fatal("unable to get git config after setting users: %v", err)
	}
	fmt.Printf("user.name=%s\n", value)

	value, err = gitConfig(tempGitConfigFile.Name(), "user.email")
	if err != nil {
		log.Fatal("unable to get git config after setting users: %v", err)
	}
	fmt.Printf("user.email=%s\n", value)

	// Output:
	// Michael Bluth <mb@example.com>
	// user.name=Michael Bluth
	// user.email=mb@example.com
}
