package util

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// 1.从配置文件获取文件路径
// 2.用多个子域名枚举工具扫描
// 3.将扫描结果覆盖到文件
func ScanSubdomains() error {
	config3, _ := ReadConfig()
	for _, engagement := range config3.Engagements {
		if engagement.SubdomainMonitor.Enabled {
			subdomainPath := engagement.SubdomainMonitor.SubdomainsDirectory
			if _, err := os.Stat(subdomainPath); os.IsNotExist(err) {
				// 如果文件不存在，则创建文件
				err := os.MkdirAll(subdomainPath, os.ModePerm)
				if err != nil {
					return err
				}
			}
			// 执行子域名扫描命令：首先获取资产允许范围，然后检查当前系统是Mac还是Windows
			for _, asset := range engagement.SubdomainMonitor.Assets {
				if runtime.GOOS == "darwin" {
					// Mac系统
					err := exec.Command("bash", "-c", fmt.Sprintf("rm -rf %s/%s.txt", engagement.SubdomainMonitor.SubdomainsDirectory, asset)).Run()
					err = exec.Command("bash", "-c", fmt.Sprintf("reconutil/subfinder -d %s -silent | tee %s/%s.txt", asset, engagement.SubdomainMonitor.SubdomainsDirectory, asset)).Run()
					err = exec.Command("bash", "-c", fmt.Sprintf("reconutil/assetfinder --subs-only %s -silent | tee -a %s/%s.txt", asset, engagement.SubdomainMonitor.SubdomainsDirectory, asset)).Run()
					if err != nil {
						_ = SendTelegramMsg("【BugBountyMonitor】❌❌❌执行子域名扫描命令失败" + err.Error())
						return err
					}
				} else if runtime.GOOS == "windows" {
					// Windows系统
					err := exec.Command("cmd", "/C", fmt.Sprintf("del /f /q %s\\%s.txt", engagement.SubdomainMonitor.SubdomainsDirectory, asset)).Run()
					cmd := exec.Command("cmd", "/C", fmt.Sprintf("reconutil\\assetfinder.exe --subs-only %s -silent > %s\\%s.txt 2>&1", asset, engagement.SubdomainMonitor.SubdomainsDirectory, asset))
					err = exec.Command("cmd", "/C", fmt.Sprintf("reconutil\\subfinder.exe -d %s -silent >> %s/%s.txt", asset, engagement.SubdomainMonitor.SubdomainsDirectory, asset)).Run()
					cmd.Environ()
					err = cmd.Run()
					//cmd /C reconutil\assetfinder.exe --subs-only nasa.gov -silent > 1.txt 2>&1

					if err != nil {
						_ = SendTelegramMsg("【BugBountyMonitor】❌❌❌执行子域名扫描命令失败" + err.Error())
						return err
					}
				}
			}
		}
	}
	return nil
}
