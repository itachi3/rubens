package controllers

import (
	"agentdesks"
	"agentdesks/datasources"
	"agentdesks/models"
	"agentdesks/utils"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type image_controller struct {
	router     *httprouter.Router
	uploadDs   *datasources.UploadDataSource
	downloadDs *datasources.DownloadDataSource
	deleteDs   *datasources.DeleteDataSource
}

func NewImageController(ru *httprouter.Router, conn *models.Connections, config *agentdesks.Config) *image_controller {
	return &image_controller{
		router:     ru,
		uploadDs:   datasources.NewUploadDataSource(conn, config),
		downloadDs: datasources.NewDownloadDataSource(conn, config),
		deleteDs:  	datasources.NewDeleteDataSource(conn, config),
	}
}

func (ic image_controller) InitializeHooks() {

	ic.router.POST("/v1.0/image", utils.ErrHandle(ic.uploadImage).ToHandle())

	ic.router.GET("/v1.0/imageUrls", utils.ErrHandle(ic.getImagesLocation).ToHandle())

	ic.router.DELETE("/v1.0/image", utils.ErrHandle(ic.deleteImages).ToHandle())

}

func (ic image_controller) uploadImage(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	if !utils.IsValidRequest(r, w) {
		return nil
	}

	return ic.uploadDs.UploadFile(w, r, p)
}

func (ic image_controller) getImagesLocation(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	if !utils.IsValidRequest(r, w) {
		return nil
	}

	return ic.downloadDs.GetImagesLocation(w, r)
}

func (ic image_controller) deleteImages(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	if !utils.IsValidRequest(r, w) {
		return nil
	}

	return ic.deleteDs.DeleteImages(w, r)
}
