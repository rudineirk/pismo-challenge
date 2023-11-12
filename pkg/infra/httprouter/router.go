package httprouter

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

var defaultTimeout = 60 * time.Second    //nolint:gochecknoglobals // default value
var defaultIdleTimeout = 5 * time.Minute //nolint:gochecknoglobals // default value

func StructuredLogger(logger *zerolog.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now() // Start timer
		path := ctx.Request.URL.Path

		rawQuery := ctx.Request.URL.RawQuery
		if rawQuery != "" {
			path = path + "?" + rawQuery
		}

		ctx.Next()

		var logEvent *zerolog.Event
		if ctx.Writer.Status() >= http.StatusInternalServerError {
			logEvent = logger.Error() //nolint:zerologlint // it's being used bellow
		} else {
			logEvent = logger.Info() //nolint:zerologlint // it's being used bellow
		}

		logEvent.Str("client_id", ctx.ClientIP()).
			Str("method", ctx.Request.Method).
			Int("status_code", ctx.Writer.Status()).
			Int("body_size", ctx.Writer.Size()).
			Str("path", path).
			Int64("latency_us", time.Since(start).Microseconds()).
			Msg(ctx.Errors.ByType(gin.ErrorTypePrivate).String())
	}
}

func NewRouter(logger *zerolog.Logger, isProduction bool) *gin.Engine {
	if isProduction {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(StructuredLogger(logger), gin.Recovery())
	_ = router.SetTrustedProxies([]string{})

	return router
}

func NewServer(httpPort int, router http.Handler) *http.Server {
	return &http.Server{
		Addr:              fmt.Sprintf(":%d", httpPort),
		Handler:           router,
		ReadTimeout:       defaultTimeout,
		ReadHeaderTimeout: defaultTimeout,
		WriteTimeout:      defaultTimeout,
		IdleTimeout:       defaultIdleTimeout,
	}
}
