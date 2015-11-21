package datasources

import (
	"agentdesks"
	"agentdesks/helpers"
	"agentdesks/models"
	"agentdesks/utils"
	"bytes"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// In MB
const MAX_SIZE int = 1

type UploadDataSource struct {
	conn   *models.Connections
	helper *helpers.ImageHelper
}

func NewUploadDataSource(conn *models.Connections, config *agentdesks.Config) *UploadDataSource {
	return &UploadDataSource{
		conn:   conn,
		helper: helpers.NewImageHelper(config),
	}
}

func (ds UploadDataSource) ValidateData(w http.ResponseWriter, r *http.Request, p httprouter.Params) (bool, error) {
	size := r.Header.Get(utils.SIZE)
	/* Check for the accepted size of the file */
	if size == "0" {
		log.Println("No file to upload")
		utils.WrapResponse(w, nil, http.StatusLengthRequired)
		return false, nil
	}

	// Converting to MB
	sizeInt, err := strconv.Atoi(size)
	if err != nil {
		log.Println("Error parsing content-length :", size)
		return false, err
	}
	sizeInt = sizeInt / (1024 * 1024)
	if sizeInt > MAX_SIZE {
		log.Println("Image too large Size: ", sizeInt, "(MB)")
		utils.WrapResponse(w, utils.GetErrorContent(1), http.StatusBadRequest)
		return false, nil
	}

	format := r.Header.Get(utils.FORMAT)
	if !utils.Search(format, utils.ACCEPTED_FORMAT) {
		log.Println("Unknown format: ", format)
		utils.WrapResponse(w, nil, http.StatusUnsupportedMediaType)
		return false, nil
	}

	return true, nil
}

func (ds UploadDataSource) UploadFile(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	// Decode file based on the format of the image
	fileType := strings.Split(r.Header.Get(utils.FORMAT), "/")

	//Dump to file
	/*file , err:= os.Create("test.png")
	io.Copy(file, r.Body)
	log.Println(err)*/

	img, err := ds.helper.Decode(r, fileType[1])
	if err != nil {
		return err
	}
	defer r.Body.Close()

	// Create the file path for writing
	baseURL, path := ds.helper.GetPath(p)
	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)
	ds.helper.SetFileName(timeStamp)

	//TODO Upload the image to Ec2 instead
	/*
	   Use encoders to encode the content to the file
	   Note: use io.Copy to just dump the body to file. Supports huge file
	*/
	buf := new(bytes.Buffer)
	err = ds.helper.Encode(buf, img, fileType[1])
	if err != nil {
		log.Println("Error encoding original : ", err)
		return err
	}

	fileName := path + utils.S3_SEPARATOR + timeStamp + "." + fileType[1]
	err = ds.helper.UploadToS3(buf, fileName, r.Header.Get(utils.FORMAT))
	if err != nil {
		return err
	}
	/*
			1. Send response after uploading original
		   	2. So that client doesn't have to wait for compressions
			3. Store the path in redis and send back the location response
	*/
	key := p.ByName(utils.IMAGE_KEY)
	value := baseURL + path + "/{width}x{height}_" + timeStamp + "." + fileType[1]
	_, err = ds.conn.RedisConn.Do("RPUSH", key, value)
	if err != nil {
		log.Println("Error writing (redis) : ", err)
	}
	successResponse := map[string]interface{} {
		"FileURL" : value,
	}
	utils.WrapResponse(w, successResponse, http.StatusCreated)

	/*
	   Scale images to most expected dimensions and upload to S3
	*/
	ds.helper.ScaleImage(path, img, r.Header.Get(utils.FORMAT))

	return nil
}
