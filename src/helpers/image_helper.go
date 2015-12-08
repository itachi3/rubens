package helpers

import (
	"agentdesks"
	"agentdesks/utils"
	"bytes"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"mime/multipart"
	"agentdesks/models"
	"strconv"
	"github.com/nfnt/resize"
)

type ImageHelper struct {
	Config   *agentdesks.Config
	s3Conn   *s3.S3
	fileName string
}

func NewImageHelper(conf *agentdesks.Config) *ImageHelper {
	return &ImageHelper{
		Config: conf,
		s3Conn: s3.New(session.New(), &aws.Config{Region: aws.String(conf.GetAmazonS3Region())}),
	}
}

func (ih *ImageHelper) GetS3Path() string {
	return utils.IMAGES_BASE_URL + ih.Config.GetAmazonS3Bucket() + utils.S3_SEPARATOR
}

func (ih *ImageHelper) Decode(file multipart.File, fileType string) (image.Image, error) {
	if fileType == "png" {
		return png.Decode(file)
	} else if fileType == "jpeg" || fileType == "jpg" {
		return jpeg.Decode(file)
	} else if fileType == "gif" {
		return gif.Decode(file)
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

func (ih *ImageHelper) ScaleImage(img image.Image, config models.UploadModel) {
	bucket := ih.Config.GetImageBucket(config.GetBucket())
	var width, height string
	var w, h int
	fileType := config.GetConversionType()
	for index, _ := range bucket.Width {
		width = bucket.Width[index]
		height = bucket.Height[index]
		buf := new(bytes.Buffer)
		w, _ = strconv.Atoi(width)
		h, _ = strconv.Atoi(height)
		m := resize.Resize(uint(w), uint(h), img, resize.MitchellNetravali)
		ih.Encode(buf, m, fileType)
		pathname := config.GetFilePath() + utils.S3_SEPARATOR + ih.fileName + "_" + width + "x" + height + "." + fileType
		go ih.UploadToS3(buf, pathname, "image/" + fileType)
	}
}

func (ih *ImageHelper) UploadToS3(buf *bytes.Buffer, pathname string, fileFormat string) error {
	params := &s3.PutObjectInput{
		Bucket:      aws.String(ih.Config.GetAmazonS3Bucket()),
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

func (ih *ImageHelper) BatchDelete(key string) error {
	params := &s3.ListObjectsInput{
		Bucket: aws.String(ih.Config.GetAmazonS3Bucket()),
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
		for i := 0; i < len(objectIdentifiers); i += batchSize {
			if i + batchSize > len(objectIdentifiers) {
				batchDelete = &s3.Delete{
					Objects: objectIdentifiers,
				}
			} else {
				batchDelete = &s3.Delete{
					Objects: objectIdentifiers[i:i + batchSize],
				}
			}

			batchParams = &s3.DeleteObjectsInput{
				Bucket: aws.String(ih.Config.GetAmazonS3Bucket()),
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