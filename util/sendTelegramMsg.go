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

type MessageRequest struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

func SendTelegramMsg(message string) error {
	//从.env文件读取
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

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)
	messageRequest := MessageRequest{
		ChatID:    chatID,
		Text:      message,
		ParseMode: "html",
	}
	jsonData, err := json.Marshal(messageRequest)
	if err != nil {
		err = errors.New("json解析失败" + err.Error())
		return err
	}

	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		err = errors.New("telegramBOT消息发送失败:" + err.Error())
		return err
	}

	if resp.StatusCode != http.StatusOK {
		err = errors.New("telegramBOT消息发送失败,响应码：" + resp.Status)
		return err
	}

	defer resp.Body.Close()

	return nil
}
