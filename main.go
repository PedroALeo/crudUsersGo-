package main

import (
	"github.com/PedroALeo/CRUDUsers/database"
	"github.com/PedroALeo/CRUDUsers/routes"
)

func main() { //main function call the ConectDB function and defers the CloseDb function and also call the RequestHandlers function
	database.ConectDB()
	defer database.CloseDB()
	routes.RequestHandlers()
}
