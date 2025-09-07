package utils

import "os"

// ReadFileAsString is a function that reads a file, then converts it into a string
func ReadFileAsString(fileName string) (string, error) {
	file, err := os.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	return string(file), nil
}
