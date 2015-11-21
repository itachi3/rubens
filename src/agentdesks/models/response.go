package models

type (
	// Structure of our resource
	HTTPErrorResponse struct {
		ErrorCode string `json:"errorCode"`
		Message   string `json:"message"`
	}

	ImageUrlResponse struct {
		FileURL []string `json:"fileURL"`
	}
)
