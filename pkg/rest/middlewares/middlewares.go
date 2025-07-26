package rest_middlewares

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/krasov-rf/goback/pkg/rps"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Domain, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func UnescapeQueryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		strUrl, err := url.QueryUnescape(c.Request.URL.String())
		if err != nil {
			return
		}
		u, err := url.Parse(strUrl)
		if err != nil {
			return
		}
		c.Request.URL = u
	}
}

func HostMiddleware(
	hostAllowed map[string]bool,
	logger Logger,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		host := strings.Split(c.Request.Host, ":")
		if len(host) == 0 {
			c.AbortWithStatus(http.StatusBadGateway)
			return
		}
		h := strings.Trim(strings.ToLower(host[0]), " ")

		if ex, ok := hostAllowed[h]; !ok || !ex {
			c.AbortWithStatus(http.StatusBadGateway)
			msg := fmt.Sprintf("host not allowed %s", host)
			if logger != nil {
				logger.Error(errors.New(msg))
			}
			return
		}
		c.Next()
	}
}

type Logger interface {
	Error(any)
}

// ограничение обращений в минуту
func RpsMiddleware[T rps.IRps, L Logger](
	rps T,
	logger L,
	service_name string,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := rps.Wait(c, c.ClientIP()); err != nil {
			logger.Error(err)
			return
		}
		c.Next()
	}
}
