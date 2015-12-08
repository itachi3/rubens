package models

type (
	UploadModel struct {
		Path string `json:"filePath"`
		Config Config
	}

	Config struct {
		Bucket string
		ConvertTo string
	}
)

func (model UploadModel) GetConversionType() string {
	return model.Config.ConvertTo
}

func (model UploadModel) GetFilePath() string {
	return model.Path
}

func (model UploadModel) GetBucket() string {
	return model.Config.Bucket
}
