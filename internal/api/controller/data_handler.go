package controller

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/raman-vhd/arvan-challenge/internal/model"
	"github.com/raman-vhd/arvan-challenge/internal/service"
)

type IDataHandlerController interface {
	ProcessData(ctx echo.Context) error
}

type dataHandlerController struct {
	svc service.IDataHandlerService
}

func NewDataHandler(
	svc service.IDataHandlerService,
) IDataHandlerController {
	return dataHandlerController{
		svc: svc,
	}
}

func (c dataHandlerController) ProcessData(ctx echo.Context) error {
    context := ctx.Request().Context()
	var data model.Data

	err := ctx.Bind(&data)
	if err != nil {
		return ctx.JSON(
			http.StatusUnprocessableEntity,
			model.APIResponse{
				Msg:  "failed to get dataHandler",
				Data: nil,
			})
	}

	err = c.svc.ProcessData(context, data)
	if err != nil {
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
