package util

import (
	"BugBountyMonitor/config"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

/*--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
															读取配置文件
-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- */

func ReadConfig() (config.Config, error) {
	configFile, err := os.Open("./config/config.json")
	if err != nil {
		return config.Config{}, errors.New("读取配置文件失败")
	}
	defer configFile.Close()

	byteValue, err := ioutil.ReadAll(configFile)
	if err != nil {
		return config.Config{}, errors.New("读取文件内容失败")
	}

	var configContent config.Config
	err = json.Unmarshal(byteValue, &configContent)
	if err != nil {
		return config.Config{}, errors.New("解析配置文件失败")
	}

	return configContent, nil
}
