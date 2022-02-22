package api

import (
	"github.com/HYY-yu/seckill.sk/internal/pkg/core"
)

func (s *Server) Route(c *Handlers) {

	v1Group := s.Engine.Group("/v1")
	{
		v1Group.Use(core.WrapAuthHandler(s.Middles.Jwt))

	}
}
