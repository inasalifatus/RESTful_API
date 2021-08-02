package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

func init() {
	InitDB()
	InitialMigration()
}

type Config struct {
	DB_Username string
	DB_Password string
	DB_Port     string
	DB_Host     string
	DB_Name     string
}

func InitDB() {
	config := Config{
		DB_Username: "root",
		DB_Password: "2796Intus@",
		DB_Port:     "3306",
		DB_Host:     "localhost",
		DB_Name:     "crud_go",
	}
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		config.DB_Username,
		config.DB_Password,
		config.DB_Host,
		config.DB_Port,
		config.DB_Name,
	)
	var err error
	DB, err = gorm.Open(mysql.Open,(connectionString), &gorm.Config{})
	if err != nil {
		panic(err)
	}
}

type User struct {
	gorm.Model
	Id int `json: "id" form: "id"`
	Name     string `json: "name" form: "name"`
	Email    string `json: "email" form: "email"`
	Password string `json: "password" form: "password"`
}

func InitialMigration() {
	DB.AutoMigrate(&User{})
}

//get all users
func GetUsersController(c echo.Context) error {
	var users []User

	if err := DB.Find(&users).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success get all users",
		"users":   users,
	})
}

//get user by id
func GetUserController(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "your id is wrong",
	})
	var count int64
	DB.Model(&User{}).Where("id = ?", id).Count(&count)
	if count == 0 {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "id not found",
		})
} else {
	var user User
	if x := DB.Find(&user, "id = ?", id); x.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "cannot find",
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"user": &user,
	})
}

//create new user
func CreateUserController(c echo.Context) error {
	user := Users{}
	c.Bind(&user)

	if err := DB.Save(&user).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success create new user",
		"user":    user,
	})
}

// delete user by id
func DeleteUserController(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "invalid id",
		})
	}
	var user User
	if x := DB.Find(&user, "id=?", id); x.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "cannot find",
		})
	}
	if x := DB.Delete(&user); x.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "cannot post",
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success create new user",
		"user":    user,
	})
}

func UpdateUserController(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "invalid id",
		})
	}
	var user User
	if x := DB.Find(&user, "id=?", id); x.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "cannot find",
		})
	}
	c.Bind(&user)
	if x := DB.Save(&user); x.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "cannot update",
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success create new user",
		"user":    user,
	})

}

func main() {
	//create a new echo instance
	e := echo.New()
	//route/ to handler function
	e.GET("/users", GetUsersController)
	e.GET("/users/:id", GetUserController)
	e.POST("users", CreateUserController)
	e.DELETE("/users/:id", DeleteUserController)
	e.PUT("/users/:id", UpdateUserController)

	//start the server, and log if it fails
	e.Logger.Fatal((e.Start(":8080")))
}
