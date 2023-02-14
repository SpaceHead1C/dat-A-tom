package handlers

type InfoResponseSchema struct {
	Service     string `json:"service"`
	Version     string `json:"version"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
