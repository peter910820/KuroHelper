package vndb

import vndbmodels "kurohelper/models/vndb"

func StaffFuzzySearch(keyword string, roleType string) {
	req := vndbmodels.VndbCreate()

	filters := []interface{}{}
	if roleType != "" {
		filters = append(filters, "and")
		// 傳進來的直接就是API篩選項規格的字串
		filters = append(filters, []string{"type", "=", roleType})
		filters = append(filters, []string{"search", "=", keyword})
	} else {
		filters = []interface{}{"search", "=", keyword}
	}

	req.Filters = filters

}
