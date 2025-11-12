package bangumi

func BangumiCharacterCreate() *CharacterRequest {
	return &CharacterRequest{
		Keyword: "",
		Filter: Filter{
			NSFW: false,
		},
	}
}

func NewCharacter() *Character {
	return &Character{
		ID:        0,
		Name:      "未收錄",
		NameCN:    "",
		Aliases:   []string{"無"},
		Age:       "未收錄",
		Image:     "",
		Summary:   "無相關介紹",
		Gender:    "未收錄",
		BirthDay:  "未收錄",
		Height:    "未收錄",
		Weight:    "未收錄",
		BWH:       "未收錄",
		BloodType: "未收錄",
		Other:     []string{},
		Game:      []string{},
		CV:        []string{},
	}
}
