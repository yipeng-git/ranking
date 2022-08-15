package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

const port = ":8080"

func main() {
	initRedis()
	// initLog()
	// initMySQL()

	g := gin.New()
	g.Use(log())

	// g.Use(middlewares)

	rankingGroup := g.Group("ranking")
	registerRankingRoutes(rankingGroup)

	server := &http.Server{Addr: port, Handler: g}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func log() gin.HandlerFunc {
	return gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: func(params gin.LogFormatterParams) string {
			logFormat := map[string]interface{}{
				"response_time": params.TimeStamp.Format("2006/01/02 - 15:04:05"),
				"http_code":     params.StatusCode,
				"latency_time":  params.Latency,
				"client_ip":     params.Latency,
				"method":        params.Method,
				"path":          params.Path,
				"error_message": params.ErrorMessage,
			}
			b, _ := json.Marshal(logFormat)
			str := string(b)
			fmt.Println(str)

			return str
		},
		Output: os.Stdin,
		SkipPaths: []string{
			"/ping",
		},
	})
}
