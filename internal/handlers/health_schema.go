package handlers

type InfoResponseSchema struct {
	Service     string `json:"service"`
	Name        string `json:"tom_name"`
	Version     string `json:"version"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
