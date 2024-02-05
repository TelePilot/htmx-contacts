package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/TelePilot/htmx-contacts/handler"
	"github.com/TelePilot/htmx-contacts/model"
	"github.com/labstack/echo/v4"
)

func main() {
	app := echo.New()
	app.Static("/static", "static")
	app.Static("/db", "db")
	redirectHandler := handler.RedirectHandler{}
	app.GET("/", redirectHandler.HandleRedirect)
	app.GET("/contacts", func(c echo.Context) error {
		con := handler.GenerateContacts(c)
		return handler.HandleUserShow(c, con)
	})
	app.GET("/contacts/:id", func(c echo.Context) error {
		con := handler.ReadContacts()
		var contactSingle model.Contact
		for i := 0; i < len(con); i++ {
			if fmt.Sprint(con[i].Id) == c.Param("id") {
				contactSingle = con[i]
				break
			}
		}
		return handler.HandleSingleContact(c, contactSingle)
	})
	app.GET("/contacts/:id/edit", func(c echo.Context) error {
		con := handler.ReadContacts()
		var contactSingle model.Contact
		for i := 0; i < len(con); i++ {
			if fmt.Sprint(con[i].Id) == c.Param("id") {
				contactSingle = con[i]
				break
			}
		}
		return handler.EditContact(c, contactSingle)
	})
	app.POST("/contacts/:id/edit", func(c echo.Context) error {
		contacts := handler.ReadContacts()
		index := -1
		for i := 0; i < len(contacts); i++ {
			if fmt.Sprint(contacts[i].Id) == c.Param("id") {
				index = i
				break
			}
		}
		if index != -1 {
			t := contacts[index]

			contacts[index] = model.Contact{
				Id:     t.Id,
				First:  c.FormValue("first"),
				Last:   c.FormValue("last"),
				Phone:  c.FormValue("phone"),
				Email:  c.FormValue("email"),
				Errors: t.Errors,
			}
			newFile, err := json.MarshalIndent(contacts, "", "  ")
			if err != nil {
				log.Fatal("SHIT")
			}
			os.WriteFile("db/contacts.json", newFile, os.ModePerm)
		}
		return c.Redirect(303, "/contacts")
	})
	app.GET("/contacts/new", func(c echo.Context) error {
		return handler.HandleNewContact(c)
	})
	app.POST("/contacts/new", func(c echo.Context) error {
		contacts := handler.ReadContacts()
		newId := contacts[len(contacts)-1].Id + 1
		new := model.Contact{
			Id:     newId,
			First:  c.FormValue("first"),
			Last:   c.FormValue("last"),
			Phone:  c.FormValue("phone"),
			Email:  c.FormValue("email"),
			Errors: map[string]string{},
		}
		contacts = append(contacts, new)

		newFile, err := json.MarshalIndent(contacts, "", "  ")
		if err != nil {

			return handler.HandleNewContact(c)
		}
		er := os.WriteFile("db/contacts.json", newFile, os.ModePerm)
		if er != nil {
			return handler.HandleNewContact(c)
		}
		return c.Redirect(303, "/contacts")
	})
	app.Start(":3000")

}
