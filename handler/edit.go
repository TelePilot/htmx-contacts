package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
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
	fmt.Println("hello")
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
			if fmt.Sprint(v.Id) == c.Param("id") {
				continue
			}
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
	fmt.Println("hello")
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
		if c.Request().Header.Get("HX-targer") == "delete-button" {
			return c.Redirect(303, "/contacts")
		}
		return c.String(200, "")
	}
	return c.String(404, "Not found")

}

func (h EditHandler) BulkDelete(c echo.Context) error {

	val, err := c.MultipartForm()
	if err != nil {
		return c.Redirect(303, "/contacts?page=1")
	}
	if len(val.Value["selected"]) <= 0 {
		return c.Redirect(303, "/contacts?page=1")
	}
	idMap := make(map[int]string)
	for _, v := range val.Value["selected"] {
		if s, err := strconv.Atoi(v); err == nil {
			idMap[s] = ""
		}
	}
	contacts := ReadContacts()

	var modified []model.Contact
	for i := 0; i < len(contacts); i++ {
		_, ok := idMap[contacts[i].Id]
		if !ok {
			modified = append(modified, contacts[i])
		}
	}

	newFile, err := json.MarshalIndent(modified, "", "  ")
	if err != nil {
		log.Fatal("SHIT")
	}
	os.WriteFile("db/contacts.json", newFile, os.ModePerm)

	return c.Redirect(303, "/contacts?page=1")
}

func (h EditHandler) ValidateEmail(c echo.Context) error {
	contacts := ReadContacts()
	for _, v := range contacts {
		if fmt.Sprint(v.Id) == c.Param("id") {
			continue
		}
		if strings.EqualFold(c.FormValue("email"), v.Email) {
			return c.String(http.StatusOK, "email must be unique")
		}
	}
	return c.String(http.StatusOK, "")
}

const letters = "abcdefghijklmnopqrstuvwxyz"

func GenerateContactView(c echo.Context, err map[string]string) error {
	cont := model.Contact{
		First:  c.FormValue("first"),
		Last:   c.FormValue("last"),
		Phone:  c.FormValue("phone"),
		Email:  c.FormValue("email"),
		Errors: err,
	}
	if c.Request().Header.Get("HX-trigger") == "fake" {

		b := make([]byte, 4)
		for i := 0; i < len(b); i++ {
			b[i] = letters[rand.Intn(len(letters))]
		}
		random := string(b)
		cont.First = random
		cont.Last = random
		cont.Phone = random
		cont.Email = random + "@" + random + "." + random
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
	return c.Redirect(303, "/contacts?page=1")
}
func (h EditHandler) HandleNewContact(c echo.Context) error {
	return GenerateContactView(c, map[string]string{})
}
