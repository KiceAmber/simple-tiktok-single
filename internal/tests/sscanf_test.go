package tests

import (
	"fmt"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"testing"
)

func TestSscanf(t *testing.T) {
	var arr = []string{
		"user:1002:favorite_video",
		"user:1001:favorite_video",
		"user:1003:favorite_video",
	}
	for _, key := range arr {
		//var userId int64 = 0
		split := strings.Split(key, ":")
		userId, err := strconv.ParseInt(split[1], 10, 64)
		if err != nil {
			zap.L().Error("ParseInt Failed", zap.Error(err))
			continue
		}
		fmt.Println("userId =>", userId)
	}
}
