package main

import (
	"github.com/TelePilot/htmx-contacts/handler"
	"github.com/labstack/echo/v4"
)

func main() {
	app := echo.New()
	app.Static("/static", "static")
	app.Static("/db", "db")

	redirectHandler := handler.RedirectHandler{}
	contactHandler := handler.ContactHandler{}
	editHandler := handler.EditHandler{}
	app.GET("/", redirectHandler.HandleRedirect)
	app.GET("/contacts", contactHandler.HandleUserShow)
	app.DELETE("/contacts", editHandler.BulkDelete)
	app.GET("/contacts/count", contactHandler.ContactTotal)

	app.GET("/contacts/:id", contactHandler.HandleSingleContact)
	app.DELETE("/contacts/:id", editHandler.ValidateDelete)

	app.GET("/contacts/:id/edit", editHandler.EditContact)
	app.POST("/contacts/:id/edit", editHandler.ValidateEdit)

	app.GET("/contacts/:id/email", editHandler.ValidateEmail)

	app.GET("/contacts/new", editHandler.HandleNewContact)
	app.POST("/contacts/new", editHandler.ValidateNewContact)
	app.Start(":3000")

}
