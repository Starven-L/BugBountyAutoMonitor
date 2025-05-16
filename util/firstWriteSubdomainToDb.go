package util

import (
	"bufio"
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

func FirstWriteSubdomainToDb() error {
	//读取配置文件路径的txt文件数据，逐行写入数据库
	config, err := ReadConfig()
	filename := ""
	envInfo, err := readEnv()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", envInfo.MYSQL_USER, envInfo.MYSQL_PASSWORD, envInfo.MYSQL_HOST, envInfo.MYSQL_PORT, envInfo.MYSQL_DATABASE)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	for _, engagement := range config.Engagements {
		if !db.Migrator().HasTable(engagement.EngagementCode) {
			if err := db.Table(engagement.EngagementCode).AutoMigrate(&Subdomain{}); err != nil {
				err = errors.New("创建" + engagement.EngagementCode + "资产数据表失败")
			}
		}

		filePath := engagement.SubdomainMonitor.SubdomainsDirectory
		files, _ := ioutil.ReadDir(filePath)
		for _, file := range files {
			filename = filePath + "/" + file.Name()
			data, _ := os.Open(filename)
			defer data.Close()
			var subdomains []string
			scanner := bufio.NewScanner(data)

			// 按行读取文件内容
			for scanner.Scan() {
				// 获取每一行并处理
				subdomain := scanner.Text()

				// 去除首尾空白字符
				subdomain = strings.TrimSpace(subdomain)

				// 去除换行符为空格（如果存在）
				subdomain = strings.ReplaceAll(subdomain, "\r", " ")
				subdomain = strings.ReplaceAll(subdomain, "\n", " ")

				// 移除http://或https://协议部分
				subdomain = strings.Replace(subdomain, "http://", "", 1)
				subdomain = strings.Replace(subdomain, "https://", "", 1)

				// 将处理后的子域名添加到数组中
				subdomains = append(subdomains, subdomain)
			}

			//fmt.Println("subdomains共有:" + fmt.Sprintf("%d", len(subdomains)))

			if err := scanner.Err(); err != nil {
				fmt.Println("读取文件时出错:", err)
			}

			//subdomains共有6582个，但是写进数据库的只有3667行:问题排查到其实是去了重的
			for _, subdomain := range subdomains {
				if err := db.Table(engagement.EngagementCode).Where("subdomain = ?", subdomain).First(&Subdomain{}).Error; err != nil {
					// 如果没有找到，则插入engagement.EngagementCode表中subdomain字段的值为新的子域名
					if errors.Is(err, gorm.ErrRecordNotFound) {
						if err := db.Table(engagement.EngagementCode).Create(&Subdomain{Subdomain: subdomain}).Error; err != nil {
							fmt.Println("插入数据时出错:", err)
						}
					}
				}
			}
			//计算数据表总数据数
			var count int64
			if err := db.Table(engagement.EngagementCode).Count(&count).Error; err != nil {
				fmt.Println("获取数据表总数据数时出错:", err)
			}
			tgMsg := fmt.Sprintf("【BugBountyMonitor数据库通知】已将%s资产的子域名写入数据表\n", filename)
			tgMsg += fmt.Sprintf("数据表%s共有%d条数据", engagement.EngagementCode, count)
			_ = SendTelegramMsg(tgMsg)

		}
		tgMsg := fmt.Sprintf("【BugBountyMonitor数据库通知】数据库准备就绪：已将%s所有资产的子域名写入数据表\n", engagement.EngagementCode)
		_ = SendTelegramMsg(tgMsg)
		time.Sleep(2 * time.Second) // 避免频繁请求

	}
	if err != nil {
		return err
	}
	return nil
}
