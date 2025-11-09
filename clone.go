package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

func main() {
	user := flag.String("user", "", "GitHub username")
	dir := flag.String("dir", "", "Directory to clone repos into")
	workers := flag.Int("workers", 3, "Number of concurrent clone operations")
	flatClone := flag.Bool("flat", false, "Clone repos to parent folder (no owner subdirectory)")
	flag.Parse()

	githubUser := *user
	cloneDir := *dir

	if githubUser == "" {
		githubUser = promptUser("Enter your GitHub username: ")
	}
	if cloneDir == "" {
		cloneDir = promptUser("Enter the directory where you want to clone the repos: ")
	}

	cloneToParent := *flatClone
	if !cloneToParent {
		cloneToParent = promptYesNo(
			"Clone to parent folder without owner subdirectory? (y/n): ",
		)
	}

	if githubUser == "" || cloneDir == "" {
		fmt.Println("Error: GitHub username and clone directory are required")
		os.Exit(1)
	}

	if err := os.MkdirAll(cloneDir, 0755); err != nil {
		fmt.Printf("Failed to create directory: %v\n", err)
		os.Exit(1)
	}

	repos, err := fetchRepositories(githubUser)
	if err != nil {
		fmt.Printf("Failed to fetch repositories: %v\n", err)
		os.Exit(1)
	}

	if len(repos) == 0 {
		fmt.Println("No repositories found")
		return
	}

	fmt.Printf("Found %d repositories. Starting clones with %d workers...\n",
		len(repos), *workers)

	semaphore := make(chan struct{}, *workers)
	var wg sync.WaitGroup

	for _, repo := range repos {
		wg.Add(1)
		go func(r string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			cloneRepository(r, cloneDir, cloneToParent)
		}(repo)
	}

	wg.Wait()
	fmt.Printf("Completed! All repos have been cloned into %s\n", cloneDir)
}

func promptUser(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func promptYesNo(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	answer := strings.TrimSpace(strings.ToLower(input))
	return answer == "y" || answer == "yes"
}

func fetchRepositories(githubUser string) ([]string, error) {
	cmd := exec.Command(
		"gh", "repo", "list", githubUser,
		"--limit", "100",
		"--json", "nameWithOwner,isFork,isArchived",
		"--jq", ".[] | select(.isFork == false and .isArchived == false) | .nameWithOwner",
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("gh command failed: %w", err)
	}

	var repos []string
	for _, line := range strings.Split(strings.TrimSpace(string(output)), "\n") {
		if trimmed := strings.TrimSpace(line); trimmed != "" {
			repos = append(repos, trimmed)
		}
	}
	return repos, nil
}

func cloneRepository(repo, baseDir string, flatClone bool) {
	var clonePath string

	if flatClone {
		repoName := filepath.Base(repo)
		clonePath = filepath.Join(baseDir, repoName)
	} else {
		clonePath = filepath.Join(baseDir, repo)
	}

	if err := os.MkdirAll(clonePath, 0755); err != nil {
		fmt.Printf("❌ Failed to create directory for %s: %v\n", repo, err)
		return
	}

	cmd := exec.Command("git", "clone", fmt.Sprintf("git@github.com:%s", repo),
		clonePath)

	if output, err := cmd.CombinedOutput(); err != nil {
		fmt.Printf("❌ Failed to clone %s: %v\n", repo, err)
		if len(output) > 0 {
			fmt.Printf("   %s\n", string(output))
		}
	} else {
		fmt.Printf("✅ %s cloned successfully\n", repo)
	}
}