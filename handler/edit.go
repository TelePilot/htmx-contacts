package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/TelePilot/htmx-contacts/model"
	"github.com/TelePilot/htmx-contacts/view/edit"
	newContact "github.com/TelePilot/htmx-contacts/view/new"
	"github.com/labstack/echo/v4"
)

type EditHandler struct{}

func (h EditHandler) EditContact(c echo.Context) error {
	con := ReadContacts()
	var contactSingle model.Contact
	for i := 0; i < len(con); i++ {
		if fmt.Sprint(con[i].Id) == c.Param("id") {
			contactSingle = con[i]
			break
		}
	}
	return render(c, edit.View(contactSingle))
}

func (h EditHandler) ValidateEdit(c echo.Context) error {
	contacts := ReadContacts()
	index := -1
	for i := 0; i < len(contacts); i++ {
		if fmt.Sprint(contacts[i].Id) == c.Param("id") {
			index = i
			break
		}
	}
	if index != -1 {
		errs := make(map[string]string)
		for _, v := range contacts {
			if c.FormValue("email") == v.Email {
				errs["email"] = "email must be unique"
				return h.EditContact(c)
			}
		}
		t := contacts[index]

		contacts[index] = model.Contact{
			Id:     t.Id,
			First:  c.FormValue("first"),
			Last:   c.FormValue("last"),
			Phone:  c.FormValue("phone"),
			Email:  strings.ToLower(c.FormValue("email")),
			Errors: t.Errors,
		}
		newFile, err := json.MarshalIndent(contacts, "", "  ")
		if err != nil {
			log.Fatal("SHIT")
		}
		os.WriteFile("db/contacts.json", newFile, os.ModePerm)
	}
	return c.Redirect(303, "/contacts")
}

func (h EditHandler) ValidateDelete(c echo.Context) error {
	contacts := ReadContacts()
	index := -1
	var modified []model.Contact
	for i := 0; i < len(contacts); i++ {
		if fmt.Sprint(contacts[i].Id) == c.Param("id") {
			index = i
		} else {
			modified = append(modified, contacts[i])
		}
	}
	if index != -1 {
		newFile, err := json.MarshalIndent(modified, "", "  ")
		if err != nil {
			log.Fatal("SHIT")
		}
		os.WriteFile("db/contacts.json", newFile, os.ModePerm)
	}
	return c.Redirect(303, "/contacts")
}

func (h EditHandler) ValidateEmail(c echo.Context) error {
	contacts := ReadContacts()
	for _, v := range contacts {
		if strings.EqualFold(c.FormValue("email"), v.Email) {
			return c.String(http.StatusOK, "email must be unique")
		}
	}
	return c.String(http.StatusOK, "")
}

func GenerateContactView(c echo.Context, err map[string]string) error {
	cont := model.Contact{
		First:  c.FormValue("first"),
		Last:   c.FormValue("last"),
		Phone:  c.FormValue("phone"),
		Email:  c.FormValue("email"),
		Errors: err,
	}
	return render(c, newContact.View(cont))
}

func (h EditHandler) ValidateNewContact(c echo.Context) error {
	contacts := ReadContacts()
	errs := make(map[string]string)
	for _, v := range contacts {
		if c.FormValue("email") == v.Email {
			errs["email"] = "email must be unique"
			return GenerateContactView(c, errs)
		}
	}
	newId := contacts[len(contacts)-1].Id + 1
	new := model.Contact{
		Id:     newId,
		First:  c.FormValue("first"),
		Last:   c.FormValue("last"),
		Phone:  c.FormValue("phone"),
		Email:  strings.ToLower(c.FormValue("email")),
		Errors: map[string]string{},
	}
	contacts = append(contacts, new)

	newFile, err := json.MarshalIndent(contacts, "", "  ")
	if err != nil {
		return GenerateContactView(c, errs)
	}
	er := os.WriteFile("db/contacts.json", newFile, os.ModePerm)
	if er != nil {
		return GenerateContactView(c, errs)
	}
	return c.Redirect(303, "/contacts")
}
func (h EditHandler) HandleNewContact(c echo.Context) error {
	return GenerateContactView(c, map[string]string{})
}
