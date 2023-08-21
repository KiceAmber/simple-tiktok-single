package tests

import (
	"errors"
	"fmt"
	"simple_tiktok_rime/internal/model"
	"simple_tiktok_rime/internal/model/entity"
	"testing"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (*User) TableName() string {
	return "user"
}

func TestGormFindOne(t *testing.T) {

	dsn := "root:root@tcp(120.77.176.211:53306)/test_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal("connect DB error, ", err)
	}

	user := &User{
		Username: "tom",
	}

	result := db.First(user)
	if result.Error != nil {
		t.Fatal("query error: ", result.Error)
	}

	// 显示查询结果
	if result.RowsAffected > 0 {
		t.Logf("\nUser info: %#v\n", user)
	} else {
		t.Log("User not found")
	}
}

func TestGetUserInfo(t *testing.T) {
	dsn := "root:root@tcp(127.0.0.1:3306)/simple_tiktok?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal("connect DB error, ", err)
	}

	user := &entity.User{}
	result := db.Where("id = ?", 2).First(user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			t.Log("No Data Selected")
		}
		t.Log("Error: ", err)
	}

	t.Logf("\nuser => %#v\n", user)
}

// TestGormAssociatedQuery 测试 gorm 关联查询
func TestGormAssociatedQuery(t *testing.T) {
	dsn := "root:root@tcp(127.0.0.1:3306)/simple_tiktok?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal("connect DB error, ", err)
	}

	// 查询视频列表，并关联查询作者信息
	var videos = make([]*model.VideoItem, 0, 30)
	db.Preload("Author").Find(&videos)

	for _, video := range videos {
		fmt.Printf("Video ID: %d, Title: %s, Author: %s\n", video.Id, video.Title, video.Author.Name)
	}
}

func TestGormTransaction(t *testing.T) {
	dsn := "root:root@tcp(127.0.0.1:3306)/simple_tiktok?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal("connect DB error, ", err)
	}

	newComment := &entity.Comment{
		Id:        1001,
		Content:   "这是 milan 的评论，评论的是 tom 发布的视频",
		AuthorId:  81588023308849152,
		VideoId:   81588215466692608,
		CreatedAt: time.Now(),
	}

	// 使用事务来更新数据
	db.Transaction(func(tx *gorm.DB) error {
		// 插入评论数据
		if err := tx.Create(newComment).Error; err != nil {
			t.Fatal("tx.Create err =>", err)
			return err
		}
		// 更新视频表的评论数量
		video := &entity.Video{}
		if err = tx.Where("id = ?", newComment.VideoId).First(video).Error; err != nil {
			t.Fatal("tx.Where err =>", err)
			return err
		}
		if err := tx.Model(&entity.Video{}).Where("id = ?", newComment.VideoId).Update("comment_count", video.CommentCount+1).Error; err != nil {
			t.Fatal("tx.Model(\"video\").Where err =>", err)
			return err
		}
		return nil
	})
	t.Log("执行完毕")
}
