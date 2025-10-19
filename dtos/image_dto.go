package dtos

type ImageDto struct {
	URLs   ImageURLs   `json:"urls"`
	Author ImageAuthor `json:"author"`
}

type ImageURLs struct {
	Small   string `json:"small"`
	Regular string `json:"regular"`
	Raw     string `json:"raw"`
}

type ImageAuthor struct {
	Name       string `json:"name"`
	Username   string `json:"username"`
	ProfileURL string `json:"profileURL"`
}
