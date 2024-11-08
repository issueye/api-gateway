package bootstrap

import (
	"api-gateway/utils"
	"fmt"
	"os"
)

func InitRuntime() {
	// 检查本地是否存在runtime文件夹
	// 获取当前程序的路径
	path := GetWorkDir()
	rtPath := isExistsCreatePath(path, "runtime")
	ROOT_PATH = rtPath
	DB_PATH = isExistsCreatePath(rtPath, "data")
	CONFIG_PATH = isExistsCreatePath(rtPath, "config")
	LOG_PATH = isExistsCreatePath(rtPath, "logs")
}

// GetWorkDir
// 获取程序运行目录
func GetWorkDir() string {
	pwd, _ := os.Getwd()
	return pwd
}

func isExistsCreatePath(path, name string) string {
	p := fmt.Sprintf("%s/%s", path, name)
	exists, err := utils.PathExists(p)
	if err != nil {
		panic(err.Error())
	}

	if !exists {
		panic("创建【config】文件夹失败")
	}

	return p
}
