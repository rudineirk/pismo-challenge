package httprouter

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rudineirk/pismo-challenge/pkg/infra/config"
)

var defaultTimeout = 60 * time.Second    //nolint:gochecknoglobals // default value
var defaultIdleTimeout = 5 * time.Minute //nolint:gochecknoglobals // default value

func StructuredLogger(logger *zerolog.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now() // Start timer
		path := ctx.Request.URL.Path
		raw := ctx.Request.URL.RawQuery

		// Process request
		ctx.Next()

		// Fill the params
		param := gin.LogFormatterParams{}

		param.TimeStamp = time.Now() // Stop timer
		param.Latency = param.TimeStamp.Sub(start)

		if param.Latency > time.Minute {
			param.Latency = param.Latency.Truncate(time.Second)
		}

		param.ClientIP = ctx.ClientIP()
		param.Method = ctx.Request.Method
		param.StatusCode = ctx.Writer.Status()
		param.ErrorMessage = ctx.Errors.ByType(gin.ErrorTypePrivate).String()
		param.BodySize = ctx.Writer.Size()

		if raw != "" {
			path = path + "?" + raw
		}

		param.Path = path

		// Log using the params
		var logEvent *zerolog.Event
		if ctx.Writer.Status() >= http.StatusInternalServerError {
			logEvent = logger.Error() //nolint:zerologlint // it's being used bellow
		} else {
			logEvent = logger.Info() //nolint:zerologlint // it's being used bellow
		}

		logEvent.Str("client_id", param.ClientIP).
			Str("method", param.Method).
			Int("status_code", param.StatusCode).
			Int("body_size", param.BodySize).
			Str("path", param.Path).
			Int64("latency_us", param.Latency.Microseconds()).
			Msg(param.ErrorMessage)
	}
}

func NewRouter(cfg *config.Config, logger *zerolog.Logger) *gin.Engine {
	if cfg.IsProduction {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(StructuredLogger(logger), gin.Recovery())
	_ = router.SetTrustedProxies([]string{})

	router.GET("/status", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"ok": true,
		})
	})

	return router
}

func NewServer(cfg *config.Config, router http.Handler) *http.Server {
	return &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.HTTPPort),
		Handler:           router,
		ReadTimeout:       defaultTimeout,
		ReadHeaderTimeout: defaultTimeout,
		WriteTimeout:      defaultTimeout,
		IdleTimeout:       defaultIdleTimeout,
	}
}
