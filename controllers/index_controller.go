package controllers

import (
	"github.com/f-alotaibi/go-starter/views"
	"github.com/labstack/echo/v4"
)

type IndexController struct{}

func NewIndexController() *IndexController {
	return &IndexController{}
}

func (c *IndexController) Show(ctx echo.Context) error {
	return views.Index().Render(ctx.Request().Context(), ctx.Response())
}
