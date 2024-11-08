package middleware

import (
	"api-gateway/internal/model"
	"api-gateway/internal/services"
	"bytes"
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
)

type TrafficMiddleware struct {
	TrafficService services.TrafficService
}

func (tm *TrafficMiddleware) TrafficStatsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录入站流量
		var inTraffic bytes.Buffer
		teeReader := io.TeeReader(c.Request.Body, &inTraffic)
		c.Request.Body = io.NopCloser(teeReader)

		// 代理请求（这里假设你有代理请求的逻辑）
		//...

		// 记录出站流量
		var outTraffic bytes.Buffer
		c.Writer = &responseWriter{
			ResponseWriter: c.Writer,
			buffer:         &outTraffic,
		}

		c.Next()

		// 计算流量大小并记录流量统计信息
		apiName := c.Request.URL.Path
		inSize := inTraffic.Len()
		outSize := outTraffic.Len()
		stats := model.TrafficStats{
			API:        apiName,
			InTraffic:  int64(inSize),
			OutTraffic: int64(outSize),
		}
		ctx := c.Request.Context()
		err := tm.TrafficService.RecordTrafficStats(ctx, &stats)
		if err != nil {
			fmt.Println("Error recording traffic stats:", err)
		}
	}
}

// 自定义的ResponseWriter用于捕获出站流量
type responseWriter struct {
	gin.ResponseWriter
	buffer *bytes.Buffer
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.buffer.Write(b)
	return rw.ResponseWriter.Write(b)
}
