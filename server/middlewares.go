package server

import (
	"fmt"
	"runtime"
	"time"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func loggerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req, res := c.Request(), c.Response()

			start := time.Now()
			if err := next(c); err != nil {
				c.Error(err)
			}
			end := time.Now()

			p := req.URL.Path
			if p == "" {
				p = "/"
			}

			bytesIn := req.Header.Get(echo.HeaderContentLength)
			if bytesIn == "" {
				bytesIn = "0"
			}

			headers := fmt.Sprintf("%+v", req.Header)
			entry := log.WithFields(log.Fields{
				"start":     start.Format(time.RFC3339),
				"end":       end.Format(time.RFC3339),
				"remote-ip": c.RealIP(),
				"host":      req.Host,
				"uri":       req.RequestURI,
				"method":    req.Method,
				"path":      p,
				"headers":   headers,
				"status":    res.Status,
				"latency":   end.Sub(start).String(),
				"bytes-in":  bytesIn,
				"bytes-out": res.Size,
			})

			switch {
			case res.Status < 400:
				entry.Info("Handled request")
			case res.Status < 500:
				entry.Warn("Handled request")
			default:
				entry.Error("Handled request")
			}

			return nil
		}
	}
}

// Same as echo's RecoverWithConfig middleware, with DefaultRecoverConfig
func recoverMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}
					stack := make([]byte, 4<<10) // 4 KB
					length := runtime.Stack(stack, true)
					log.WithError(err).Errorf("[PANIC RECOVER] %s", stack[:length])
					c.Error(err)
				}
			}()
			return next(c)
		}
	}
}
