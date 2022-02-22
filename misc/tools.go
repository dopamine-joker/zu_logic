package misc

import (
	"fmt"
	"github.com/go-basic/uuid"
	"net"
	"os"
)

const (
	tokenPrefix = "token_"

	networkSplit = "@"
	Network      = "tcp"
)

//GenerateTokenKey 生成tokenKey
func GenerateTokenKey(userId int32) string {
	return fmt.Sprintf("%s%d_%s", tokenPrefix, userId, uuid.New())
}

//GetTokenKeyPrefix 删除相关用户token
func GetTokenKeyPrefix(userId int32) string {
	return fmt.Sprintf("%s%d_", tokenPrefix, userId)
}

// GetLocalIP 获取本机IP
func GetLocalIP() string {
	var localIP string
	// 获取本机ip,用于在etcd注册服务
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				localIP = ipnet.IP.String()
			}
		}
	}
	return localIP
}

//PathExists 文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
