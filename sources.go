package ipsumru

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
)

const (
	nameNews2023 = "rus_news_2023_1M-sentences.txt"
	nameNews2024 = "rus_news_2023_1M-sentences.txt"
	// https://wortschatz.uni-leipzig.de/en/download/Russian?utm_source=chatgpt.com
	baseURL = "https://github.com/fpawel/ipsumru/releases/download/v0.0.0.1/"
	url2023 = baseURL + nameNews2023
	url2024 = baseURL + nameNews2024
)

func ensureFile(url, path string) error {
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}
	if errors.Is(err, os.ErrNotExist) {
		if err = download(url, path); err != nil {
			return fmt.Errorf("could not download %q to %q: %w", url, path, err)
		}
		return nil
	}
	return fmt.Errorf("could not stat %q to %q: %w", url, path, err)
}

func download(url, path string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("http get: %w", err)
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			slog.Error("Failed to close response body", "error", err)
		}
	}()

	var fOut *os.File
	if fOut, err = os.Create(path); err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer func() {
		if err = fOut.Close(); err != nil {
			slog.Error("Failed to close file", "error", err)
		}
	}()
	if _, err = io.Copy(fOut, resp.Body); err != nil {
		return fmt.Errorf("copy response body to file: %w", err)
	}
	return nil
}
