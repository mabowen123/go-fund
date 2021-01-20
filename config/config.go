package config

import (
	"fmt"
	"github.com/Unknwon/goconfig"
)

func Load(key string) (config map[string]string) {
	cfg, err := goconfig.LoadConfigFile("./env.ini")
	if err != nil {
		fmt.Println("读取文件错误", err)
	}

	config, _ = cfg.GetSection(key)
	return
}
