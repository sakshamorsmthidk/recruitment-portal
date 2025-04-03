package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func init() {
	err := godotenv.Load("cmd/server/.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

var oauthConfig = &oauth2.Config{
	ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
	RedirectURL:  "http://localhost:8080/auth/google/callback",
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
	Endpoint:     google.Endpoint,
}

var oauthState = "randomstate"

func main() {
	r := gin.Default()



	log.Println("Client ID:", os.Getenv("GOOGLE_OAUTH_CLIENT_ID"))
	log.Println("Client Secret:", os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"))

	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Sign in with google",
		})
	})

	url := oauthConfig.AuthCodeURL(oauthState)
	log.Println("Generated Auth URL:", url)

	r.GET("/auth/google", func(c *gin.Context) {
		url := oauthConfig.AuthCodeURL(oauthState)
		c.Redirect(http.StatusTemporaryRedirect, url)
	})

	r.GET("/auth/google/callback", func(c *gin.Context) {
		handleGoogleCallback(c)
	})

	fmt.Println("Server is running on http://localhost:8080")
	r.Run(":8080")

}

func handleGoogleCallback(c *gin.Context) {
	state := c.Query("state")
	if state != oauthState {
		log.Println("Invalid Oauth state")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	code := c.Query("code")
	if code == "" {
		log.Println("Authorization code not found")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Println("Failed to exchange authorization code with token: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	client := oauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		log.Println("Failed to fetch user info: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		log.Println("Failed to decode user info: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": userInfo,
	})
}
