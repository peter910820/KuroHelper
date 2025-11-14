package vndb

var (
	Gender = map[string]string{
		"m": "男",
		"f": "女",
		"o": "非二元",
		"a": "不明",
	}
	Sex = map[string]string{
		"m": "男",
		"f": "女",
		"b": "非二元",
		"a": "不明",
	}
	Role = map[string]string{
		"main":    "主角",
		"primary": "主角",
		"side":    "配角",
		"appears": "配角",
	}
	RolePriority = map[string]int{
		"main":    1,
		"primary": 2,
		"side":    3,
		"appears": 4,
	}
)
