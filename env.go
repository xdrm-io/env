package env

import (
	"os"
)

// Read returns :
//   - the value of the environment variable {key} if it exists
//   - the contents of the file located at the path from the environment variable
//     {key}_FILE if it exists
func Read(key string) (string, bool) {
	if raw, ok := os.LookupEnv(key); ok {
		return raw, true
	}

	path, ok := os.LookupEnv(key + "_FILE")
	if !ok {
		return "", false
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		return "", false
	}
	return string(raw), true
}
