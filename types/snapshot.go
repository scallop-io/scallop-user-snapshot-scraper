package types

type Snapshot struct {
	Users        []UserSnapshot `json:"users"`
	Eligible     []UserSnapshot `json:"eligible"`
	PeriodNumber uint           `json:"periodNumber"`
	PoolBase     string         `json:"poolBase"`
	Timestamp    int64          `json:"timestamp"`
	BlockTime    int64          `json:"blockTime"`
	Winner       string         `json:"winner"`
	RandomSeed   string         `json:"randomSeed"`
	LuckyNumber  string         `json:"luckyNumber"`
	AwardTx      string         `json:"awardTx"`
	VrfTx        string         `json:"vrfTx"`
}

type UserSnapshot struct {
	User     string `json:"user"`
	UserData string `json:"userData"`
	PreUser  string `json:"preUser"`
	NextUser string `json:"nextUser"`
	Amount   string `json:"amount"`
	EditedAt uint64 `json:"editedAt"`
}
