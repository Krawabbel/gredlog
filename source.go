package main

import (
	"os"
)

type Source interface {
	Read() (string, error)
	String() string
}

type FileSource struct {
	Path string
}

func (f FileSource) Read() (string, error) {
	raw, err := os.ReadFile(f.Path)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}

func (f FileSource) String() string {
	return f.Path
}

// type WebSource struct {
// 	Path string
// }

// func (w WebSource) Read() (string, error) {
// 	return "", fmt.Errorf("web source not yet implemented")
// }

// func (w WebSource) String() string {
// 	return w.Path
// }
