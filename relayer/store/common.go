package store

import "strings"

func GetKey(keys []string) []byte {
	return []byte(strings.Join(keys, "-"))
}
