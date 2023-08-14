package toolx

import (
	"crypto/md5"
	"fmt"
	"github.com/google/uuid"
	"strings"
)

// toolx.go 自定义工具类

// GenUUID 生成随机的UUID
func GenUUID() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

// Md5Encrypt 密码加密
func Md5Encrypt(password string) string {
	hash := md5.Sum([]byte(password))
	return fmt.Sprintf("%x", hash)
}
