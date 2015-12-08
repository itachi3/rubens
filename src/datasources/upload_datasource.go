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
	"encoding/json"
)

// In MB
const MAX_SIZE int = 2

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

func (ds UploadDataSource) checkSize(w http.ResponseWriter, r *http.Request) bool {
	size, err := strconv.Atoi(r.Header.Get("Content-Length"));

	/* Check for the accepted size of the file */
	if err != nil || size <= 0 {
		log.Println("Invalid content length : ", err)
		utils.WrapResponse(w, nil, http.StatusLengthRequired)
		return false
	}

	// Converting to MB
	size = size / (1024 * 1024)
	if size > MAX_SIZE {
		log.Println("Request size too large Size: ", size, "(MB)")
		utils.WrapResponse(w, utils.GetErrorContent(1), http.StatusBadRequest)
		return false
	}

	return true
}

func (ds UploadDataSource) UploadFile(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	if !ds.checkSize(w, r) {
		return nil
	}

	r.ParseMultipartForm(32 << 15)
	file, handler, err := r.FormFile("imageFile")
	if err != nil {
		log.Println("Error opening multipart file : ", err)
		return err
	}
	defer file.Close()

	format := handler.Header["Content-Type"]
	if !utils.Search(format, utils.ACCEPTED_FORMAT) {
		log.Println("Unknown format: ", format)
		utils.WrapResponse(w, nil, http.StatusUnsupportedMediaType)
		return nil
	}

	/* Process request information details */
	requestString := r.FormValue("requestBody");
	if requestString == "" {
		log.Println("No configuration details provided");
		utils.WrapResponse(w, utils.GetErrorContent(3), http.StatusBadRequest)
		return nil
	}

	var config models.UploadModel;
	json.Unmarshal([]byte(requestString), &config)

	// Decode file based on the format of the image

	/*Dump to file
	out , err:= os.Create("test.png")
	io.Copy(out, file)
	log.Println(err)*/

	fileType := strings.Split(format[0], "/")
	img, err := ds.helper.Decode(file, fileType[1])
	if err != nil {
		return err
	}
	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)
	ds.helper.SetFileName(timeStamp)

	/*
	   Use encoders to encode the content to the file and upload to s3
	   Note: use io.Copy to just dump the body to file. Supports huge file
	*/
	buf := new(bytes.Buffer)
	err = ds.helper.Encode(buf, img, config.GetConversionType())
	if err != nil {
		log.Println("Error encoding original : ", err)
		return err
	}

	fileName := config.GetFilePath() + utils.S3_SEPARATOR + timeStamp + "." + config.GetConversionType()
	err = ds.helper.UploadToS3(buf, fileName, "image/" + config.GetConversionType())
	if err != nil {
		return err
	}

	/*
			1. Send response after uploading original
		   	2. So that client doesn't have to wait for compressions
			3. Store the path in redis and send back the location response
	*/
	s3Path := ds.helper.GetS3Path()
	dynamicFileName := s3Path + config.GetFilePath() + utils.S3_SEPARATOR + timeStamp + "_{width}x{height}." + config.GetConversionType()
	_, err = ds.conn.RedisConn.Do("RPUSH", config.GetFilePath(), dynamicFileName)
	if err != nil {
		log.Println("Error writing (redis) : ", err)
	}
	successResponse := models.UploadResponse{
		ImageURL: s3Path + fileName,
		DynamicImageURL: dynamicFileName,
	}
	utils.WrapResponse(w, successResponse, http.StatusCreated)

	/*
	   Scale images to most expected dimensions and upload to S3
	*/


	ds.helper.ScaleImage(img, config)

	return nil
}