package main

import (
	"crypto/tls"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
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

func main() {
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
	r.Run("localhost:8080")
}
