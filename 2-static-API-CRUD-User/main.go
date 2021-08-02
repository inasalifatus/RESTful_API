package main

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type User struct {
	Id       int    `json:"id" form:"id"`
	Name     string `json:"name" form:"name"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

var users []User

//------------------controller--------------------------

// to get all users
func GetUsersController(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success get all users",
		"users":   users,
	})
}

// to get user by id
func GetUserController(c echo.Context) error {
	userId, _ := strconv.Atoi(c.Param("userId"))
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success get user",
		"users":   users[userId],
	})
}

// to delete user by id
func DeleletUserController(c echo.Context) error {
	userId := c.Param("id")
	j, _ := strconv.Atoi(userId)

	for i := j; i < len(users); i++ {
		users[i].Id -= 1
	}
	newUsers := append(users[:j-1], users[j:]...)
	users = newUsers
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success deleted user",
	})
}

// to update user by id
func UpdateUserController(c echo.Context) error {
	user := new(User)
	if err := c.Bind(user); err != nil {
		return err
	}
	userId, _ := strconv.Atoi(c.Param("userId"))
	if users[userId].Name != "" {
		users[userId].Name = user.Name
	}
	if users[userId].Email != "" {
		users[userId].Email = user.Email
	}
	if users[userId].Password != "" {
		users[userId].Password = user.Password
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success get user",
		"users":   users[userId],
	})
}

// to create new user
func CreateUserController(c echo.Context) error {
	user := User{}
	c.Bind(&user)

	if len(users) == 0 {
		user.Id = 1
	} else {
		newId := users[len(users)-1].Id + 1
		user.Id = newId
	}
	users = append(users, user)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"messages": "success create user",
		"user":     user,
	})
}

//----------------------main------------------------------
func main() {
	e := echo.New()
	//routing with query parameter
	e.GET("/users", GetUsersController)
	e.GET("/users/:userId", GetUserController)
	e.POST("/users", CreateUserController)
	e.PUT("/users/:userId", UpdateUserController)
	e.DELETE("/users/:userId", DeleletUserController)

	// start the server, and log if it fails
	e.Logger.Fatal(e.Start(":8080"))
}
