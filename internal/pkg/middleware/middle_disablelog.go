package middleware

import "github.com/HYY-yu/seckill.sk/internal/pkg/core"

func (m *middleware) DisableLog() core.HandlerFunc {
	return func(c core.Context) {
		c.DisableLog(true)
	}
}