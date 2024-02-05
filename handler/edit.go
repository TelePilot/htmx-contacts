package handler

import (
	"fmt"

	"github.com/TelePilot/htmx-contacts/model"
	"github.com/TelePilot/htmx-contacts/view/edit"
	newContact "github.com/TelePilot/htmx-contacts/view/new"
	"github.com/labstack/echo/v4"
)

func EditContact(c echo.Context, con model.Contact) error {
	return render(c, edit.View(con))
}

type NewHandler struct{}

func generateContact(c echo.Context) model.Contact {
	var errors map[string]string
	fmt.Println(c.FormValue("errors"))
	return model.Contact{
		First:  c.FormValue("first"),
		Last:   c.FormValue("last"),
		Phone:  c.FormValue("phone"),
		Email:  c.FormValue("email"),
		Errors: errors,
	}
}
func HandleNewContact(c echo.Context) error {
	return render(c, newContact.View(generateContact(c)))
}
