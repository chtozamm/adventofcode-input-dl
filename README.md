# Advent of Code Input Downloader

A command-line tool for downloading the input for the Advent of Code challenges.

## Prerequisites

- [Go](https://go.dev/)
- A valid session cookie from [Advent of Code](https://adventofcode.com/). You can obtain this by logging into your account and copying the session cookie.
- Set the session cookie as an environmental variable `AOC_SESSION`.

## Installation

```
go install github.com/chtozamm/adventofcode-input-dl/cmd/adventofcode@latest
```

## Usage

```
adventofcode <year> <day>
```

## Examples

To download the input for Day 1 of Advent of Code 2023, run:

```
adventofcode 2023 1
```

You should see a message like:

```
Successfully created aoc_2023_day_1.txt
```
