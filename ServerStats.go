package main

import (
	log "github.com/Sirupsen/logrus"
    "encoding/json"
	// "fmt"
	"net/http"
	"os"
	. "sync/atomic"
	"time"
)

type serverStats struct {
    StartTime time.Time `json:"StartTime"`
	Uptime string `json:"Uptime"`  
	Hostname string `json:"Hostname"`
	ImagesServed int64 `json:"ImagesServed"`
	ChannelLength int `json:"ChannelLength"`
}

func initialServerStats() *serverStats {
    stats := new(serverStats)
    stats.StartTime = time.Now()
	stats.Hostname, _ = os.Hostname()
	
	return stats
}

func (s *serverStats) handler(w http.ResponseWriter, r *http.Request) {
    s.Uptime = time.Since(s.StartTime).String()
	s.ImagesServed = LoadInt64(&(s.ImagesServed))
    
    js, err := json.MarshalIndent(s, "", "  ")
    
	if err != nil {
		log.Error("Error marshalling ServerStats %v", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Write(js)
}