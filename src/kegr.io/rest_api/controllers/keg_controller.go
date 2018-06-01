package controllers

import (
	"context"
	"net/http"

	"kegr.io/protobuf/model/storage/keg"
	"kegr.io/protobuf/server/storage"

	"github.com/gin-gonic/gin"
	"kegr.io/storage_client"
)

// KegController is the controller responsible for the Keg endpoints
type KegController struct {
	c storage_client.IClient
	IController
}

// NewKegController creates a new instance of the KegController struct
func NewKegController(c storage_client.IClient) *KegController {
	return &KegController{
		c: c,
	}
}

// Register registers the necessary API endpoints this controller serves
func (kc *KegController) Register(router *gin.Engine) {
	group := router.Group("/keg")
	{
		group.POST("/", kc.create)
		group.GET("/", kc.getKegs)
		group.GET("/:kegID", kc.get)
		group.GET("/:kegID/liquid", kc.getLiquids)
		group.PUT("/:kegID", kc.update)
		group.DELETE("/:kegID", kc.delete)
		// group.GET("/:kegID/liquids", kc.getLiquids)
	}
}

func (kc *KegController) create(ctx *gin.Context) {
	kegID := ctx.Param("kegID")
	var options *keg.Options
	ctx.BindJSON(&options)

	res, err := kc.c.Get().CreateKeg(
		context.Background(),
		&storage.CreateKegRequest{
			KegId:   kegID,
			Options: options,
		},
	)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	ctx.String(http.StatusOK, res.GetKegID())
}

func (kc *KegController) get(ctx *gin.Context) {
	res, err := kc.c.Get().GetKeg(
		context.Background(),
		&storage.GetKegRequest{
			KegId: ctx.Param("kegID"),
		},
	)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, res.Options)
}

func (kc *KegController) getKegs(ctx *gin.Context) {
	res, err := kc.c.Get().GetKegs(
		context.Background(),
		&storage.GetKegsRequest{},
	)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, res.Kegs)
}

func (kc *KegController) update(ctx *gin.Context) {
	var options *keg.Options
	ctx.BindJSON(&options)

	_, err := kc.c.Get().UpdateKegOptions(
		context.Background(),
		&storage.UpdateKegOptionsRequest{
			KegId:   ctx.Param("kegID"),
			Options: options,
		},
	)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	ctx.Status(http.StatusOK)
}

func (kc *KegController) delete(ctx *gin.Context) {
	_, err := kc.c.Get().DeleteKeg(
		context.Background(),
		&storage.DeleteKegRequest{
			KegId: ctx.Param("kegID"),
		},
	)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	ctx.Status(http.StatusOK)
}

func (kc *KegController) getLiquids(ctx *gin.Context) {
	res, err := kc.c.Get().GetKegLiquids(
		context.Background(),
		&storage.GetKegLiquidsRequest{
			KegId: ctx.Param("kegID"),
		},
	)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, res.Liquids)
}
