package util

import (
	"fmt"
	"os"
	"os/exec"
)

func test() error {
	config5, _ := ReadConfig()
	for _, engagement := range config5.Engagements {
		for _, asset := range engagement.SubdomainMonitor.Assets {
			cmd := exec.Command("cmd", "/C", fmt.Sprintf("del /f /q %s\\%s.txt", engagement.SubdomainMonitor.SubdomainsDirectory, asset))
			cmd.Env = os.Environ()
			err := cmd.Run()
			cmd = exec.Command("cmd", "/C", fmt.Sprintf("reconutil\\assetfinder.exe --subs-only %s -silent > %s\\%s.txt 2>&1", asset, engagement.SubdomainMonitor.SubdomainsDirectory, asset))
			cmd.Env = os.Environ()
			err = cmd.Run()
			//cmd /C reconutil\assetfinder.exe --subs-only nasa.gov -silent > 1.txt 2>&1
			if err != nil {
				_ = SendTelegramMsg("【BugBountyMonitor】❌❌❌执行子域名扫描命令失败" + err.Error())
				return err
			}
		}
	}
	return nil
}
