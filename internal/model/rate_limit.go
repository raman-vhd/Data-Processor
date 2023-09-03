package model

type TokenBucketInfo struct {
	Rate       int `redis:"rate"`
	LastRefill int `redis:"lastRefill"`
	Bucket     int `redis:"bucket"`
}

type RateLimit struct {
	UserID    string `json:"userid"`
	ReqPerMin string `json:"ReqPerMin"`
	ReqPerMon string `json:"ReqPerMon"`
}

type ReqSizePerMonInfo struct {
	Max  int `redis:"max"`
	Size int `redis:"size"`
	Date int `redis:"date"`
}
