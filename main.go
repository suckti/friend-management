package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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

func getDB() *mgo.Session {
	// session, err := mgo.Dial("mongodb://127.0.0.1:27017/")
	session, err := mgo.Dial("mongodb://35.240.191.219:27017/")
	if err != nil {
		fmt.Println(err.Error())
	}

	return session
}

func home(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Friend management api",
	})
}

func friendRequest(c *gin.Context) {
	var f Friend
	status := true
	message := ""

	//init db
	sess := getDB()
	defer sess.Close()
	db := sess.DB("imd").C("users")

	//get request from post
	c.BindJSON(&f)

	//get value in friend array
	u1 := f.Friends[0]
	u2 := f.Friends[1]
	block := false

	//first user
	user1 := User{}
	err := db.Find(bson.M{"email": u1}).One(&user1)
	if err != nil { //if not found
		var friend []string

		friend = append(friend, u2)
		db.Insert(&User{u1, friend, []string{}, []string{}})

	} else { //update friend if already exist
		exist := checkSliceExist(user1.Friends, u2) //check if already friend
		if exist == false {
			block := checkSliceExist(user1.Block, u2) //check if blocked
			if block == false {
				friend := append(user1.Friends, u2)
				db.Update(&user1, &User{u1, friend, user1.Block, user1.Subscribe})
			} else {
				status = false
				block = true
				message = "This two user can't be friend"
			}
		} else {
			status = false
			message = "This two user already be friends"
		}
	}

	//second user
	user2 := User{}
	err = db.Find(bson.M{"email": u2}).One(&user2)
	if err != nil { //if not found
		var friend []string

		friend = append(friend, u1)
		db.Insert(&User{u2, friend, []string{}, []string{}})

	} else if block == false { //update friend if already exist and not blocked
		exist := checkSliceExist(user2.Friends, u1)
		if exist == false {
			friend := append(user2.Friends, u1)
			db.Update(&user2, &User{u2, friend, user2.Block, user2.Subscribe})
		} else {
			status = false
			message = "This two user already be friends"
		}
	}

	if status == true {
		c.JSON(200, gin.H{
			"success": status,
		})
	} else {
		c.JSON(200, gin.H{
			"status":  status,
			"message": message,
		})
	}
}

func friendList(c *gin.Context) {
	var user User
	//init db
	sess := getDB()
	defer sess.Close()
	db := sess.DB("imd").C("users")

	//get request from post
	c.BindJSON(&user)
	email := user.Email
	err := db.Find(bson.M{"email": email}).One(&user)
	if err != nil {
		c.JSON(200, gin.H{
			"message": err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"success": true,
			"friends": user.Friends,
			"count":   len(user.Friends),
		})
	}
}

func friendCommon(c *gin.Context) {
	//init db
	sess := getDB()
	defer sess.Close()
	db := sess.DB("imd").C("users")

	var f Friend
	user1 := User{}
	user2 := User{}

	c.BindJSON(&f)
	u1 := f.Friends[0]
	u2 := f.Friends[1]

	err_user1 := db.Find(bson.M{"email": u1}).One(&user1)
	err_user2 := db.Find(bson.M{"email": u2}).One(&user2)
	if err_user1 != nil && err_user2 != nil {
		c.JSON(200, gin.H{
			"message": "user not found",
		})
	} else {
		common := intersection(user1.Friends, user2.Friends)
		c.JSON(200, gin.H{
			"success": true,
			"friends": common,
			"count":   len(common),
		})
	}
}

func subscribe(c *gin.Context) {
	//init db
	sess := getDB()
	defer sess.Close()
	db := sess.DB("imd").C("users")

	var s Subscribe
	c.BindJSON(&s)

	target := User{}

	err := db.Find(bson.M{"email": s.Target}).One(&target)
	if err != nil {
		c.JSON(200, gin.H{
			"message": "target not found",
		})
	} else {
		exist := checkSliceExist(target.Subscribe, s.Requestor)
		if exist == false {
			subscribe := append(target.Subscribe, s.Requestor)
			db.Update(&target, &User{target.Email, target.Friends, target.Block, subscribe}) //being subscriber to target
			c.JSON(200, gin.H{
				"success": true,
			})
		} else {
			c.JSON(200, gin.H{
				"success": false,
				"message": "this requestor already subscribe to target",
			})
		}
	}
}

func block(c *gin.Context) {
	//init db
	sess := getDB()
	defer sess.Close()
	db := sess.DB("imd").C("users")

	var s Subscribe
	c.BindJSON(&s)

	r := User{}

	err := db.Find(bson.M{"email": s.Requestor}).One(&r)
	if err != nil {
		c.JSON(200, gin.H{
			"message": "requestor not found",
		})
	} else {
		exist := checkSliceExist(r.Block, s.Target)
		if exist == false {
			block := append(r.Block, s.Target)
			db.Update(&r, &User{r.Email, r.Friends, block, r.Subscribe})
			c.JSON(200, gin.H{
				"success": true,
			})
		} else {
			c.JSON(200, gin.H{
				"success": false,
				"message": "this requestor already block target",
			})
		}
	}
}

func notification(c *gin.Context) {

}

// HELPER
func intersection(a []string, b []string) (inter []string) {
	low, high := a, b
	if len(a) > len(b) {
		low = b
		high = a
	}

	done := false
	for i, l := range low {
		for j, h := range high {
			f1 := i + 1
			f2 := j + 1
			if l == h {
				inter = append(inter, h)
				if f1 < len(low) && f2 < len(high) {
					if low[f1] != high[f2] {
						done = true
					}
				}
				high = high[:j+copy(high[j:], high[j+1:])]
				break
			}
		}
		if done {
			break
		}
	}
	return
}

func checkSliceExist(elements []string, email string) bool {
	// Create a map of all unique elements.
	for v := range elements {
		if elements[v] == email {
			return true
		}
	}

	return false
}
