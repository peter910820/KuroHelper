package erogs

import (
	"kurohelper/cache"
)

func zhtwToJp(search string) string {
	runes := []rune(search)
	for i, r := range runes {
		if jp, ok := cache.ZhtwToJp[r]; ok {
			runes[i] = jp
		}
	}
	return string(runes)
}
