package service

import (
	"os"
	"strings"
)

func createFile(code string, lang string) (*os.File, error) {
	var filename string
	if lang == "Go" {
		filename = "sourceFile." + "go"
	} else {
		filename = "sourceFile." + "java"
	}

	f, err := os.Create(filename)
	if err != nil {
		return f, err
	}

	lines := strings.Split(code, "\n")
	for _, line := range lines {
		f.WriteString(line + "\n")
	}
	return f, err
}
