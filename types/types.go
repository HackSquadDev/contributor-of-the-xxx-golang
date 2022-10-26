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