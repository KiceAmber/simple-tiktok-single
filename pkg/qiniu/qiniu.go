package qiniu

import (
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"simple_tiktok_rime/manifest/config"
)

type Qiniu struct {
	UploadToken  string
	FormUploader *storage.FormUploader
	HostName     string
	FolderName   *struct {
		Video string
		Cover string
	}
}

var qiniu *Qiniu

func New() *Qiniu {
	if qiniu == nil {
		Init(config.Conf.QiniuConfig)
	}
	return qiniu
}

// Init 初始化七牛云配置
func Init(configFile *config.QiniuConfig) {
	// 生成上传凭证 UploadToken 和 上传器 FormUploader
	putPolicy := storage.PutPolicy{
		Scope: configFile.Bucket,
	}
	mac := qbox.NewMac(configFile.AccessKey, configFile.SecretKey)
	qiniu = &Qiniu{
		UploadToken: putPolicy.UploadToken(mac),
		FormUploader: storage.NewFormUploader(&storage.Config{
			Zone: &storage.ZoneHuanan,
		}),
		HostName: configFile.HostName,
		FolderName: &struct {
			Video string
			Cover string
		}{Video: configFile.FolderName.Video, Cover: configFile.FolderName.Cover},
	}
}
