package main

import (
	"fmt"
	"errors"
	"net/http"
	"io"
	"os"
	"strings"
	"strconv"
)

//-------------------------------------------------------------------------

func getHTML(rawURL string) (string, error) {
	res, err := http.Get(rawURL)
	if err != nil {
	    return "", err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return "", errors.New(fmt.Sprintf("error status: %s", res.Status))
	}
	if !strings.Contains(res.Header.Get("Content-Type"), "text/html")  {
		return "", errors.New(fmt.Sprintf("error Header: not text/html"))
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	return string(data), nil
}

//-------------------------------------------------------------------------

func getArg() (string, int, int, error) {
	arg := os.Args[1:]
	if len(arg) < 1 {
		fmt.Println("no website provided -> website, concurrency, max_visit")
		os.Exit(1)
	} else if len(arg) < 3 {
		fmt.Println("too few arguments provided -> website, concurrency, max_visit")
		os.Exit(1)
	} else if len(arg) > 3 {
		fmt.Println("too many arguments provided -> website, concurrency, max_visit")
		os.Exit(1)
	} else {
		fmt.Println("starting crawl of: " + arg[0])
	}
	num_1, err := strconv.Atoi(arg[1])
	if err != nil {
		return "", 0, 0, err
	}
	num_2, err := strconv.Atoi(arg[2])
	if err != nil {
		return "", 0, 0, err
	}
	return arg[0], num_1, num_2, nil
}

//-------------------------------------------------------------------------
