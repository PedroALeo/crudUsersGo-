package routes

import (
	"github.com/PedroALeo/CRUDUsers/handlers"
	"github.com/gin-gonic/gin"
)

func RequestHandlers() { //RequestHandlers function declares the routes needed and call the ru function
	r := gin.Default()
	r.GET("/users", handlers.GetAll)            // Declares a route with the GET method to get all users from the database
	r.GET("/users/:id", handlers.GetById)       // Declares a route with the GET method to get a user in the database specified by id
	r.POST("/users", handlers.CreatUser)        // Declares a route with the POST method to Create a user in the database
	r.DELETE("/users/:id", handlers.DeleteUser) // Declares a route with the DELETE method to delete a user from the database specified by id
	r.PUT("/users/:id", handlers.PutUser)       // Declares a route with the PUT method to update all the fields of an user from the database specified by id
	r.PATCH("/users/:id", handlers.PatchUser)   // Declares a route with the PATCH method to uptade specified fields of an user from the database specfied by id
	r.NoRoute(handlers.Route404)                // Declares a route to return a NOT FOUND status to all the unmapped routes
	r.Run()
}
