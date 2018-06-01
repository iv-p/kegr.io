package controllers

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
	"kegr.io/protobuf/model/storage/liquid"
	"kegr.io/protobuf/server/storage"
	"kegr.io/storage_client"
	"kegr.io/storage_controller/util"
)

// LiquidController is the controller responsible for the Liquid endpoints
type LiquidController struct {
	c storage_client.IClient
	IController
}

// NewLiquidController creates a new instance of the LiquidController struct
func NewLiquidController(c storage_client.IClient) *LiquidController {
	return &LiquidController{
		c: c,
	}
}

// Register registers the necessary API endpoints this controller serves
func (fc *LiquidController) Register(router *gin.Engine) {
	group := router.Group("/liquid")
	{
		group.POST("/", fc.create)
		group.PUT("/:liquidId/keg/:kegId", fc.update)
		group.DELETE("/:liquidId/keg/:kegId", fc.delete)
	}
}

func (fc *LiquidController) update(ctx *gin.Context) {
	var options *liquid.Options
	ctx.BindJSON(&options)

	_, err := fc.c.Get().UpdateLiquidOptions(
		context.Background(),
		&storage.UpdateLiquidOptionsRequest{
			KegId:    ctx.Param("kegId"),
			LiquidId: ctx.Param("liquidId"),
			Options:  options,
		},
	)

	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	ctx.Status(http.StatusCreated)
}

func (fc *LiquidController) delete(ctx *gin.Context) {
	_, err := fc.c.Get().DeleteLiquid(
		context.Background(),
		&storage.DeleteLiquidRequest{
			KegId:    ctx.Param("kegId"),
			LiquidId: ctx.Param("liquidId"),
		},
	)

	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	ctx.Status(http.StatusOK)
}

func (fc *LiquidController) create(ctx *gin.Context) {
	info, _ := ctx.FormFile("file")

	file, _ := info.Open()
	defer file.Close()

	ext := path.Ext(info.Filename)[1:]
	name := strings.TrimSuffix(info.Filename, ext)
	sz := len(name)
	name = strings.Replace(name, " ", "_", -1)[:sz-1]

	content := bytes.NewBuffer(nil)
	_, err := io.Copy(content, file)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	hash := util.GetFileHash(file)

	size := info.Size
	cache := int64(120)
	gzip := false

	options := &liquid.Options{}
	options.Name = name
	options.Ext = ext
	options.Cache = cache
	options.Gzip = gzip

	liquid := &liquid.Liquid{}
	liquid.FileHash = hash
	liquid.Size = size
	liquid.Content = content.Bytes()
	liquid.Options = options
	liquid.FileHash = hash

	_, err = fc.c.Get().CreateLiquid(
		context.Background(),
		&storage.CreateLiquidRequest{
			KegId:  ctx.PostForm("kegID"),
			Liquid: liquid,
		},
	)

	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx.Status(http.StatusCreated)
}
