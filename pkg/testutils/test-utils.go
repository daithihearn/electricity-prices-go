package testutils

import (
	"encoding/json"
	"os"
)

func ReadJsonStringFromFile(filePath string) string {
	bytes, _ := os.ReadFile(filePath)
	return string(bytes)
}

func ReadJSONFromFile(filePath string, result interface{}) error {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, result)
	if err != nil {
		return err
	}
	return nil
}
