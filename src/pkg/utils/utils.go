package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
	"rttmas-backend/pkg/utils/logger"
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
