package fileutils

import (
	"bufio"
	"os"
)

func ReadFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var html string = ""
	for scanner.Scan() {
		line := scanner.Text() + "\n" // <-- Salto de linea
		html += line
	}

	return html, nil
}
