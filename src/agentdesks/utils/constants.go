package utils

const (
	USER_AGENT      string = "X-User-Agent"
	FORMAT          string = "Content-Type"
	SIZE            string = "Content-Length"
	USER_ID         string = "userId"
	PROPERTY_ID     string = "propertyId"
	IMAGES_BASE_URL string = "https://s3.amazonaws.com/"
	S3_SEPARATOR    string = "/"
)

var ACCEPTED_FORMAT = []string{"image/jpeg", "image/png", "image/gif"}

var ErrorCodes = map[int]string{
	1: "File too large",
	2: "Invalid request headers",
}
