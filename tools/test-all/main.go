package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

var modules = []string{
	"backend/proto/gen/go",
	"backend/api-gateway",
	"backend/auth-service",
	"backend/user-service",
	"backend/ride-service",
	"backend/location-service",
	"backend/matching-service",
	"backend/notification-service",
	"backend/pricing-service",
	"backend/payment-service",
}

func main() {
	root, err := repositoryRoot()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to find repository root: %v\n", err)
		os.Exit(1)
	}

	for _, module := range modules {
		fmt.Printf("Testing %s\n", module)
		if err := run(filepath.Join(root, module), "go", "test", "./..."); err != nil {
			fmt.Fprintf(os.Stderr, "%s failed: %v\n", module, err)
			os.Exit(1)
		}
	}

	fmt.Println("All backend module tests passed.")
}

func repositoryRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if exists(filepath.Join(dir, "backend")) && exists(filepath.Join(dir, ".git")) {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("could not find backend and .git directories")
		}
		dir = parent
	}
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func run(dir string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
