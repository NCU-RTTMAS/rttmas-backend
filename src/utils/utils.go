package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
	"rttmas-backend/utils/logger"
	"strings"
)

func Jsonalize(v any) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		logger.Info(err)
	}
	return string(b)
}

func GetWorkingDirectory() string {
	binaryDir, err := os.Executable()
	if err != nil {
		logger.Error(err)
	}
	dir := filepath.Dir(binaryDir)
	if err != nil {
		logger.Error(err)
	}
	return dir
}

func ParseCommaSeparatedString(input string) []string {
	// Trim any leading or trailing whitespace
	input = strings.TrimSpace(input)

	// Split the string by commas
	parts := strings.Split(input, ",")

	// Trim whitespace from each part
	for i, part := range parts {
		parts[i] = strings.TrimSpace(part)
	}

	return parts
}
