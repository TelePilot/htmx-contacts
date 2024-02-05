package handler

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/TelePilot/htmx-contacts/model"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func render(c echo.Context, component templ.Component) error {
	return component.Render(c.Request().Context(), c.Response())
}

type RedirectHandler struct{}

func (h RedirectHandler) HandleRedirect(c echo.Context) error {
	return c.Redirect(303, "/contacts")
}

func GenerateContacts(context echo.Context) []model.Contact {
	query := context.QueryParam("q")
	val := ""
	if len(query) >= 0 {
		val = strings.ToUpper(query)
	}
	c := ReadContacts()
	var s []model.Contact
	for _, v := range c {
		values := reflect.ValueOf(v)
		for i := 0; i < values.NumField(); i++ {
			if strings.Contains(strings.ToUpper(values.Field(i).String()), val) {
				s = append(s, v)
				break
			}
		}
	}

	return s
}

func ReadContacts() []model.Contact {
	input, err := os.Open("db/contacts.json")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = input.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	byteValue, err := io.ReadAll(input)
	if err != nil {
		log.Fatal(err)
	}
	var contacts []model.Contact

	er := json.Unmarshal(byteValue, &contacts)
	if er != nil {
		log.Fatal(er, "marshal")
	}
	return contacts
}
