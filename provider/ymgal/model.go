package ymgal

type (
	basicResp[T any] struct {
		Success bool `json:"success"`
		Code    int  `json:"code"`
		Data    T    `json:"data"`
	}

	tokenResp struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
		Scope       string `json:"scope"`
	}

	randomGameResp struct {
		GID         int    `json:"gid"`
		DeveloperID int    `json:"developerId"`
		Name        string `json:"name"`
		ChineseName string `json:"chineseName"`
		HaveChinese bool   `json:"haveChinese"`
		MainImg     string `json:"mainImg"`
		ReleaseDate string `json:"releaseDate"`
		State       string `json:"state"`
	}
)
