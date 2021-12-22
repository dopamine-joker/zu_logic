package misc

import (
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"net/url"
	"os"
)

var CosClient *cos.Client

//InitCos 初始化cos
func InitCos() {
	u, err := url.Parse(os.Getenv("RAWURL"))
	if err != nil {
		panic(err)
	}
	b := &cos.BaseURL{BucketURL: u}
	CosClient = cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv("SECRETID"),
			SecretKey: os.Getenv("SECRETKEY"),
		},
	})
}
