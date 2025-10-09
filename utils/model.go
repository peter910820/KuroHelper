package utils

type CustomIDType int

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
		GameID      int
		BrandID     int
		UserID      int
		ConfirmMark bool
	}

	AddHasPlayedArgs struct {
		GameID      int
		BrandID     int
		ConfirmMark bool
	}
)

const (
	CustomIDTypePage CustomIDType = iota + 1
	CustomIDTypeSort
	CustomIDTypeAddWish
	CustomIDTypeAddHasPlayed
)
