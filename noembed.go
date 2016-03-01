package main

type NoEmbed struct {
	AuthorName string `json:"author_name"`
	Width int `json:"width"`
	AuthorURL string `json:"author_url"`
	ProviderURL string `json:"provider_url"`
	Version string `json:"version"`
	ThumbnailWidth int `json:"thumbnail_width"`
	ProviderName string `json:"provider_name"`
	ThumbnailURL string `json:"thumbnail_url"`
	Height int `json:"height"`
	ThumbnailHeight int `json:"thumbnail_height"`
	HTML string `json:"html"`
	URL string `json:"url"`
	Title string `json:"title"`
	Type string `json:"type"`
}
