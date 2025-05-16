package util

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

/*--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
															获取bugcrowd官网资产更新数据
-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- */

// tip：注意到扩展性-engagement的扩展性
func CheckAssetUpdateData() error {
	config1, err := ReadConfig()
	if err != nil {
		return errors.New("读取配置文件失败")
	}

	// 声明一个数组存放所有的 engagementCode
	// 对比原始js代码传入的是Engagements数组，此处二开变更为不传参，直接在函数里面读取得到Engagements数组
	for i := range config1.Engagements {
		Engagement := &config1.Engagements[i] // 获取指向元素的指针

		url := "https://bugcrowd.com/engagements/" + Engagement.EngagementCode + "/announcements.json"
		// 获取响应的 JSON 数据
		resp, _ := http.Get(url)
		bodyBytes, _ := io.ReadAll(resp.Body)
		defer resp.Body.Close()
		var respjson map[string]interface{}
		_ = json.Unmarshal(bodyBytes, &respjson)

		announcements, _ := respjson["announcements"].([]interface{})

		if Engagement.Announcements.LastAnnouncementId == nil {
			if len(announcements) == 0 {
				_ = SendTelegramMsg("【Announcement监控通知】\n" + Engagement.Name + "监控announcement无变化(无公告)")
				continue
			}
			id, _ := announcements[0].(map[string]interface{})["id"].(string)
			Engagement.Announcements.LastAnnouncementId = &id
			_ = SendTelegramMsg("【Announcement监控通知】\n" + Engagement.Name + "监控发现新公告\n公告url为：https://bugcrowd.com/engagements/" + Engagement.EngagementCode + "/announcements.json")

		} else {
			id, _ := announcements[0].(map[string]interface{})["id"].(string)
			// 如果没有新公告，直接返回
			if *Engagement.Announcements.LastAnnouncementId == id {
				_ = SendTelegramMsg("【Announcement监控通知】\n" + Engagement.Name + "监控announcement无变化")
				continue
			}
			Engagement.Announcements.LastAnnouncementId = &id
			// 发送通知
			_ = SendTelegramMsg("【Announcement监控通知】\n" + Engagement.Name + "监控发现新公告\n公告url为：https://bugcrowd.com/engagements/" + Engagement.EngagementCode + "/announcements.json")
		}
		// 将更新的数据写入config.json文件
		if err := writeConfig(config1); err != nil {
			return err
		}
	}

	return nil
}
