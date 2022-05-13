package api

import (
	"github.com/HYY-yu/seckill.pkg/core"
)

func (s *Server) Route(c *Handlers) {

	v1Group := s.Engine.Group("/v1")
	{
		skGroup := v1Group.Group("/")
		skGroup.Use(core.WrapAuthHandler(s.Middles.Jwt), c.loginHandler.CheckBlackList)
		skGroup.GET("/list", c.skHandler.List)
		skGroup.PUT("/resource", c.skHandler.Add)
		skGroup.DELETE("/resource", c.skHandler.Delete)
		skGroup.GET("/unlogin", c.loginHandler.Unlogin)

		orderGroup := skGroup.Group("/order")
		orderGroup.GET("/list", c.orderHandler.List)
		orderGroup.GET("/join", c.orderHandler.Join)
	}

	{
		v1Group.GET("/login", c.loginHandler.Login)
		v1Group.POST("/refresh", c.loginHandler.RefreshToken)
	}
}
