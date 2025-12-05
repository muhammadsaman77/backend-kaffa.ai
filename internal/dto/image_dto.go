package dto

type CreateImageMetadata struct {
	ID           string `json:"id"`
	OriginalName string `json:"original_name"`
	Size         int32  `json:"size"`
	MimeType     string `json:"mime_type"`
	Path         string `json:"path"`
}
