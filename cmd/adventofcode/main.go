package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const (
	minYear = 2015
	maxYear = 2024
	minDay  = 1
	maxDay  = 25
)

func main() {
	// Validate command-line arguments
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing arguments.")
		fmt.Printf("Usage: %s <year> <day>\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	if len(os.Args) > 3 {
		fmt.Println("Error: Too many arguments provided.")
		fmt.Printf("Usage: %s <year> <day>\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	year, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Printf("Error: Invalid year argument: %v\n", os.Args[1])
		os.Exit(1)
	}

	if year < minYear || year > maxYear {
		fmt.Printf("Error: Year must be between %d and %d\n", minYear, maxYear)
		os.Exit(1)
	}

	day, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Printf("Error: Invalid day argument: %v\n", os.Args[2])
		os.Exit(1)
	}

	if day < minDay || day > maxDay {
		fmt.Printf("Error: Day must be between %d and %d\n", minDay, maxDay)
		os.Exit(1)
	}

	// Get session value from environment variable
	session := os.Getenv("AOC_SESSION")
	if session == "" {
		fmt.Fprintln(os.Stderr, "Error: AOC_SESSION not found in environment variables.")
		os.Exit(1)
	}

	// Create HTTP client and send request
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://adventofcode.com/%d/day/%d/input", year, day), nil)
	if err != nil {
		fmt.Printf("Error: Failed to create HTTP request: %s", err)
		os.Exit(1)
	}

	req.AddCookie(&http.Cookie{Name: "session", Value: session})

	resp, err := client.Do(req)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			fmt.Printf("Error: The request exceeded time out of %.1f seconds.\n", client.Timeout.Seconds())
		} else {
			fmt.Printf("Error: Failed to make a request: %s\n", err)
		}
		os.Exit(1)
	}
	defer resp.Body.Close()

	// Handle response status code
	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusBadRequest:
			fmt.Printf("Error: AOC_SESSION is expired.\n")
			os.Exit(1)
		case http.StatusNotFound:
			fmt.Printf("Error: Required input not found. The puzzle may not be available yet.\n")
			os.Exit(1)
		default:
			fmt.Printf("Error: Unexpected status code: %d\n", resp.StatusCode)
			os.Exit(1)
		}
	}

	// Create file for the input and write the response
	filename := fmt.Sprintf("aoc_%d_%d.txt", year, day)
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error: Failed to create input file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Printf("Error: Failed to write response to input file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully created %s\n", file.Name())
}
