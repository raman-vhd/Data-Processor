package middleware

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/raman-vhd/arvan-challenge/internal/model"
	"github.com/raman-vhd/arvan-challenge/internal/service"
)

type RateLimitMiddleware struct {
	svc service.IRateLimitService
}

func NewRateLimitMiddleware(svc service.IRateLimitService) RateLimitMiddleware {
	return RateLimitMiddleware{
		svc: svc,
	}
}

func (m RateLimitMiddleware) Handler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		var data model.Data
		err := c.Bind(&data)
		if err != nil {
			log.Println(err)
			return c.JSON(
				http.StatusUnprocessableEntity,
				model.APIResponse{
					Msg:  "failed to get data",
					Data: nil,
				})
		}

		c.Set("userid", data.UserID)
		c.Set("dataid", data.ID)

		ok, err := m.svc.CheckRateLimitPerMin(ctx, data.UserID)
		if err != nil {
			log.Println(err)
			return c.JSON(
				http.StatusInternalServerError,
				model.APIResponse{
					Msg:  "internal server error. try again",
					Data: nil,
				},
			)
		}
		if !ok {
			return c.JSON(
				http.StatusTooManyRequests,
				model.APIResponse{
					Msg:  "too many requests. exceed rate limit per minute",
					Data: nil,
				})
		}

		size := int(c.Request().ContentLength)
		ok, err = m.svc.CheckReqSizeLimitPerMon(ctx, data.UserID, size)
		if err != nil {
			log.Println(err)
			return c.JSON(
				http.StatusInternalServerError,
				model.APIResponse{
					Msg:  "internal server error. try again",
					Data: nil,
				},
			)
		}
		if !ok {
			return c.JSON(
				http.StatusRequestEntityTooLarge,
				model.APIResponse{
					Msg:  "too many requests. exceed request size limit per month",
					Data: nil,
				})
		}

		return next(c)
	}
}
