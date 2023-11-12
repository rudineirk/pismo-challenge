package healthcheck

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

func SetupHealthCheck(router *gin.Engine, db *bun.DB) {
	dbHealthCheck := func(ctx *gin.Context) {
		var result int
		err := db.QueryRowContext(ctx, "SELECT 1").Scan(&result)

		if err != nil || result != 1 {
			ctx.JSON(http.StatusServiceUnavailable, gin.H{
				"ok": false,
			})
		} else {
			ctx.JSON(http.StatusOK, gin.H{
				"ok": true,
			})
		}
	}

	router.GET("/status", dbHealthCheck)
	router.GET("/healthcheck/readiness", dbHealthCheck)
	router.GET("/healthcheck/liveliness", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"ok": true,
		})
	})
}
