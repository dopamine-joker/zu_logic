package misc

import (
	"github.com/tencentyun/cos-go-sdk-v5"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"os"
)

var CosClient *cos.Client

//initCos 初始化cos
func initCos() {
	secretID := os.Getenv("SECRETID")
	secretKey := os.Getenv("SECRETKEY")
	rawUrl := os.Getenv("RAWURL")
	Logger.Info("secretId", zap.String("id", secretID))
	Logger.Info("secretKey", zap.String("key", secretKey))
	Logger.Info("rawUrl", zap.String("rawurl", rawUrl))
	u, err := url.Parse(rawUrl)
	if err != nil {
		panic(err)
	}
	b := &cos.BaseURL{BucketURL: u}
	CosClient = cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  secretID,
			SecretKey: secretKey,
		},
	})
}
