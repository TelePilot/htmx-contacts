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
	query := c.QueryParam("q")
	val := ""
	if len(query) >= 0 {
		val = query
	}

	return render(c, contacts.View(val, con, c))
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
