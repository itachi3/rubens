package datasources

import (
	"agentdesks"
	"agentdesks/helpers"
	"agentdesks/models"
	"agentdesks/utils"
	"github.com/garyburd/redigo/redis"
	"log"
	"net/http"
)

type DownloadDataSource struct {
	conn   *models.Connections
	helper *helpers.ImageHelper
}

func NewDownloadDataSource(conn *models.Connections, config *agentdesks.Config) *DownloadDataSource {
	return &DownloadDataSource{
		conn:   conn,
		helper: helpers.NewImageHelper(config),
	}
}

func (download DownloadDataSource) GetImagesLocation(w http.ResponseWriter, r *http.Request) error {
	queryValues := r.URL.Query()
	key := queryValues.Get("key")
	//Always get the latest location
	fileLocation, err := redis.Strings(download.conn.RedisConn.Do("LRANGE", key, 0, -1))
	if err != nil {
		log.Println("Error retrieving record (redis) : ", err)
		return err
	}

	if len(fileLocation) == 0 {
		utils.WrapResponse(w, nil, http.StatusNotFound)
		return nil
	}

	success := models.ImageUrlResponse{
		DynamicImageURL: fileLocation,
	}
	utils.WrapResponse(w, success, http.StatusOK)
	return nil
}
