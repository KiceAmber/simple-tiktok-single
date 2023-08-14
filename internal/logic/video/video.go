package video

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"os/exec"
	"simple_tiktok_rime/internal/consts"
	"simple_tiktok_rime/internal/dao/mysql"
	"simple_tiktok_rime/internal/dao/redis"
	"simple_tiktok_rime/internal/model"
	"simple_tiktok_rime/internal/service"
	myqiniu "simple_tiktok_rime/pkg/qiniu"
	"simple_tiktok_rime/pkg/snowflake"
	"simple_tiktok_rime/pkg/toolx"
	"strings"

	"github.com/qiniu/go-sdk/v7/storage"
)

type sVideo struct{}

func init() {
	service.RegisterVideo(New())
}

func New() *sVideo {
	return &sVideo{}
}

// PublishVideo 发布视频投稿
// 1. 检测文件扩展名是否合法，如 mp4、flv 等视频格式
// 1. 使用雪花算法生成视频ID `id`
// 2. 将视频存储到第三方 OSS，得到视频的播放链接 `play_url`
// 3. 使用 ffmpeg 从视频解析出封面 `cover_url`
// 4. 数据库入库操作
func (*sVideo) PublishVideo(in *model.PublishVideoInput) (err error) {

	// 检测文件扩展名是否合法，不合法直接返回
	filename := in.VideoData.Filename
	if !isVideoExtValid(filename) {
		err = consts.ErrInvalidVideoExt
		return
	}

	// 生成雪花算法
	in.Id = snowflake.GenID()

	// 上传到第三方 OSS 七牛云，得到视频的播放链接
	if in.PlayUrl, err = uploadVideo(in.VideoData); err != nil {
		return
	}

	// 视频截图生成封面图上传到 OSS
	if in.CoverUrl, err = uploadCover(in.PlayUrl, "1"); err != nil {
		return
	}

	// 数据库操作
	err = mysql.Video().InsertVideoInfo(in)
	return
}

// isVideoExtValid 检测文件扩展名是否合法
func isVideoExtValid(filename string) bool {
	validVideoExt := []string{"mp4", "flv", "avi"}

	// 文件后缀名的判断
	indexOfDot := strings.LastIndex(filename, ".")
	if indexOfDot < 0 {
		return false
	}

	suffix := strings.ToLower(filename[indexOfDot+1:])
	for _, fileExt := range validVideoExt {
		if suffix == fileExt {
			return true
		}
	}
	return false
}

// uploadVideo 上传视频到第三方 OSS
// 通过 videoFile 视频文件的元数据，将视频上传到 OSS
// 返回的 string 为视频的播放链接
func uploadVideo(videoFile *multipart.FileHeader) (string, error) {

	srcFile, err := videoFile.Open()
	if err != nil {
		return "", err
	}
	defer srcFile.Close()

	// 获取到视频后缀
	suffix := strings.ToLower(videoFile.Filename[strings.LastIndex(videoFile.Filename, ".")+1:])

	// 七牛云的 GOSDK 上传操作
	qiniu := myqiniu.New()
	ret := storage.PutRet{}
	err = qiniu.FormUploader.Put(
		context.Background(),
		&ret,
		qiniu.UploadToken,
		qiniu.FolderName.Video+"/"+toolx.GenUUID()+"."+suffix,
		srcFile,
		videoFile.Size,
		&storage.PutExtra{},
	)
	if err != nil {
		return "", err
	}

	playUrl := fmt.Sprintf("http://%s/%s", qiniu.HostName, ret.Key)
	return playUrl, nil
}

// parseCoverToLocal 解析视频的封面到本地
// playUrl 为目标视频的播放链接，time 为截取视频第 time 秒的帧
func uploadCover(playUrl string, time string) (string, error) {
	// ffmpeg 解析视频获取截图
	cmd := exec.Command("ffmpeg", "-i", playUrl, "-ss", time, "-frames:v", "1", "./output.jpg")
	_, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	// 本地图片会有残留，在该函数最后需要调用 os.Remove() 来删除本地的图片
	defer func() {
		_ = os.Remove("./output.jpg")
	}()

	// 将本地生成的 output.jpg 上传到 OSS
	qiniu := myqiniu.New()
	ret := storage.PutRet{}
	if err = qiniu.FormUploader.PutFile(
		context.Background(),
		&ret,
		qiniu.UploadToken,
		qiniu.FolderName.Cover+"/"+toolx.GenUUID()+".jpg",
		"./output.jpg",
		nil,
	); err != nil {
		return "", err
	}

	coverUrl := fmt.Sprintf("http://%s/%s", qiniu.HostName, ret.Key)
	return coverUrl, nil
}

func (*sVideo) GetVideoFeed(in *model.VideoFeedInput) (out *model.VideoFeedOutput, err error) {
	out = new(model.VideoFeedOutput)

	// 首先获取到视频流
	out, err = mysql.Video().GetVideoList(in)
	if err != nil {
		return nil, err
	}
	// 将 NextTime 设置为视频列表中发布最早的视频时间
	out.NextTime = out.VideoList[len(out.VideoList)-1].CreatedAt

	// 查询 redis，查看用户是否点赞视频，如果 in.UserId == -1 则说明为未登录用户，点赞默认为 false
	if in.UserId != -1 {
		err = redis.Video().IsUserFavoriteVideo(in.UserId, out.VideoList)
		if err != nil {
			return nil, err
		}
	}

	for _, video := range out.VideoList {
		video.Author = &model.UserItem{}
		userOut, err := mysql.User().QueryUserById(&model.GetUserInfoInput{UserId: video.AuthorId})
		if err != nil {
			fmt.Println("AuthorId = ", userOut.UserItem.Id, " 的作者信息查询失败")
			continue
		}
		video.Author = userOut.UserItem
	}
	return
}

// GetVideoPublishedList 获取视频发布列表
func (*sVideo) GetVideoPublishedList(in *model.GetVideoPublishedListInput) (out *model.GetVideoPublishedListOutput, err error) {

	// 因为视频的作者都是一致的，所以只需查出一次，直接赋给结构体的字段即可
	userOut, err := mysql.User().QueryUserById(&model.GetUserInfoInput{UserId: in.UserId})
	if err != nil {
		return nil, err
	}

	out = new(model.GetVideoPublishedListOutput)

	// 首先获取到视频流
	out, err = mysql.Video().GetVideoPublishedList(in)
	if err != nil {
		return nil, err
	}

	for _, video := range out.VideoList {
		video.Author = userOut.UserItem
	}
	return
}
