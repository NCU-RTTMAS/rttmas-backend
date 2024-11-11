package utils

import (
	"encoding/json"
	"rttmas-backend/pkg/utils/logger"
)

func Jsonalize(v any) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		logger.Info(err)
	}
	return string(b)
}
