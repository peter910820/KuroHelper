package handlers

// 資料分頁
func pagination[T any](result *[]T, page int, useCache bool) bool {
	resultLen := len(*result)
	expectedMin := page * 15
	expectedMax := page*15 + 15

	if !useCache || page == 0 {
		if resultLen > 15 {
			*result = (*result)[:15]
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
