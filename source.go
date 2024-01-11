package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

type Source interface {
	Read() (string, error)
	String() string
}

func NewSource(input string) (Source, error) {

	url, err := url.Parse(input)
	if err == nil {
		return WebSource{url: *url}, nil
	}

	if filepath.IsAbs(input) || filepath.IsLocal(input) {
		return FileSource{path: input}, nil
	}

	return nil, fmt.Errorf("cannot deduce source type of '%s'", input)
}

type FileSource struct {
	path string
}

func (f FileSource) Read() (string, error) {
	raw, err := os.ReadFile(f.path)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}

func (f FileSource) String() string {
	return f.path
}

type WebSource struct {
	url url.URL
}

func (w WebSource) Read() (string, error) {
	res, err := http.Get(w.url.String())
	if err != nil {
		return "", err
	}

	content, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return "", err
	}
	return string(content), nil

}

func (w WebSource) String() string {
	return w.url.String()
}
