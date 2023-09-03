package controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/raman-vhd/arvan-challenge/internal/model"
	"github.com/raman-vhd/arvan-challenge/internal/service"
)

type IDataHandlerController interface {
	ProcessData(ctx echo.Context) error
}

type dataHandlerController struct {
	svc            service.IDataHandlerService
	rateLimiterSvc service.IRateLimitService
}

func NewDataHandler(
	svc service.IDataHandlerService,
	rateLimiterSvc service.IRateLimitService,
) IDataHandlerController {
	return dataHandlerController{
		svc:            svc,
		rateLimiterSvc: rateLimiterSvc,
	}
}

func (c dataHandlerController) ProcessData(ctx echo.Context) error {
	context := ctx.Request().Context()
	data := model.Data{
		ID:     ctx.Get("dataid").(string),
		UserID: ctx.Get("userid").(string),
	}

	err := c.svc.ProcessData(context, data)
	if err != nil {
		log.Println(err)
		return ctx.JSON(
			http.StatusInternalServerError,
			model.APIResponse{
				Msg:  "internal error. try again",
				Data: nil,
			})
	}

	err = c.rateLimiterSvc.UseRateLimitPerMinToken(context, data.UserID)
	if err != nil {
		log.Println(err)
		return ctx.JSON(
			http.StatusInternalServerError,
			model.APIResponse{
				Msg:  "internal error. try again",
				Data: nil,
			})
	}

	size := int(ctx.Request().ContentLength)
	err = c.rateLimiterSvc.AddToCurrentReqSizePerMon(context, data.UserID, size)
	if err != nil {
		log.Println(err)
		return ctx.JSON(
			http.StatusInternalServerError,
			model.APIResponse{
				Msg:  "internal error. try again",
				Data: nil,
			})
	}

	return ctx.JSON(
		http.StatusOK,
		model.APIResponse{
			Msg:  "Data with ID: " + fmt.Sprint(data.ID) + " has been successfully sent for processing",
			Data: data,
		})
}
