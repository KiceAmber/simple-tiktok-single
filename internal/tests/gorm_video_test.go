package tests

import (
	"simple_tiktok_rime/internal/model/entity"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// TestVideoCreate 测试插入 video 数据
func TestVideoCreate(t *testing.T) {
	dsn := "root:root@tcp(120.77.176.211:53306)/simple_tiktok?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal("connect DB error, ", err)
	}

	var newVideo = new(entity.Video)
	newVideo.Id = 10001
	newVideo.AuthorId = 20001
	newVideo.Title = "这是测试视频标题"
	newVideo.PlayUrl = "这是测试视频播放地址"
	newVideo.CoverUrl = "这是测试视频封面地址"

	result := db.Create(newVideo)
	if result.Error != nil {
		t.Fatal("query error: ", result.Error)
	}
}
