package handlers

const (
	Played = 1 << iota
	Wish
)

// 資料分頁
func pagination[T any](result *[]T, page int, useCache bool) bool {
	resultLen := len(*result)
	expectedMin := page * 10
	expectedMax := page*10 + 10

	if !useCache || page == 0 {
		if resultLen > 10 {
			*result = (*result)[:10]
			return true
		}
		return false
	} else {
		if resultLen > expectedMax {
			*result = (*result)[expectedMin:expectedMax]
			return true
		} else {
			*result = (*result)[expectedMin:]
			return false
		}
	}
}
