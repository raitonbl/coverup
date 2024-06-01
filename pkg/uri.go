package pkg

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func ReadURI(directory, uri string) ([]byte, error) {
	if strings.HasPrefix(uri, "http://") || strings.HasPrefix(uri, "https://") {
		return readFromURL(uri)
	} else if strings.HasPrefix(uri, "file://") {
		return readFromFile(directory, uri)
	} else {
		return nil, fmt.Errorf("unsupported URI: %s", uri)
	}
}

func readFromURL(uri string) ([]byte, error) {
	resp, err := http.Get(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to get content from %s: %w", uri, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}

func readFromFile(directory string, uri string) ([]byte, error) {
	path := uri[7:]
	if directory != "" {
		path = filepath.Join(directory, path)
	}
	return os.ReadFile(path)
}
