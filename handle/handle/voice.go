package handle

import (
	"bytes"
	"errors"
	"github.com/dopamine-joker/zu_logic/misc"
	"go.uber.org/zap"
	"time"

	asr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/asr/v20190614"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	terrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

//processVoice 通过sdk调用接口获取语音识别结果
func processVoice(cosUrl string) (txt string, err error) {

	request := asr.NewCreateRecTaskRequest()

	request.EngineModelType = common.StringPtr("16k_zh")
	request.ChannelNum = common.Uint64Ptr(1)
	request.ResTextFormat = common.Uint64Ptr(2)
	request.SourceType = common.Uint64Ptr(0)
	request.Url = common.StringPtr(cosUrl)

	response, err := misc.CloudClient.CreateRecTask(request)
	if _, ok := err.(*terrors.TencentCloudSDKError); ok {
		misc.Logger.Error("An API error has returned: %s", zap.Error(err))
		return
	}
	if err != nil {
		misc.Logger.Error("An API error has returned: %s", zap.Error(err))
		return "", err
	}

	taskId := *response.Response.Data.TaskId

	misc.Logger.Info("voice txt taskId", zap.Uint64("taskId", taskId))

	request2 := asr.NewDescribeTaskStatusRequest()

	request2.TaskId = common.Uint64Ptr(taskId)

	var response2 *asr.DescribeTaskStatusResponse

	t := time.After(8 * time.Second)

	for {
		select {
		case <-t:
			return "", errors.New("fail to process voice")
		default:
			time.Sleep(1 * time.Second)

			response2, err = misc.CloudClient.DescribeTaskStatus(request2)
			if _, ok := err.(*terrors.TencentCloudSDKError); ok {
				misc.Logger.Error("An API error has returned: %s", zap.Error(err))
				return "", errors.New("api error")
			}
			if err != nil {
				misc.Logger.Error("An API error has returned: %s", zap.Error(err))
				return "", err
			}
			misc.Logger.Info("api response2", zap.Any("rsp", response2))
			status := *response2.Response.Data.Status
			if status == 2 {
				dataList := response2.Response.Data.ResultDetail

				buffer := bytes.Buffer{}

				for _, data := range dataList {
					buffer.WriteString(*data.FinalSentence)
				}

				txt = buffer.String()

				misc.Logger.Info("voice to txt result", zap.String("txt", txt))
				return txt, nil
			} else if status == 3 {
				return "", errors.New("fail to process voice")
			}
		}

	}
}
