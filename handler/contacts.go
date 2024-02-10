package handler

import (
	"fmt"

	"github.com/TelePilot/htmx-contacts/model"
	"github.com/TelePilot/htmx-contacts/view/contacts"
	"github.com/labstack/echo/v4"
)

type ContactHandler struct{}

func (h ContactHandler) HandleUserShow(c echo.Context) error {
	con := GenerateContacts(c)

	if c.Request().Header.Get("HX-trigger") == "search" {
		return render(c, contacts.Rows(c, con))
	}
	return render(c, contacts.View(c, con))
}
func (h ContactHandler) ContactTotal(c echo.Context) error {
	contacts := ReadContacts()
	return c.String(200, "("+fmt.Sprint(len(contacts))+" total contacts)")
}

func (h ContactHandler) HandleSingleContact(c echo.Context) error {
	con := ReadContacts()
	var contactSingle model.Contact
	for i := 0; i < len(con); i++ {
		if fmt.Sprint(con[i].Id) == c.Param("id") {
			contactSingle = con[i]
			break
		}
	}
	return render(c, contacts.Contact(contactSingle))
}
