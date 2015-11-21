package helpers

import (
	"agentdesks"
	"agentdesks/utils"
	"bytes"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/julienschmidt/httprouter"
	"github.com/nfnt/resize"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type ImageHelper struct {
	config   *agentdesks.Config
	s3Conn   *s3.S3
	fileName string
}

func NewImageHelper(conf *agentdesks.Config) *ImageHelper {
	return &ImageHelper{
		config: conf,
		s3Conn: s3.New(session.New(), &aws.Config{Region: aws.String(conf.GetAmazonS3Region())}),
	}
}

func (ih *ImageHelper) GetPath(p httprouter.Params) (string, string) {
	userId := p.ByName(utils.USER_ID)
	propertyId := p.ByName(utils.PROPERTY_ID)
	return utils.IMAGES_BASE_URL +
			ih.config.GetAmazonS3Bucket() +
			utils.S3_SEPARATOR,
		userId +
			utils.S3_SEPARATOR +
			propertyId
}

func (ih *ImageHelper) Decode(r *http.Request, fileType string) (image.Image, error) {
	if fileType == "png" {
		return png.Decode(r.Body)
	} else if fileType == "jpeg" || fileType == "jpg" {
		return jpeg.Decode(r.Body)
	} else if fileType == "gif" {
		return gif.Decode(r.Body)
	}
	log.Println("Unsupported media for decoding ", fileType)
	return nil, errors.New("Unsupported media type")
}

func (ih *ImageHelper) Encode(w io.Writer, img image.Image, fileType string) error {
	if fileType == "png" {
		return png.Encode(w, img)
	} else if fileType == "jpeg" || fileType == "jpg" {
		return jpeg.Encode(w, img, nil)
	} else if fileType == "gif" {
		return gif.Encode(w, img, nil)
	}
	log.Println("Unsupported media for encoding ", fileType)
	return errors.New("Unsupported media type")
}

func (ih *ImageHelper) ScaleImage(path string, img image.Image, fileFormat string) {
	var width, height string
	var w, h int
	fileType := strings.Split(fileFormat, "/")
	for index, _ := range ih.config.Img.Width {
		width = ih.config.Img.Width[index]
		height = ih.config.Img.Height[index]
		buf := new(bytes.Buffer)
		w, _ = strconv.Atoi(width)
		h, _ = strconv.Atoi(height)
		m := resize.Resize(uint(w), uint(h), img, resize.MitchellNetravali)
		ih.Encode(buf, m, fileType[1])
		pathname := path + utils.S3_SEPARATOR + width + "x" + height + "_" + ih.fileName + "." + fileType[1]
		go ih.UploadToS3(buf, pathname, fileFormat)
	}
}

func (ih *ImageHelper) UploadToS3(buf *bytes.Buffer, pathname string, fileFormat string) error {
	params := &s3.PutObjectInput{
		Bucket:      aws.String(ih.config.GetAmazonS3Bucket()),
		Key:         aws.String(pathname),
		ACL:         aws.String("public-read"),
		Body:        bytes.NewReader(buf.Bytes()),
		ContentType: aws.String(fileFormat),
	}
	_, err := ih.s3Conn.PutObject(params)
	if err != nil {
		log.Println("Error uploading to S3 : ", err)
		return errors.New("Oops ! Something went wrong")
	}
	return nil
}

func (ih *ImageHelper) BatchDelete(bucket string, key string) error {
	params := &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String(key),
	}

	res, err := ih.s3Conn.ListObjects(params)

	if err == nil && res.Contents != nil {
		var objectIdentifiers []*s3.ObjectIdentifier
		var batchDelete *s3.Delete
		var batchParams *s3.DeleteObjectsInput

		for _, object := range res.Contents {
			objectIdentifiers = append(objectIdentifiers, &s3.ObjectIdentifier{Key: object.Key})
		}

		batchSize := 2 //s3Max batch size
		for i:=0; i<len(objectIdentifiers); i+=batchSize {
			if i+batchSize > len(objectIdentifiers) {
				batchDelete = &s3.Delete{
					Objects: objectIdentifiers,
				}
			} else {
				batchDelete = &s3.Delete{
					Objects: objectIdentifiers[i:i+batchSize],
				}
			}

			batchParams = &s3.DeleteObjectsInput{
				Bucket: aws.String(bucket),
				Delete: batchDelete,
			}

			_, err := ih.s3Conn.DeleteObjects(batchParams)
			if err != nil {
				log.Println("Error deleting objects (s3) : ", err)
				return err
			}
		}
	}

	return err
}

func (ih *ImageHelper) SetFileName(name string) {
	ih.fileName = name
}
