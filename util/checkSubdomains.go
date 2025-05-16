package util

import (
	"BugBountyMonitor/config"
	"bufio"
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"
)

type Subdomain struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Subdomain string    `gorm:"unique;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

var wg sync.WaitGroup

func CheckSubdomains() error {
	envInfo, _ := readEnv()
	if envInfo.MYSQL_HOST == "" || envInfo.MYSQL_USER == "" || envInfo.MYSQL_DATABASE == "" || envInfo.MYSQL_PORT == "" {
		err := errors.New("MYSQL相关配置为空请检查配置")
		return err
	}

	configInfo, _ := ReadConfig()

	//连接数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", envInfo.MYSQL_USER, envInfo.MYSQL_PASSWORD, envInfo.MYSQL_HOST, envInfo.MYSQL_PORT, envInfo.MYSQL_DATABASE)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("数据库连接失败: %v", err)
	}

	for _, engagement := range configInfo.Engagements {
		processAsset(engagement, db)
		//获取数据表中总数据条数
		var count int64
		if err := db.Table(engagement.EngagementCode).Count(&count).Error; err != nil {
			fmt.Println("获取数据表总数据数时出错:", err)
		}
		tgMsg := fmt.Sprintf("【BugBountyMonitor数据库检查更新通知】数据表%s目前共有%d条数据", engagement.EngagementCode, count)
		_ = SendTelegramMsg(tgMsg)
	}
	tgMsg := "【子域名检查通知】该轮子域名枚举检查资产监控已完成"
	_ = SendTelegramMsg(tgMsg)
	return nil
}

func processAsset(engagement config.Engagement, db *gorm.DB) {
	// 读取文件路径的子域名并丢到对应数据表
	filePath := engagement.SubdomainMonitor.SubdomainsDirectory
	files, _ := ioutil.ReadDir(filePath)

	for _, file := range files {
		processFile(engagement, db, file.Name())
	}
}

func processFile(engagement config.Engagement, db *gorm.DB, filename string) {
	//defer wg.Done()
	// 读取文件内容
	data, _ := os.Open(engagement.SubdomainMonitor.SubdomainsDirectory + "/" + filename)

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

	if err := scanner.Err(); err != nil {
		fmt.Println("读取文件时出错:", err)
	}

	//协程并发处理数组的数据是否存在于数据库中
	for _, subdomain := range subdomains {
		//wg.Add(1)
		processSubdomain(subdomain, engagement, db, filename)
	}
}

func processSubdomain(subdomain string, engagement config.Engagement, db *gorm.DB, filename string) {
	if err := db.Table(engagement.EngagementCode).Where("subdomain = ?", subdomain).First(&Subdomain{}).Error; err != nil {
		// 如果没有找到，则插入engagement.EngagementCode表中subdomain字段的值为新的子域名
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := db.Table(engagement.EngagementCode).Create(&Subdomain{Subdomain: subdomain}).Error; err != nil {
				fmt.Println("插入数据时出错:", err)
			}
		}
		//NotifyNewSubdomain(subdomain, engagement)
		fmt.Println("新的活跃子域名：" + subdomain)
		tgMsg := "<b>🚨监控发现新的活跃子域名🚨: " + subdomain + " </b>\n"
		_ = SendTelegramMsg(tgMsg)
	}
	return
}

func NotifyNewSubdomain(subdomain string, engagement config.Engagement) {
	//var url string
	//var hideCodes []int
	//for _, hideCode := range engagement.SubdomainMonitor.HideCodes {
	//	hideCodes = append(hideCodes, hideCode)
	//}

	//// 处理url
	//if strings.HasPrefix(subdomain, "https://") || strings.HasPrefix(subdomain, "http://") {
	//	url = subdomain
	//} else {
	//	url = fmt.Sprintf("https://%s", subdomain)
	//	subdomain = fmt.Sprintf("https://%s", subdomain)
	//}
	//
	//// 1. 创建浏览器配置（禁用无头模式）
	//opts := append(chromedp.DefaultExecAllocatorOptions[:],
	//	chromedp.Flag("headless", true),
	//	chromedp.Flag("disable-gpu", true),
	//)
	//
	//// 2. 创建浏览器实例
	//allocCtx, cancelAlloc := chromedp.NewExecAllocator(context.Background(), opts...)
	//defer cancelAlloc()
	//
	//// 3. 创建 ChromeDP 上下文
	//ctx, cancelCtx := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	//defer cancelCtx()
	//
	//// 4. 定义变量存储响应码
	//var responseCode int64
	//var gotResponse bool
	//
	//// 5. 创建单独的监听上下文
	//listenCtx, cancelListen := context.WithCancel(ctx)
	//defer cancelListen()
	//
	//// 6. 监听网络响应，获取第一个状态码后只取消监听
	//chromedp.ListenTarget(listenCtx, func(ev interface{}) {
	//	if ev, ok := ev.(*network.EventResponseReceived); ok {
	//		if !gotResponse {
	//			responseCode = ev.Response.Status
	//			gotResponse = true
	//			cancelListen() // 只取消监听，不影响主上下文
	//		}
	//	}
	//})
	//
	//// 7. 访问目标网页
	//err := chromedp.Run(ctx,
	//	network.Enable(), // 启用网络监听
	//	chromedp.Navigate(url),
	//	chromedp.Sleep(2*time.Second), // 等待页面加载
	//)
	////fmt.Println("跑到这里说明已经完成使命啦1")
	//if err != nil {
	//	//log.Fatal(err)
	//	return
	//}
	////fmt.Println("跑到这里说明已经完成使命啦2")
	//
	////检查responseCode是否在hideCodes中
	//for _, hideCode := range hideCodes {
	//	if int(responseCode) == hideCode {
	//		return
	//	}
	//}
	//fmt.Println("跑到这里说明已经完成使命啦3")

	// 8.发送BOT消息
	//tgMsg := "<b>🚨监控发现新的活跃子域名🚨: " + subdomain + " </b>\n"
	//_ = SendTelegramMsg(tgMsg)
	//fmt.Println("跑到这里说明已经完成使命啦4")

}
