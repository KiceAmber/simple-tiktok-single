package tests

import (
	"simple_tiktok_single/pkg/toolx"
	"testing"
)

func TestMd5Encrypt(t *testing.T) {
	password := toolx.Md5Encrypt("hello,world")
	t.Log("Password Encrypted => ", password)
}
