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
