package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func main() {
	head := os.Getenv("GITHUB_SHA")
	eventName := os.Getenv("GITHUB_EVENT_NAME")
	base := os.Getenv("GITHUB_BASE_REF")

	var cmd *exec.Cmd
	if eventName == "pull_request" {
		cmd = exec.Command("git", "log", "--pretty=format:%s", fmt.Sprintf("%s..%s", base, head))
	} else {
		cmd = exec.Command("git", "log", "--pretty=format:%s", "-1", "HEAD")
	}

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error running git log:", err)
		os.Exit(1)
	}

	commits := strings.Split(out.String(), "\n")
	pattern := regexp.MustCompile(`^(feat|fix|docs|style|refactor|test|chore|ci)(\([a-z0-9_-]+\))?: [a-z].+`)

	failed := false
	for _, c := range commits {
		c = strings.TrimSpace(c)
		if c == "" {
			continue
		}
		if strings.HasPrefix(c, "Merge pull request") {
			break // 不檢查 merge commit
		}
		if !pattern.MatchString(c) {
			fmt.Printf("❌ Invalid commit: %s\n", c)
			failed = true
		}
	}

	if failed {
		fmt.Println("❌ Commit lint failed. Please follow Conventional Commits format: type(scope): subject")
		os.Exit(1)
	}

	fmt.Println("✅ All commits passed commit lint.")
}
