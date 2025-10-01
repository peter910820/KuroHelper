package erogs

import "kurohelper/config"

func zhtwToJp(search string) string {
	runes := []rune(search)
	for i, r := range runes {
		if jp, ok := config.ZhtwToJp[r]; ok {
			runes[i] = jp
		}
	}
	return string(runes)
}
