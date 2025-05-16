package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

func SendTelegramMsgWithImage(message string, imageURL string) error {
	// 从.env文件读取
	envInfo, _ := readEnv()
	botToken := envInfo.TG_TOKEN
	chatID := envInfo.TG_CHATID

	if botToken == "" || chatID == "" {
		err := errors.New("tgBOT相关配置为空请检查配置")
		return err
	}

	// 从 config.proxy 获取代理地址
	config4, _ := ReadConfig()
	proxyAddr := config4.Proxy
	var client *http.Client
	if proxyAddr != "" {
		proxyURL, err := url.Parse(proxyAddr)
		if err != nil {
			return fmt.Errorf("解析代理地址失败: %v", err)
		}
		transport := &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
		client = &http.Client{
			Transport: transport,
			Timeout:   10 * time.Second,
		}
	} else {
		client = &http.Client{
			Timeout: 10 * time.Second,
		}
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendPhoto", botToken)

	// 创建消息请求
	messageRequest := map[string]interface{}{
		"chat_id":    chatID,
		"photo":      imageURL, // 直接使用图片 URL
		"caption":    message,
		"parse_mode": "html",
	}

	// 将请求结构体转换为 JSON
	jsonData, err := json.Marshal(messageRequest)
	if err != nil {
		return fmt.Errorf("json解析失败: %v", err)
	}

	// 发送 POST 请求
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("telegramBOT消息发送失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegramBOT消息发送失败，响应码：%v", resp.Status)
	}

	return nil
}
