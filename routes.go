package main

import (
	"github.com/gin-gonic/gin"
)

func registerRankingRoutes(g *gin.RouterGroup) {
	g.GET("/updatemyscore", updateMyScore)
	g.GET("/getmyranking", getMyRanking)
}
