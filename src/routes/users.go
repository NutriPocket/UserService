// Package routes defines the routes for the API endpoints and the handlers for each route.
package routes

import (
	"net/http"

	controller "github.com/NutriPocket/UserService/controller/users"
	"github.com/NutriPocket/UserService/model"
	"github.com/NutriPocket/UserService/service"
	"github.com/gin-gonic/gin"
)

func UsersRoutes(router *gin.Engine) {
	{
		users_routes := router.Group("/users")
		users_routes.GET("/", getUsers)
		users_routes.GET("/:username", getUser)
		users_routes.PATCH("/:username", updateUser)
	}
}

func getUsers(c *gin.Context) {
	var params model.GetUsersParams

	params.SearchUsername = c.Query("searchUsername")

	service, err := service.NewUserService(nil)
	if err != nil {
		c.Error(err)
		return
	}

	users, err := service.GetAllUsers(params)

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, users)
}

func getUser(c *gin.Context) {
	username := c.Param("username")

	controller := controller.UserController{}

	controller.ValidateString(username, "username")

	service, err := service.NewUserService(nil)
	if err != nil {
		c.Error(err)
		return
	}

	user, err := service.GetUser(username)

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func updateUser(c *gin.Context) {
	username := c.Param("username")

	controller := controller.UserController{}

	controller.ValidateString(username, "username")

	var user *model.EditableUser

	if err := c.ShouldBindJSON(&user); err != nil {
		c.Error(err)
		return
	}

	service, err := service.NewUserService(nil)
	if err != nil {
		c.Error(err)
		return
	}

	ret, err := service.UpdateUser(username, user)

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, ret)
}
