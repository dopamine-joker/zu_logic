package misc

import (
	"fmt"
	"github.com/go-basic/uuid"
	"net"
)

const (
	tokenPrefix = "token_"

	networkSplit = "@"
	Network      = "tcp"
)

//GenerateTokenKey 生成tokenKey
func GenerateTokenKey(userId int32) string {
	return fmt.Sprintf("%s_%d_%s", tokenPrefix, userId, uuid.New())
}

//GetTokenKeyPrefix 删除相关用户token
func GetTokenKeyPrefix(userId int32) string {
	return fmt.Sprintf("%s_%d_", tokenPrefix, userId)
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
