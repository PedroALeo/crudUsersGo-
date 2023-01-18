package handlers

import (
	"net/http"
	"strconv"

	"github.com/PedroALeo/CRUDUsers/models"
	"github.com/gin-gonic/gin"
)

func GetAll(c *gin.Context) { //GetAll recive a pointer to gin.Context and call the function GetAllUsers from models and shows a http status and the json of users in the database
	users := models.GetAllUsers()
	c.JSON(http.StatusOK, users)
}

func GetById(c *gin.Context) { //GetById recive a pointer to gin.Context, get the id from the URL and call the function GetUserById
	id := c.Params.ByName("id")
	user, err := models.GetUserById(id)
	switch err {
	case models.ErrBadRequestInvalidID:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
	case models.ErrUserNotFound:
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error()})
	default:
		c.JSON(http.StatusOK, user)
	}
}

func CreatUser(c *gin.Context) { //CreateUser recive a pointer to gin.Context, gets the JSON body in the request and call the NewUser function and returns a http status and the created user JSON if there is no error.
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	if err := models.UserValidation(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}

	user, err := models.NewUser(&user)

	if err == models.ErrBadRequest {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)

}

func DeleteUser(c *gin.Context) { //DeleteUser recive a pointer to gin.Context, gets the ID from the URL and call the function DeleteUserById and returns a http status and the deleted user JSON if there is no error.
	id := c.Params.ByName("id")
	user, err := models.DeleteUserById(id)
	switch err {
	case models.ErrBadRequestInvalidID:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
	case models.ErrUserNotFound:
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error()})
	default:
		c.JSON(http.StatusOK, user)
	}
}

func PutUser(c *gin.Context) { //PutUser recive a pointer to gin.Context, gets the ID from the URL and the JSON sended in the request and call the function PutUserById and returns a http status and the updated user JSON if there is no error.
	var user models.User
	id := c.Params.ByName("id")
	user.ID, _ = strconv.Atoi(id)
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	if err := models.UserValidation(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}

	user, err := models.PutUserById(&user)
	switch err {
	case models.ErrBadRequestInvalidID:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
	case models.ErrUserNotFound:
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error()})
	default:
		c.JSON(http.StatusOK, user)
	}

}

func PatchUser(c *gin.Context) { //PatchUser recive a pointer to gin.Context, gets the ID from the URL and the JSON sended in the request and call the function PutUserById and returns a http status and the updated user JSON if there is no error.
	var user models.User
	id := c.Params.ByName("id")
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}

	user.ID, _ = strconv.Atoi(id)

	user, err := models.PatchUserById(&user)
	switch err {
	case models.ErrBadRequestInvalidID:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
	case models.ErrUserNotFound:
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error()})
	default:
		c.JSON(http.StatusOK, user)
	}
}

func Route404(c *gin.Context) { //Route404 recive a pointer to gin.Context, and returns a http.StatusNotFound its called every time a unmapped route is requeried
	c.JSON(http.StatusNotFound, "404 NOT FOUND")
}
