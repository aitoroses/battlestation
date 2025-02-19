package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type Status struct {
	Generation int  `json:"generation"`
	Available  bool `json:"available"`
}

type FireRequest struct {
	Target struct {
		X int `json:"x"`
		Y int `json:"y"`
	} `json:"target"`
	Enemies int `json:"enemies"`
}

type FireResponse struct {
	Casualties int `json:"casualties"`
	Generation int `json:"generation"`
}

type IonCannon struct {
	generation int
	fireTime   float64
	lastFired  time.Time
	mu         sync.RWMutex
}

func (c *IonCannon) isAvailable() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.lastFired.IsZero() {
		return true
	}

	return time.Since(c.lastFired) >= time.Duration(c.fireTime*float64(time.Second))
}

func (c *IonCannon) handleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	status := Status{
		Generation: c.generation,
		Available:  c.isAvailable(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func (c *IonCannon) handleFire(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.isAvailable() {
		http.Error(w, "Cannon not available", http.StatusServiceUnavailable)
		return
	}

	var req FireRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Record fire time
	c.lastFired = time.Now()

	// Return all enemies as casualties
	resp := FireResponse{
		Casualties: req.Enemies,
		Generation: c.generation,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	port := flag.Int("port", 8080, "HTTP server port")
	flag.Parse()

	// Get configuration from environment
	generation := 1
	if gen := os.Getenv("GENERATION"); gen != "" {
		fmt.Sscanf(gen, "%d", &generation)
	}

	fireTime := 3.5
	if ft := os.Getenv("FIRE_TIME"); ft != "" {
		fmt.Sscanf(ft, "%f", &fireTime)
	}

	cannon := &IonCannon{
		generation: generation,
		fireTime:   fireTime,
	}

	// Register handlers
	http.HandleFunc("/status", cannon.handleStatus)
	http.HandleFunc("/fire", cannon.handleFire)

	// Start server
	addr := fmt.Sprintf(":%d", *port)
	log.Printf("Ion Cannon (Generation %d) listening on %s", generation, addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
