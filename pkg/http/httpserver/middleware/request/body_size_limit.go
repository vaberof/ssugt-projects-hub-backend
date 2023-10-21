package request

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func BodySizeLimitMiddleware(maxBodySizeBytes int) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		contentLengthStr := ctx.Request.Header.Get("Content-Length")
		contentLength, err := strconv.Atoi(contentLengthStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"result":  "Error",
				"code":    "BAD_REQUEST",
				"message": fmt.Sprintf("Illegal header value. %s=%s", "Content-Length", contentLengthStr),
			})

			return
		}

		if contentLength > maxBodySizeBytes {
			ctx.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"result":  "Error",
				"code":    "REQUEST_ENTITY_TOO_LARGE",
				"message": fmt.Sprintf("Illegal request body size. Given request body size in bytes=%d, allowable max request body size in bytes=%d", contentLength, maxBodySizeBytes),
			})

			return
		}

		ctx.Next()
	}
}
