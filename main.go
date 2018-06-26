package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"
)

type User struct {
	Email     string   `json:"email"`
	Friends   []string `json:"friends"`
	Block     []string `json:"block"`
	Subscribe []string `json:"subscribe"`
}

type Friend struct {
	Friends []string
}

type Subscribe struct {
	Requestor string
	Target    string
}

type Notification struct {
	Sender string
	Text   string
}

type Person struct {
	Name  string
	Phone string
}

func main() {
	r := gin.Default()

	r.GET("/", home)

	r.POST("/friend_request", friendRequest)
	r.POST("/friend_list", friendList)
	r.POST("/friend_common", friendCommon)
	r.POST("/subscribe", subscribe)
	r.POST("/block", block)
	r.POST("/notification", notification)

	// r.Run() //this code is for local testing run

	//below code using for google cloud app engine run
	http.Handle("/", r)
	appengine.Main()
}

func home(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Friend management api",
	})
}

func friendRequest(c *gin.Context) {

}

func friendList(c *gin.Context) {

}

func friendCommon(c *gin.Context) {

}

func subscribe(c *gin.Context) {

}

func block(c *gin.Context) {

}

func notification(c *gin.Context) {

}
