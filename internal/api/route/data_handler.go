package route

import (
	"github.com/raman-vhd/arvan-challenge/internal/api/controller"
	"github.com/raman-vhd/arvan-challenge/internal/api/middleware"
	"github.com/raman-vhd/arvan-challenge/internal/lib"
)

type dataHandlerRoute struct {
	handler         lib.RequestHandler
	ctrl            controller.IDataHandlerController
	rateLimitM      middleware.RateLimitMiddleware
	duplicateCheckM middleware.DuplicateCheckMiddleware
}

func NewDataHandler(
	handler lib.RequestHandler,
	ctrl controller.IDataHandlerController,
	rateLimitM middleware.RateLimitMiddleware,
	duplicateCheckM middleware.DuplicateCheckMiddleware,
) dataHandlerRoute {
	return dataHandlerRoute{
		handler:         handler,
		ctrl:            ctrl,
		rateLimitM:      rateLimitM,
		duplicateCheckM: duplicateCheckM,
	}
}

func (a dataHandlerRoute) Setup() {
	api := a.handler.Echo.Group("/api")
	api.Use(
		a.rateLimitM.Handler,
		a.duplicateCheckM.Handler,
	)

	api.POST("/data", a.ctrl.ProcessData)
}
