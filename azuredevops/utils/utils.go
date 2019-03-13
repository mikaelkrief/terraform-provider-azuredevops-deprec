package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

//PrettyPrint json
func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		log.Println(string(b))
	}
	return
}

//PeriodicFunc tick
func PeriodicFunc(tick time.Time) {
	fmt.Println("Tick at: ", tick)
}

func Bool(input bool) *bool {
	return &input
}

func Int32(input int32) *int32 {
	return &input
}

func Int64(input int64) *int64 {
	return &input
}

func Float(input float64) *float64 {
	return &input
}

func String(input string) *string {
	return &input
}