package util

import (
	"errors"
	"github.com/joho/godotenv"
	"os"
)

type Env struct {
	TG_TOKEN       string
	TG_CHATID      string
	MYSQL_HOST     string
	MYSQL_PORT     string
	MYSQL_USER     string
	MYSQL_PASSWORD string
	MYSQL_DATABASE string
}

func readEnv() (Env, error) {
	// 加载 .env 文件
	err := godotenv.Load(".env")
	if err != nil {
		err = errors.New("加载 .env 文件失败")
		return Env{}, err
	}

	return Env{
		TG_TOKEN:       os.Getenv("TELEGRAM_BOT_TOKEN"),
		TG_CHATID:      os.Getenv("TELEGRAM_CHAT_ID"),
		MYSQL_HOST:     os.Getenv("MYSQL_HOST"),
		MYSQL_PORT:     os.Getenv("MYSQL_PORT"),
		MYSQL_USER:     os.Getenv("MYSQL_USER"),
		MYSQL_PASSWORD: os.Getenv("MYSQL_PASSWORD"),
		MYSQL_DATABASE: os.Getenv("MYSQL_DATABASE"),
	}, nil
}
