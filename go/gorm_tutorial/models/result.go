package models

import "time"

type Result struct {
	Time   time.Time     `json:"lastest"`
	Error  string        `json:"error"`
	Result []interface{} `json:"result"` // This is a dynamic array of any type
}
