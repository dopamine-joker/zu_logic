package misc

import (
	"os"

	asr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/asr/v20190614"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

var CloudClient *asr.Client

func initTencentSDK() {
	credential := common.NewCredential(
		os.Getenv("SECRETID"),
		os.Getenv("SECRETKEY"),
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "asr.tencentcloudapi.com"
	CloudClient, _ = asr.NewClient(credential, "", cpf)
}
