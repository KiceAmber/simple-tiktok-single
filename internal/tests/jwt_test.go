package tests

import (
	"simple_tiktok_rime/pkg/jwt"
	"testing"
)

func TestJWT(t *testing.T) {
	var userId int64 = 20001

	// 生成 token
	token, _ := jwt.GenToken(userId)
	t.Logf("token => %s", token)

	// 解析 token
	myClaims, err := jwt.ParseToken(token)
	if err != nil {
		t.Fatal("解析 token 失败")
		return
	}
	t.Log("myClaims.id => ", myClaims.Id)
}
