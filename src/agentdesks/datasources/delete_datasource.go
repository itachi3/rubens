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
	userId := queryValues.Get(utils.USER_ID)
	propertyId := queryValues.Get(utils.PROPERTY_ID)
	key := userId + "/" + propertyId
	fileLocation, err := redis.Strings(deleteDs.connect.RedisConn.Do("LRANGE", key, 0, -1))
	if err != nil {
		log.Println("Error retrieving records (redis) : ", err)
		return err
	}

	if len(fileLocation) == 0 {
		utils.WrapResponse(w, nil, http.StatusNotFound)
		return nil
	}

	return deleteDs.helper.BatchDelete("image.agentdesks.com",key)
}
