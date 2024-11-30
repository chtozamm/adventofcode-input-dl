package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

const (
	minYear = 2015
	maxYear = 2024
	minDay  = 1
	maxDay  = 25
)

var outputFilename string

var rootCmd = &cobra.Command{
	Use:   "adventofcode <year> <day>",
	Short: "Download Advent of Code input",
	Long:  "A command-line tool for downloading the input for the Advent of Code challenges.",
	RunE: func(cmd *cobra.Command, args []string) error {
		switch len(args) {
		case 0:
			return fmt.Errorf("missing arguments")
		case 1:
			return fmt.Errorf("not enough arguments: missing <day>")
		}

		year, err := strconv.Atoi(args[0])
		if err != nil || year < minYear || year > maxYear {
			return fmt.Errorf("year must be between %d and %d", minYear, maxYear)
		}

		day, err := strconv.Atoi(args[1])
		if err != nil || day < minDay || day > maxDay {
			return fmt.Errorf("day must be between %d and %d", minDay, maxDay)
		}

		session := os.Getenv("AOC_SESSION")
		if session == "" {
			return fmt.Errorf("AOC_SESSION not found in environment variables")
		}

		err = fetchInput(year, day, session)
		return err
	},
}

func init() {
	rootCmd.Flags().StringVarP(&outputFilename, "output", "o", "", "Output filename (default: aoc_<year>_<day>.txt)")
}

func fetchInput(year, day int, session string) error {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://adventofcode.com/%d/day/%d/input", year, day), nil)
	if err != nil {
		return fmt.Errorf("create HTTP request: %s", err)
	}

	req.AddCookie(&http.Cookie{Name: "session", Value: session})

	resp, err := client.Do(req)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return fmt.Errorf("request exceeded time out of %.1f seconds", client.Timeout.Seconds())
		} else {
			return fmt.Errorf("make request: %s", err)
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return handleResponseError(resp.StatusCode)
	}

	filename := outputFilename
	if filename == "" {
		filename = fmt.Sprintf("aoc_%d_%d.txt", year, day)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("create input file: %v", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("write response to input file: %v", err)
	}

	fmt.Printf("Successfully written %s\n", file.Name())
	return nil
}

func handleResponseError(statusCode int) error {
	switch statusCode {
	case http.StatusBadRequest:
		return fmt.Errorf("AOC_SESSION is expired")
	case http.StatusNotFound:
		return fmt.Errorf("required input not found. The puzzle may not be available yet")
	default:
		return fmt.Errorf("unexpected status code: %d", statusCode)
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
