package main

import (
	"fmt"
	"strconv"
)

type QueryResult struct {
	T         Ticker `json:"ticker"`
	Timestamp int    `json:"timestamp"` //Unix timestamp
	Success   bool   `json:"success"`
	Error     string `json:"error"`
}

type Ticker struct {
	Price  string `json:"price"`
	Volume string `json:"volume"`
	Change string `json:"change"`
}

func (q QueryResult) String() string {
	return fmt.Sprintf("%v, %v", q.T, q.Timestamp)
}

func (q QueryResult) ToArray() []string {
	return []string{strconv.Itoa(q.Timestamp), q.T.Price, q.T.Volume, q.T.Change}
}

func (t Ticker) String() string {
	return fmt.Sprintf("%v, %v, %v", t.Price, t.Volume, t.Change)
}
