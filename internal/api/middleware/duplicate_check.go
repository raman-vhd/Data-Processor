package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/raman-vhd/arvan-challenge/internal/model"
	"github.com/raman-vhd/arvan-challenge/internal/service"
)

type DuplicateCheckMiddleware struct {
	svc service.IDuplicateCheckService
}

func NewDuplicateCheckMiddleware(
	svc service.IDuplicateCheckService,
) DuplicateCheckMiddleware {
	return DuplicateCheckMiddleware{
		svc: svc,
	}
}

func (m DuplicateCheckMiddleware) Handler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		data := model.Data{
			ID:     c.Get("dataid").(string),
			UserID: c.Get("userid").(string),
		}

		d, err := m.svc.IsDuplicate(ctx, data.UserID, data.ID)
		if err != nil {
			return c.JSON(
				http.StatusInternalServerError,
				model.APIResponse{
					Msg:  "internal server error. try again",
					Data: nil,
				},
			)
		}
		if d {
			return c.JSON(
				http.StatusConflict,
				model.APIResponse{
					Msg:  "data already exist",
					Data: nil,
				})
		}

		return next(c)
	}
}
