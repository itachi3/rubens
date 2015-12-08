package agentdesks

type (

	Config struct {
		Port string
		Logs Log       `json:"logs"`
		Datastore DataStores `json:"dataStores"`
		BucketPool map[string]ImageBucket    `json:"imageBucket"`
	}

	Log struct {
		ErrorLog  string
		AccessLog string
	}

	DataStores struct {
		Redis REDIS `json:"redis"`
		AwsS3 AmazonS3 `json:"amazonS3"`
	}

	REDIS struct {
		Protocol string
		Port     string
	}

	AmazonS3 struct {
		Region string
		Bucket string `json:"bucketName"`
	}

	ImageBucket struct {
		Width  []string
		Height []string
	}
)

func (config *Config) GetRedisPort() string {
	return config.Datastore.Redis.Port
}

func (config *Config) GetRedisProtocol() string {
	return config.Datastore.Redis.Protocol
}

func (config* Config) GetAmazonS3Region() string {
	return config.Datastore.AwsS3.Region
}

func (config* Config) GetAmazonS3Bucket() string {
	return config.Datastore.AwsS3.Bucket
}

func (config *Config) GetImageBucket(key string) ImageBucket {
	return config.BucketPool[key]
}