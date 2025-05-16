package util

// 实现：检查资产，检查报告，检查子域名
// 注意：检查资产和报告需要将新数据写入到config.json，检查子域名的数据存放于数据库
func execTasks() error {
	//监控资产变化
	err := CheckAssetUpdateData()
	if err != nil {
		return err
	}
	//监控报告变化
	err = CheckReportUpdataData()
	if err != nil {
		return err
	}
	//监控子域名变化
	err = CheckSubdomains()
	if err != nil {
		return err
	}
	return nil
}
