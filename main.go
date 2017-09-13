package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

type Users struct {
	ID        int    `gorm:"AUTO_INCREMENT" form:"id" json:"id"`
	Firstname string `gorm:"not null" form:"firstname" json:"firstname"`
	Lastname  string `gorm:"not null" form:"lastname" json:"lastname"`
}

func InitDb() *gorm.DB {
	db, err := gorm.Open("sqlite3", "./data.db")
	db.LogMode(true)

	if err != nil {
		panic(err)
	}

	if !db.HasTable(&Users{}) {
		db.CreateTable(&Users{})
		db.Set("gorm:table_options", "ENGINE=InnoDB").CreateTable(&Users{})
	}
	return db
}

func main() {
	r := gin.Default() // establish router

	v1 := r.Group("api/v1") // establish router by group
	{
		// create RESTful API
		v1.GET("/users", GetUsers)
		v1.GET("/users/:id", GetUser)
		v1.POST("/users", PostUser)
		v1.PUT("/users/:id", UpdateUser)
		v1.DELETE("/users/:id", DeleteUser)
	}
	// declare a port
	r.Run(":8080")
}

func GetUsers(c *gin.Context) {
	var users = []Users{
		Users{ID: 1, Firstname: "Tyler", Lastname: "Keesling"},
		Users{ID: 2, Firstname: "Berto", Lastname: "Ortega"},
	}
	c.JSON(200, users)
}

func GetUser(c *gin.Context) {
	id := c.Params.ByName("id")
	userID, _ := strconv.ParseInt(id, 0, 64)

	if userID == 1 {
		content := gin.H{"id": userID, "firstname": "Tyler", "lastname": "Keesling"}
		c.JSON(200, content)
	} else if userID == 2 {
		content := gin.H{"id": userID, "firstname": "Roberto", "lastname": "Ortega"}
		c.JSON(200, content)
	} else {
		content := gin.H{"error": "user with id#" + id + " not found"}
		c.JSON(404, content)
	}
}

func PostUser(c *gin.Context) {

}

func UpdateUser(c *gin.Context) {
	// future code
}

func DeleteUser(c *gin.Context) {
	// future code
}
