package http

type errorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

type gameListResponse struct {
	Data []int64 `json:"data"`
}

type gameInfo struct {
	ID          int64        `json:"id"`
	Title       string       `json:"title"`
	Kind        string       `json:"kind"`
	Description string       `json:"description"`
	Until       *int64       `json:"until"`
	URLs        gameInfoURLs `json:"urls"`
	Store       string       `json:"store"`
}

type gameInfoURLs struct {
	Default string `json:"default"`
	Browser string `json:"browser"`
	Org     string `json:"org"`
	Client  string `json:"client"`
}

type gameInfoResponse struct {
	Data map[string]gameInfo
}
