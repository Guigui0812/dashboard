package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-yaml"
)

type CacheEntry struct {
	Status    string
	Timestamp time.Time
}

var (
	cache      = make(map[string]CacheEntry)
	cacheMutex sync.RWMutex
	cacheTTL   = 5 * time.Minute
)

func getFromCache(url string) (string, bool) {
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()

	entry, exists := cache[url]
	if !exists {
		return "", false
	}

	if time.Since(entry.Timestamp) > cacheTTL {
		return "", false
	}

	return entry.Status, true
}

func setCache(url string, status string) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	cache[url] = CacheEntry{
		Status:    status,
		Timestamp: time.Now(),
	}
}

type servicesConfig struct {
	Services []struct {
		URL string `yaml:"url"`
	} `yaml:"services"`
}

func loadAllowedURLs() (map[string]bool, error) {
	path := os.Getenv("SERVICES_CONFIG")
	candidates := []string{path, "config/services.yaml", "../config/services.yaml"}

	var data []byte
	var err error
	for _, p := range candidates {
		if p == "" {
			continue
		}
		data, err = os.ReadFile(p)
		if err == nil {
			break
		}
	}
	if err != nil {
		return nil, fmt.Errorf("could not find services.yaml (set SERVICES_CONFIG): %w", err)
	}

	var cfg servicesConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("could not parse services.yaml: %w", err)
	}

	allowed := make(map[string]bool, len(cfg.Services))
	for _, s := range cfg.Services {
		allowed[s.URL] = true
	}
	return allowed, nil
}

func main() {
	allowedURLs, err := loadAllowedURLs()
	if err != nil {
		panic(err)
	}

	// Skip TLS verification: dashboard targets often use self-signed certs.
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		c.Next()
	})

	r.GET("/proxy", func(c *gin.Context) {
		url := c.Query("url")
		if url == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "URL manquante"})
			return
		}
		if !allowedURLs[url] {
			c.JSON(http.StatusForbidden, gin.H{"error": "URL non autorisée"})
			return
		}

		if cachedStatus, found := getFromCache(url); found {
			c.JSON(http.StatusOK, gin.H{
				"status": cachedStatus,
				"cached": true,
			})
			return
		}

		response, err := http.Get(url)
		if err != nil {
			print("Error fetching URL:", err)
			status := "error"
			setCache(url, status)
			c.JSON(http.StatusOK, gin.H{
				"status": status,
				"cached": false,
			})
			return
		}
		defer response.Body.Close()

		var status string
		if response.StatusCode < 500 {
			status = "ok"
		} else {
			status = "error"
		}

		setCache(url, status)

		c.JSON(http.StatusOK, gin.H{
			"status": status,
			"cached": false,
		})
	})
	// N'écoute que sur localhost : le proxy n'a pas vocation à être joint
	// depuis l'extérieur du conteneur, seul le frontend (même conteneur) l'appelle.
	r.Run("localhost:8080")
}
