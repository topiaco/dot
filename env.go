package dot

import "os"

func Getenv(key string, def ...string) string {
	val := os.Getenv(key)
	if val == "" {
		if len(def) > 0 {
			return def[0]
		}
	}

	return val
}
