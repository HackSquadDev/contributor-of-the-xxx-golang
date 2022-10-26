package types

type SearchResponse struct {
	Search struct {
		Nodes    []RepositoryResponse
		PageInfo struct {
			HasNextPage bool
			EndCursor   string
		}
	}
}
type RepositoryResponse struct {
	Title  string
	Url    string
	Author struct {
		AvatarURL string
		Login     string
		Url       string
	}
}
type OrganizationResponse struct {
	Organization struct {
		Name      string
		Url       string
		AvatarUrl string
		Login     string
	}
}

// data as got from imagga API
type ColorResponse struct {
	Result struct {
		Colors struct {
			BackgroundColors []struct {
				HTMLCode string `json:"html_code"`
			} `json:"background_colors"`
			ForegroundColors []struct {
				HTMLCode string `json:"html_code"`
			} `json:"foreground_colors"`
			ImageColors []struct {
				HTMLCode string `json:"html_code"`
			} `json:"image_colors"`
		} `json:"colors"`
	} `json:"result"`
}
type ThemeColors struct {
	BackgroundColors []string
	ForegroundColors []string
	ImageColors      []string
}
