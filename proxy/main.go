package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"crypto/tls"
)

func main() {
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
		print("Fetching URL:", url)
		response, err := http.Get(url)
		print("Received response with status code:", response)
		if err != nil {
			print("Error fetching URL:", err)
		}
		if response.StatusCode == 200 {
			c.JSON(http.StatusOK, gin.H{
				"status": "ok",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status": "error",
			})
		}
	})
	r.Run("localhost:8080")
}