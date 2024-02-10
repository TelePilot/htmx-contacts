package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"strconv"
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
	return c.Redirect(303, "/contacts?page=1")
}

func GenerateContacts(context echo.Context) []model.Contact {
	query := context.QueryParam("q")
	p := context.QueryParam("page")
	page := 1
	if pN, err := strconv.Atoi(p); err == nil {
		page = pN
	}
	val := ""
	if len(query) >= 0 {
		val = strings.ToUpper(query)
	}
	c := ReadContacts()

	var s []model.Contact
	for _, v := range c {
		values := reflect.ValueOf(v)
		for i := 0; i < values.NumField(); i++ {
			if values.Field(i).Kind().String() != "string" {
				continue
			}
			fieldVal := fmt.Sprint(values.Field(i))
			if strings.Contains(strings.ToUpper(fieldVal), val) {
				s = append(s, v)
				break
			}
		}
	}
	max := len(s)
	if page*10 < max {
		max = page * 10
	}
	cont := s[10*(page-1) : max]
	return cont
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
