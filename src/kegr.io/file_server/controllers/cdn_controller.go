package controllers

import (
	"cerberus/storage_controller/config"
	"cerberus/storage_controller/model"
	"cerberus/storage_controller/service"
	"cerberus/storage_controller/util"
	"fmt"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// CdnController defines the routes that manipulate Cdns
type CdnController struct {
	c          *config.Config
	kegService service.IKegService
	Controller
}

// NewCdnController is the controller responsible for serving
// files to customers and defines those API endpoints
func NewCdnController(
	c *config.Config,
	ks service.IKegService) *CdnController {
	return &CdnController{
		c:          c,
		kegService: ks,
	}
}

// Register registers the necessary API endpoints this controller serves
func (cc *CdnController) Register(router *gin.Engine) {
	router.GET("/c/*path", cc.get)
}

func (cc *CdnController) get(ctx *gin.Context) {
	req := ctx.Param("path")
	kegAccessPath := filepath.Dir(req)[1:]
	liquidAccessName := filepath.Base(req)
	var liquid model.ILiquid

	var err error

	if liquid, err = cc.kegService.GetLiquidByPath(kegAccessPath, liquidAccessName); err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}
	content := liquid.GetContent()

	if liquid.GetOptions().GetCache() > 0 {
		ctx.Header("Cache-Control", fmt.Sprintf("public, max-age=%v", liquid.GetOptions().GetCache()))
	}

	if liquid.GetOptions().GetGzip() {
		ctx.Header("Content-Encoding", "gzip")
		content = util.GzipBytes(content)
	}

	mimeType := mime.TypeByExtension(liquid.GetOptions().GetExt())
	ctx.Data(http.StatusOK, mimeType, content)
}
