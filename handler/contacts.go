package handler

import (
	"github.com/TelePilot/htmx-contacts/model"
	"github.com/TelePilot/htmx-contacts/view/contacts"
	"github.com/labstack/echo/v4"
)

func HandleUserShow(c echo.Context, con []model.Contact) error {
	query := c.QueryParam("q")
	val := ""
	if len(query) >= 0 {
		val = query
	}

	return render(c, contacts.View(val, con, c))
}

func HandleSingleContact(c echo.Context, con model.Contact) error {
	return render(c, contacts.Contact(con))
}
