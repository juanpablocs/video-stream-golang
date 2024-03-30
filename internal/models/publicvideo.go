package models

type PublicVideo struct {
	ID          string            `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Thumbnail   map[string]string `json:"thumbnail"`
	Duration    float64           `json:"duration"`
	VttTrack    string            `json:"vttTrack"`
	Status      VideoStatus       `json:"status"`
	CreatedAt   string            `json:"createdAt"`
}

type VideoListResponse struct {
	CurrentPage int64         `json:"currentPage"`
	NextPage    int64         `json:"nextPage"`
	Total       int64         `json:"total"`
	Videos      []PublicVideo `json:"videos"`
}
