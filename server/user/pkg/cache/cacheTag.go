package cache

import "strings"

const (
	SEP             = "_"
	CACHE_TAG       = "cache_tag"
	CACHE_USER      = "user"
	CACHE_PHONE     = "phone"
	CACHE_SIGNATURE = "signature"
)

func getCacheKey(keyBase, key string) string {
	keys := []string{
		keyBase,
		key,
	}
	return strings.Join(keys, SEP)
}
