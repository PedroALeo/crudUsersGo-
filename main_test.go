package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/PedroALeo/CRUDUsers/database"
	"github.com/PedroALeo/CRUDUsers/handlers"
	"github.com/PedroALeo/CRUDUsers/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var (
	ID int
)

func SetupTestRoutes() *gin.Engine { // SetupTestRoutes create test toutes for the tests
	gin.SetMode(gin.ReleaseMode)
	routes := gin.Default()
	return routes
}

func CreateMockUser() { // CreataMockUser creates a mock user and insert it in the database and saves the user id for later
	user := models.User{Name: "MockUserTest", Email: "UserTest@test.com", Password: "password", Phone: "034phone"}
	row := database.DB.QueryRow("INSERT INTO users(nome, email, password, phone) VALUES ($1, $2, $3, $4) RETURNING id",
		user.Name, user.Email, user.Password, user.Password)
	row.Scan(&ID)
}

func DeleteMockUser() { // DeleteMockUser deletes the user created in the CreatMockUser function
	database.DB.QueryRow("DELETE FROM users WHERE id=$1", ID)
}

func TestGetAll(t *testing.T) { // TestGetAll function test the controller to get all the user from the database
	database.ConectDB()
	defer database.CloseDB()
	r := SetupTestRoutes()
	r.GET("/users", handlers.GetAll)
	req, _ := http.NewRequest("GET", "/users", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestGetById(t *testing.T) { // TestGetById function test the controller to get an user from the database specified by id
	database.ConectDB()
	defer database.CloseDB()

	r := SetupTestRoutes()
	r.GET("/users/:id", handlers.GetById)
	path := "/users/" + strconv.Itoa(1)
	req, _ := http.NewRequest("GET", path, nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	var user models.User
	json.Unmarshal(resp.Body.Bytes(), &user)
	fmt.Println(user.Name)
	assert.Equal(t, "testUser", user.Name)
	assert.Equal(t, "test@user.com", user.Email)
	assert.Equal(t, "123456789", user.Password)
	assert.Equal(t, "03499999999", user.Phone)
}

func TestDelete(t *testing.T) { // TestDelete function test the controller to delete an user from the database
	database.ConectDB()
	defer database.CloseDB()
	CreateMockUser()
	r := SetupTestRoutes()
	r.DELETE("/users/:id", handlers.DeleteUser)
	path := "/users/" + strconv.Itoa(ID)
	req, _ := http.NewRequest("DELETE", path, nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestPutUserById(t *testing.T) { // TestPutUserById function test the controller to uptade all the fields of an user from the database
	database.ConectDB()
	defer database.CloseDB()

	CreateMockUser()
	defer DeleteMockUser()

	r := SetupTestRoutes()
	r.PUT("/users/:id", handlers.PutUser)

	path := "/user/" + strconv.Itoa(ID)

	user := models.User{Name: "user editado", Email: "email@editado", Password: "editado", Phone: "ediphone"}
	jsonUser, _ := json.Marshal(user)

	req, _ := http.NewRequest("PUT", path, bytes.NewBuffer(jsonUser))
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	var MockUserEdt models.User
	json.Unmarshal(resp.Body.Bytes(), &MockUserEdt)

	assert.Equal(t, "user editado", MockUserEdt.Name)
	assert.Equal(t, "email@editado", MockUserEdt.Email)
	assert.Equal(t, "editado", MockUserEdt.Password)
	assert.Equal(t, "ediphone", MockUserEdt.Phone)
}
