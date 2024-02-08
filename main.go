package main

import (
	"log/slog"
	"os"

	"github.com/erajayatech/hermes"
	"github.com/vanng822/go-premailer/premailer"
)

func main() {
	h := hermes.Hermes{
		Product: hermes.Product{
			Name: "Spobly",
			Link: "https://github.com/spobly/",
			Logo: "",
		},
	}
	email := hermes.Email{
		Body: hermes.Body{
			Name: "John",
			Intros: []string{
				"We are happy to have you on board",
			},

			Actions: []hermes.Action{
				{
					Instructions: "To get started click here",
					Button: hermes.Button{
						Color:     "#0a0a0a",
						Text:      "Heyyyyyyyyy",
						TextColor: "#fafafa",
					},
				},
			},
		},
	}

	body, err := h.GenerateHTML(email)
	if err != nil {
		slog.Error("body not made")
	}

	err = os.WriteFile("preview.html", []byte(body), 0644)
	if err != nil {
		slog.Error("Unable to write")
	}

	prem, err := premailer.NewPremailerFromString(body, &premailer.Options{
		RemoveClasses:     true,
		CssToAttributes:   true,
		KeepBangImportant: true,
	})
	if err != nil {
		slog.Error("Unable to write")
	}

	html, err := prem.Transform()
	if err != nil {
		slog.Error("Unable to write")
	}

	err = os.WriteFile("clean.html", []byte(html), 0644)
	if err != nil {
		slog.Error("Unable to write")
	}
}
