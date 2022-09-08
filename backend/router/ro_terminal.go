package router

import (
	v1 "github.com/1Panel-dev/1Panel/app/api/v1"
	"github.com/1Panel-dev/1Panel/middleware"

	"github.com/gin-gonic/gin"
)

type TerminalRouter struct{}

func (s *TerminalRouter) InitTerminalRouter(Router *gin.RouterGroup) {
	terminalRouter := Router.Group("terminals").Use(middleware.JwtAuth()).Use(middleware.SessionAuth())
	baseApi := v1.ApiGroupApp.BaseApi
	{
		terminalRouter.GET("", baseApi.WsSsh)
		terminalRouter.GET("/local", baseApi.LocalWsSsh)
	}
}