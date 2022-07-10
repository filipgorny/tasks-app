package organizer

type Task struct {
	Number   string `json:"number"`
	Label    string `json:"label"`
	Uuid     string `json:"uuid"`
	Done     bool   `json:"done"`
	CreateAt string `json:"createdAt"`
}
