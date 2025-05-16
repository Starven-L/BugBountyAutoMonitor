package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// 工具函数：安全获取字符串字段
func getStr(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok && val != nil {
		if strVal, ok := val.(string); ok {
			return strVal
		}
	}
	return "nil"
}

// 工具函数：安全获取浮点数字段
func getFloat(m map[string]interface{}, key string) float64 {
	if val, ok := m[key]; ok && val != nil {
		if floatVal, ok := val.(float64); ok {
			return floatVal
		}
	}
	return 0.0
}

func CheckReportUpdataData() error {
	config2, err := ReadConfig()
	if err != nil {
		return errors.New("读取配置文件失败")
	}

	// 声明一个数组存放所有的 engagementCode
	// 对比原始js代码传入的是Engagements数组，此处二开变更为不传参，直接在函数里面读取得到Engagements数组
	for i := range config2.Engagements {
		Engagement := &config2.Engagements[i]
		url := "https://bugcrowd.com/engagements/" + Engagement.EngagementCode + "/crowdstream.json?page=1&filter_by=" + strings.Join(Engagement.CrowdStream.FilterBy, ",")
		//获取response的json数据
		resp, _ := http.Get(url)
		bodyBytes, _ := io.ReadAll(resp.Body)
		defer resp.Body.Close()
		var respjson map[string]interface{}
		_ = json.Unmarshal(bodyBytes, &respjson)

		results, _ := respjson["results"].([]interface{})

		if Engagement.CrowdStream.LastReportId == nil {
			if len(results) == 0 {
				_ = SendTelegramMsg("【CrowdStream监控通知】\n" + Engagement.Name + "监控crowdstream无变化(无报告)")
				continue
			}
			id, _ := results[0].(map[string]interface{})["id"].(string)
			Engagement.CrowdStream.LastReportId = &id
			imagePath := results[0].(map[string]interface{})["logo_url"].(string)

			firstResult, ok := results[0].(map[string]interface{})
			if !ok {
				return fmt.Errorf("invalid result format")
			}
			tgMsg := "【CrowdStream监控通知】" + Engagement.Name + "\n"
			tgMsg += "<b>🚨监控发现新报告 <a href=\"https://bugcrowd.com" + getStr(firstResult, "engagement_path") + "\">" + Engagement.Name + "</a> 🚨 </b>\n"
			tgMsg += "<b>✅漏洞报告名称：" + getStr(firstResult, "engagement_name") + "</b>\n"
			tgMsg += "<b>✅漏洞状态：" + getStr(firstResult, "substate") + "</b>\n"
			tgMsg += "<b>✅漏洞级别：p" + fmt.Sprintf("%d", int(getFloat(firstResult, "priority"))) + "</b>\n"
			tgMsg += "<b>✅漏洞报告时间：" + getStr(firstResult, "created_at") + "</b>\n"
			tgMsg += "<b>✅漏洞公开披露时间：" + getStr(firstResult, "disclosed_at") + "</b>\n"
			tgMsg += "<b>✅漏洞赏金：" + fmt.Sprintf("%.2f", getFloat(firstResult, "amount")) + "</b>\n"
			tgMsg += "<b>✅漏洞积分：" + fmt.Sprintf("%.2f", getFloat(firstResult, "point")) + "</b>\n"
			tgMsg += "<b>✅漏洞目标站点：" + getStr(firstResult, "target") + "</b>\n"
			tgMsg += "<b>✅漏洞报告链接：<a href=\"https://bugcrowd.com/" + getStr(firstResult, "disclosure_report_url") + "\">点这里嗷</a></b>\n"
			tgMsg += "<b>✅漏洞提交研究员：" + getStr(firstResult, "researcher_username") + "</b>\n"
			err = SendTelegramMsgWithImage(tgMsg, imagePath)
		} else {
			id, _ := results[0].(map[string]interface{})["id"].(string)
			// 如果没有新公告，直接返回
			if *Engagement.CrowdStream.LastReportId == id {
				_ = SendTelegramMsg("【CrowdStream监控通知】\n" + Engagement.Name + "监控报告crowdstream无变化")
				continue
			}
			Engagement.CrowdStream.LastReportId = &id
			// 发送通知
			if config2.Notifications.Telegram {
				id, _ := results[0].(map[string]interface{})["id"].(string)
				Engagement.CrowdStream.LastReportId = &id
				imagePath := results[0].(map[string]interface{})["logo_url"].(string)

				firstResult, ok := results[0].(map[string]interface{})
				if !ok {
					return fmt.Errorf("invalid result format")
				}
				tgMsg := "【CrowdStream监控通知】" + Engagement.Name + "\n"
				tgMsg += "<b>🚨监控发现新报告 <a href=\"https://bugcrowd.com" + getStr(firstResult, "engagement_path") + "\">" + Engagement.Name + "</a> 🚨 </b>\n"
				tgMsg += "<b>✅漏洞报告名称：" + getStr(firstResult, "engagement_name") + "</b>\n"
				tgMsg += "<b>✅漏洞状态：" + getStr(firstResult, "substate") + "</b>\n"
				tgMsg += "<b>✅漏洞级别：p" + fmt.Sprintf("%d", int(getFloat(firstResult, "priority"))) + "</b>\n"
				tgMsg += "<b>✅漏洞报告时间：" + getStr(firstResult, "created_at") + "</b>\n"
				tgMsg += "<b>✅漏洞公开披露时间：" + getStr(firstResult, "disclosed_at") + "</b>\n"
				tgMsg += "<b>✅漏洞赏金：" + fmt.Sprintf("%.2f", getFloat(firstResult, "amount")) + "</b>\n"
				tgMsg += "<b>✅漏洞积分：" + fmt.Sprintf("%.2f", getFloat(firstResult, "point")) + "</b>\n"
				tgMsg += "<b>✅漏洞目标站点：" + getStr(firstResult, "target") + "</b>\n"
				tgMsg += "<b>✅漏洞报告链接：<a href=\"https://bugcrowd.com/" + getStr(firstResult, "disclosure_report_url") + "\">点这里嗷</a></b>\n"
				tgMsg += "<b>✅漏洞提交研究员：" + getStr(firstResult, "researcher_username") + "</b>\n"
				err = SendTelegramMsgWithImage(tgMsg, imagePath)
			}
		}
		//将更新的数据写入config.json文件
		if err := writeConfig(config2); err != nil {
			return err
		}
	}
	return nil
}

