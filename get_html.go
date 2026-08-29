package main

import (
	"fmt"
	"errors"
	"net/http"
	"io"
	"os"
	"strings"
	"strconv"
	
	"context"
	"time"

	"github.com/chromedp/chromedp"

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

func getHTMLdp(rawURL string) (string, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),          // try visible first (more successful)
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36"),
		chromedp.WindowSize(1920, 1080),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	var html string

	err := chromedp.Run(ctx,
		chromedp.Navigate(rawURL),

		// Wait until the page is really ready
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.Sleep(6*time.Second), // give Cloudflare + page time to settle

		// More reliable way to get the full HTML
		chromedp.Evaluate(`document.documentElement.outerHTML`, &html),
	)

	if err != nil {
		return "", fmt.Errorf("chromedp error: %w", err)
	}

	// Detect if we are still on the challenge page
	if strings.Contains(html, "you are a real person") || strings.Contains(html, "cf-challenge") {
		return "", errors.New("still on Cloudflare challenge page")
	}

	return html, nil
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
