# Advent of Code Input Downloader

A command-line tool for downloading the input for the Advent of Code challenges.

## Prerequisites

- Go is required to build or [install](#installation) the application.
- A valid session cookie from [Advent of Code](https://adventofcode.com/). You can obtain this by logging into your account and copying the session cookie.
- Set the session cookie as an environmental variable `AOC_SESSION`.

## Installation

Install with Go:

```sh
go install github.com/chtozamm/adventofcode-input-dl/cmd/adventofcode@latest
```

## Usage

```sh
adventofcode <year> <day>
```

To download the input for Day 1 of Advent of Code 2023, simply run:

```sh
adventofcode 2023 1
```

You should see a message like:

```sh
Successfully created aoc_2023_day_1.txt
```

## License

[MIT](https://github.com/chtozamm/adventofcode-input-dl/blob/main/LICENSE.md)
