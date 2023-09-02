package route

import (
	"github.com/raman-vhd/arvan-challenge/internal/api/controller"
	"github.com/raman-vhd/arvan-challenge/internal/lib"
)

type dataHandlerRoute struct {
	handler         lib.RequestHandler
	ctrl            controller.IDataHandlerController
}

func NewDataHandler(
	handler lib.RequestHandler,
	ctrl controller.IDataHandlerController,
) dataHandlerRoute {
	return dataHandlerRoute{
		handler:         handler,
		ctrl:            ctrl,
	}
}

func (a dataHandlerRoute) Setup() {
	api := a.handler.Echo.Group("/api")

	api.POST("/data", a.ctrl.ProcessData)
}
