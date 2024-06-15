package pkg

import (
	"fmt"
	http2 "github.com/raitonbl/coverup/pkg/http"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func ReadFromURL(httpClient http2.Client, uri string) ([]byte, error) {
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get content from %s: %w", uri, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}

func ReadFromFile(directory string, uri string) ([]byte, error) {
	path := uri[7:]
	if directory != "" {
		path = filepath.Join(directory, path)
	}
	return os.ReadFile(path)
}
