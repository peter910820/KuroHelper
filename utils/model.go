package utils

type (
	CustomIDArgs interface {
		PageArgs | SortArgs | AddWishArgs | AddHasPlayedArgs
	}

	NewCustomID[T CustomIDArgs] struct {
		CommandName string
		Type        CustomIDType
		Value       T
	}

	PageArgs struct {
		CacheID string
		Page    int
	}

	SortArgs struct {
		CacheID       string
		SortMethod    string
		SortDirection string
	}

	AddWishArgs struct {
		CacheID     string
		ConfirmMark bool
	}

	AddHasPlayedArgs struct {
		CacheID     string
		ConfirmMark bool
	}
)
