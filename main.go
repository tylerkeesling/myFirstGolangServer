package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

type Users struct {
	ID        int    `gorm:"AUTO_INCREMENT" form:"id" json:"id"`
	Firstname string `gorm:"not null" form:"firstname" json:"firstname"`
	Lastname  string `gorm:"not null" form:"lastname" json:"lastname"`
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Next()
	}
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

	r.Use(Cors())

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
	// connect to and disconnect from db
	db := InitDb()
	defer db.Close()

	var users []Users
	// this is SELECT * FROM users
	db.Find(&users)

	// display the result
	c.JSON(200, users)
}

func GetUser(c *gin.Context) {
	db := InitDb()
	defer db.Close()

	id := c.Params.ByName("id")
	var user Users
	// SELECT * FROM users WHERE id = 1
	db.First(&user, id)

	if user.ID != 0 {
		c.JSON(200, user)
	} else {
		c.JSON(404, gin.H{"error": "User not found."})
	}
}

func PostUser(c *gin.Context) {
	db := InitDb()
	defer db.Close()

	var user Users
	c.Bind(&user)

	if user.Firstname != "" && user.Lastname != "" {
		// insert user into db
		db.Create(&user)
		c.JSON(201, gin.H{"success": user})
	} else {
		// error
		c.JSON(422, gin.H{"error": "Fields are empty."})
	}
}

func UpdateUser(c *gin.Context) {
	// connect/disconnect from db
	db := InitDb()
	defer db.Close()

	// Get user id
	id := c.Params.ByName("id")
	var user Users
	// SELECT * FROM users WHERE id = 1;
	db.First(&user, id)

	if user.Firstname != "" && user.Lastname != "" {
		if user.ID != 0 {
			var newUser Users
			c.Bind(&newUser)

			result := Users{
				ID:        user.ID,
				Firstname: newUser.Firstname,
				Lastname:  newUser.Lastname,
			}
			// UPDATE users SET ....
			db.Save(&result)
			c.JSON(200, gin.H{"success": result})
		} else {
			c.JSON(404, gin.H{"error": "Fields are empty."})
		}
	}
}

func DeleteUser(c *gin.Context) {
	// Connect and disconnect from db
	db := InitDb()
	defer db.Close()

	// Get user ID
	id := c.Params.ByName("id")
	var user Users
	// SELECT * FROM users WHERE id = 1
	db.First(&user, id)

	if user.ID != 0 {
		// DELETE FROM users WHERE id = user.ID
		db.Delete(&user)
		c.JSON(200, gin.H{"Success": "User #" + id + " successfully deleted."})
	} else {
		c.JSON(404, gin.H{"error": "User not found."})
	}
}
