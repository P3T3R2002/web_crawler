package main

import (
	"fmt"
	"errors"
	"net/http"
	"io"
	"os"
	"strings"
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

func getArg() string{
	arg := os.Args[1:]
	if len(arg) < 1 {
		fmt.Println("no website provided")
		os.Exit(1)
	} else if len(arg) > 1 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	} else {
		fmt.Println("starting crawl of: " + arg[0])
	}
	return arg[0]
}

//-------------------------------------------------------------------------
