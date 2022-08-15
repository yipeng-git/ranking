package main

import (
	"math"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

const listSize = 10
const tsEndpoint = 2000000000
const rankingZset = "ranking:zset"

var scoreBase = math.Pow10(10)

func updateMyScore(ctx *gin.Context) {
	userId := ctx.Request.FormValue("userid")
	scoreStr := ctx.Request.FormValue("score")
	score, err := strconv.ParseFloat(scoreStr, 64)
	if err != nil {
		responseErr(ctx, "wrong score format")
		return
	}

	rs := newRankingScore(userId, score)
	err = rs.Update(ctx)
	if err != nil {
		responseErr(ctx, err.Error())
		return
	}
	responseOk(ctx, struct{}{})
}

type rankingStruct struct {
	UserId string  `json:"userId"`
	Score  float64 `json:"score"`
}

type getMyRankingResp struct {
	MyRank rankingStruct   `json:"myRank"`
	Before []rankingStruct `json:"before"`
	After  []rankingStruct `json:"after"`
}

func getMyRanking(ctx *gin.Context) {
	userId := ctx.Request.FormValue("userid")

	rs := newRankingScore(userId, 0)
	res, err := rs.GetMyRank(ctx)
	if err != nil {
		responseErr(ctx, err.Error())
	} else {
		responseOk(ctx, res)
	}
}

// =========== 业务逻辑 ===========

type rankingScore struct {
	userId string
	score  float64
}

func newRankingScore(userId string, score float64) *rankingScore {
	return &rankingScore{userId: userId, score: score}
}

// Update 更新用户在排行榜中的得分
func (rs *rankingScore) Update(ctx *gin.Context) error {
	currScore, err := redisCli.ZScore(ctx, rankingZset, rs.userId).Result()
	if err != nil && err != redis.Nil {
		return err
	}
	currScore = rs.fromRedisScore(currScore)
	if rs.score > currScore {
		return redisCli.ZAdd(ctx, rankingZset, &redis.Z{Score: rs.toRedisScore(), Member: rs.userId}).Err()
	}
	return nil
}

func (rs *rankingScore) toRedisScore() float64 {
	ts := tsEndpoint - time.Now().Unix()
	return rs.score*scoreBase + float64(ts)
}

// GetMyRank 获取用户在排行榜中的排名、得分，以及前后各10名的排名、得分
func (rs *rankingScore) GetMyRank(ctx *gin.Context) (*getMyRankingResp, error) {
	myRank, err := redisCli.ZRank(ctx, rankingZset, rs.userId).Result()
	if err == redis.Nil {
		return &getMyRankingResp{
			MyRank: rankingStruct{
				UserId: rs.userId,
				Score:  0,
			},
		}, nil
	}
	start := myRank - 10
	if start < 0 {
		start = 0
	}
	end := myRank + 10
	zRes, err := redisCli.ZRevRangeWithScores(ctx, rankingZset, start, end).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	res := &getMyRankingResp{
		Before: make([]rankingStruct, 0, listSize),
		After:  make([]rankingStruct, 0, listSize),
	}
	before := true
	for i, _ := range zRes {
		tmpScore := zRes[i].Score
		tmpUserId := zRes[i].Member.(string)
		if tmpUserId == rs.userId {
			res.MyRank = rankingStruct{UserId: tmpUserId, Score: rs.fromRedisScore(tmpScore)}
			before = false
		} else if before {
			res.Before = append(res.Before, rankingStruct{UserId: tmpUserId, Score: rs.fromRedisScore(tmpScore)})
		} else {
			res.After = append(res.After, rankingStruct{UserId: tmpUserId, Score: rs.fromRedisScore(tmpScore)})
		}
	}
	return res, nil
}

func (rs *rankingScore) fromRedisScore(redisScore float64) float64 {
	return math.Floor(redisScore / scoreBase)
}
