package config

type Config struct {
	Engagements   []Engagement `json:"engagements"`
	Notifications Notification `json:"notifications"`
	Proxy         string       `json:"proxy"`
}

type Engagement struct {
	Name             string           `json:"name"`
	EngagementCode   string           `json:"engagementCode"`
	Enabled          bool             `json:"enabled"`
	Platform         string           `json:"platform"`
	Announcements    Announcements    `json:"announcements"`
	CrowdStream      CrowdStream      `json:"crowdStream"`
	SubdomainMonitor SubdomainMonitor `json:"subdomainMonitor"`
}

type Announcements struct {
	Enabled            bool    `json:"enabled"`
	LastAnnouncementId *string `json:"lastAnnouncementId"`
}

type CrowdStream struct {
	Enabled               bool     `json:"enabled"`
	MinimumPriorityNumber int      `json:"minimumPriorityNumber"`
	FilterBy              []string `json:"filterBy"`
	LastReportId          *string  `json:"lastReportId"`
}

type SubdomainMonitor struct {
	Enabled             bool     `json:"enabled"`
	StoreMode           bool     `json:"storeMode"`
	SubdomainsDirectory string   `json:"subdomainsDirectory"`
	ScreenshotEnabled   bool     `json:"screenshotEnabled"`
	HideCodes           []int    `json:"hideCodes"`
	Assets              []string `json:"assets"`
}

type Notification struct {
	Telegram bool `json:"telegram"`
	Discord  bool `json:"discord"`
}
