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

// ConvertUnit 转化单位，将过千的数据转化为以万为单位，并以 string 的形式返回
func ConvertUnit(number int64) string {
	if number >= 10000 {
		unit := float64(number) / 10000.0
		return fmt.Sprintf("%.2fw", unit)
	} else {
		return fmt.Sprintf("%d", number)
	}
}
