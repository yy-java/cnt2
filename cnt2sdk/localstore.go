package cnt2sdk

type LocalData struct {
	App     string   `json:"app"`
	Profile string   `json:"profile"`
	Configs []Config `json:"configs"`
	Uptime  int64    `json:"uptime"`
}

type LocalFileStore struct {
	targetFile  string
	fileContent string
	updateChan  chan string
	localStore  LocalData
	client      *XClient
}
