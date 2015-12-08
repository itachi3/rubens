package datasources

import (
	"agentdesks"
	"agentdesks/helpers"
	"agentdesks/models"
	"agentdesks/utils"
	"github.com/garyburd/redigo/redis"
	"log"
	"net/http"
	"strings"
)

type DeleteDataSource struct {
	connect *models.Connections
	helper  *helpers.ImageHelper
}

func NewDeleteDataSource(conn *models.Connections, conf *agentdesks.Config) *DeleteDataSource {
	return &DeleteDataSource{
		connect: conn,
		helper:  helpers.NewImageHelper(conf),
	}
}

func (deleteDs *DeleteDataSource) DeleteImages(w http.ResponseWriter, r *http.Request) error {
	queryValues := r.URL.Query()
	imageKey := queryValues.Get("key")
	var redisKey, s3Key string
	if imageKey != "" {
		chunk := strings.Split(imageKey, deleteDs.helper.Config.GetAmazonS3Bucket() + "/")
		chunk = strings.Split(chunk[len(chunk)-1], "_")
		s3Key = chunk[0]
		chunk = strings.Split(s3Key, "/")
		redisKey = chunk[0] + "/" + chunk[1]
		log.Println(s3Key)
	} else {
		redisKey = queryValues.Get("pathKey")
		s3Key = redisKey
	}

	fileLocation, err := redis.Strings(deleteDs.connect.RedisConn.Do("LRANGE", redisKey, 0, -1))
	if err != nil {
		log.Println("Error retrieving records (redis) : ", err)
		return err
	}

	if len(fileLocation) == 0 {
		utils.WrapResponse(w, nil, http.StatusNotFound)
		return nil
	}

	err = deleteDs.helper.BatchDelete(s3Key)
	if err == nil {
		utils.WrapResponse(w, nil, http.StatusOK)
		deleteDs.connect.RedisConn.Do("DEL", redisKey)
		// Remove only that key from redis and push back the list
		if imageKey != "" {
			for _, elem := range fileLocation {
				if !strings.EqualFold(elem, imageKey) {
					deleteDs.connect.RedisConn.Do("RPUSH", redisKey, elem)
				}
			}
		}
	}
	return err
}
