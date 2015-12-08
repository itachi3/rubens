package models

type (
	// Structure of our resource
	HTTPErrorResponse struct {
		ErrorCode string `json:"errorCode"`
		Message   string `json:"message"`
	}

	ImageUrlResponse struct {
		DynamicImageURL []string `json:"dynamicImageURL"`
	}

	UploadResponse struct {
		ImageURL string `json:"imageURL"`
		DynamicImageURL string `json:"dynamicImageURL"`
	}
)