/*
tgMsg += "<b>🚨监控发现新报告 <a href=\"https://bugcrowd.com" + results[0].(map[string]interface{})["engagement_path"].(string) + "\">" + Engagement.Name + "</a> 🚨 </b>\n"
			tgMsg += "<b>✅漏洞报告名称：" + results[0].(map[string]interface{})["engagement_name"].(string) + "</b>\n"
			tgMsg += "<b>✅漏洞状态：" + results[0].(map[string]interface{})["substate"].(string) + "</b>\n"
			tgMsg += "<b>✅漏洞级别：p" + fmt.Sprintf("%d", int(results[0].(map[string]interface{})["priority"].(float64))) + "</b>\n"
			tgMsg += "<b>✅漏洞报告时间：" + results[0].(map[string]interface{})["created_at"].(string) + "</b>\n"
			tgMsg += "<b>✅漏洞公开披露时间：" + results[0].(map[string]interface{})["disclosed_at"].(string) + "</b>\n"
			tgMsg += "<b>✅漏洞赏金：" + fmt.Sprintf("%.2f", func() float64 {
				if amount, ok := results[0].(map[string]interface{})["amount"]; ok && amount != nil {
					return amount.(float64)
				}
				return 0.0
			}()) + "</b>\n"
			tgMsg += "<b>✅漏洞积分：" + fmt.Sprintf("%.2f", func() float64 {
				if amount, ok := results[0].(map[string]interface{})["point"]; ok && amount != nil {
					return amount.(float64)
				}
				return 0.0
			}()) + "</b>\n"
			tgMsg += "<b>✅漏洞目标站点：" + results[0].(map[string]interface{})["target"].(string) + "</b>\n"
			tgMsg += "<b>✅漏洞报告链接：<a href=\"https://bugcrowd.com/" + results[0].(map[string]interface{})["disclosure_report_url"].(string) + "\" >点这里嗷</a></b>\n"
			tgMsg += "<b>✅漏洞提交研究员：" + results[0].(map[string]interface{})["researcher_username"].(string) + "</b>\n"
			err = SendTelegramMsgWithImage(tgMsg, imagePath)
*/
