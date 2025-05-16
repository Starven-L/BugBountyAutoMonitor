package main

import (
	"BugBountyMonitor/util"
	"fmt"
	"time"
)

func main() {
	//Logo显示与启动通知
	util.ShowLogo()
	tgMsg := "【BugBountyMonitor】漏洞赏金资产监控已启动\n"
	tgMsg += "【+】Author：Starven\n"
	tgMsg += "【+】Email：starvenl@qq.com\n"
	tgMsg += "【+】Team：Syclover三叶草安全技术小组\n"
	err := util.SendTelegramMsg(tgMsg)
	if err != nil {
		fmt.Println("Telegram消息发送失败:", err)
		return
	}

	// 定时检查
	lastStatusTime := time.Now()

	//初次检查
	err = util.CheckAssetUpdateData()
	if err != nil {
		_ = util.SendTelegramMsg("【报错通知❌】\n" + err.Error())
		return
	}
	err = util.CheckReportUpdataData()
	if err != nil {
		_ = util.SendTelegramMsg("【报错通知❌】\n" + err.Error())
		return
	}
	err = util.ScanSubdomains()
	if err != nil {
		_ = util.SendTelegramMsg("【报错通知❌】\n" + err.Error())
		return
	}
	err = util.FirstWriteSubdomainToDb()
	if err != nil {
		_ = util.SendTelegramMsg("【报错通知❌】\n" + err.Error())
		return
	}

	for {
		_ = util.SendTelegramMsg("---------分割线-----------")
		currentTime := time.Now()
		if currentTime.Sub(lastStatusTime).Minutes() >= 60 {
			fmt.Sprintf("检查检查")
			tgMsg := fmt.Sprintf("【BugBountyMonitor】正常运行中\n最后检查时间: %s", currentTime.Format("2006-01-02 15:04:05"))
			tgMsg += "【+】Author：Starven\n"
			tgMsg += "【+】Email：starvenl@qq.com\n"
			tgMsg += "【+】Team：Syclover三叶草安全技术小组\n"
			_ = util.SendTelegramMsg(tgMsg)
			lastStatusTime = currentTime // 重置计时器
			err = util.CheckAssetUpdateData()
			err = util.CheckReportUpdataData()
			err = util.ScanSubdomains()
			err = util.CheckSubdomains()
			if err != nil {
				_ = util.SendTelegramMsg("【报错通知❌】\n" + err.Error())
			}
		}
		time.Sleep(3610 * time.Second) // 每小时检查一次
	}
}
