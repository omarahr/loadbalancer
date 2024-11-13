package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var (
	Address        = ":8080"
	ReqCount int64 = 0
)

var rootCmd = &cobra.Command{
	Use:   "bc",
	Short: "simple backend server",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatal("invalid arguments")
		}

		Address = args[0]
	},
}

func PrintStatus() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop() // Ensure the ticker is stopped properly

	// Infinite loop to keep the program running
	for {
		select {
		case <-ticker.C:
			// Print a message every time the ticker fires
			fmt.Printf("ReqCount: %d\n", ReqCount)
		}
	}
}

func main() {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		return
	}

	go PrintStatus()

	r.GET("/ping", func(c *gin.Context) {
		ReqCount++
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/health", func(c *gin.Context) {
		fmt.Printf("Request from %s\n", c.RemoteIP())
		c.JSON(200, gin.H{})
	})

	if err := r.Run(fmt.Sprintf(":%s", Address)); err != nil {
		log.Fatal(err)
	}
}
