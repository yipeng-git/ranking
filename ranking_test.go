package main

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestRankingScore_toRedisScore(t *testing.T) {
	rs := newRankingScore("", 123)
	ts := 2000000000 - time.Now().Unix()
	if rs.toRedisScore() != float64(1230000000000+ts) {
		t.Errorf("wrong rs.toRedisScore, expected 1230000000000, got %v", rs.toRedisScore())
	}
}

func TestRankingScore_fromRedisScore(t *testing.T) {
	rs := newRankingScore("", 123)
	if rs.fromRedisScore(1230000000000) != 123 {
		t.Errorf("wrong rs.toRedisScore, expected 1230000000000, got %v", rs.fromRedisScore(1230000000000))
	}
}

func TestRanking(t *testing.T) {
	initRedis()
	rs := rankingScore{userId: "aaaaa"}
	myRank, err := redisCli.ZRank(context.Background(), rankingZset, rs.userId).Result()
	fmt.Println(myRank, err)
}
