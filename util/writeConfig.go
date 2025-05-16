package util

import (
	"BugBountyMonitor/config"
	"encoding/json"
	"errors"
	"os"
)

func writeConfig(configData config.Config) error {
	// Marshal the config data to JSON
	data, err := json.MarshalIndent(configData, "", "    ")
	if err != nil {
		return errors.New("将配置数据转换为JSON失败")
	}

	// 打开配置文件，如果没有则创建
	configFile, err := os.OpenFile("config/config.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return errors.New("打开配置文件失败")
	}
	defer configFile.Close()

	// 将数据写入文件
	_, err = configFile.Write(data)
	if err != nil {
		return errors.New("写入配置文件失败")
	}

	return nil
}
