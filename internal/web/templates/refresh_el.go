package templates

import (
	"github.com/wawandco/meilo/internal/models"
	. "maragu.dev/gomponents"
	//hx "maragu.dev/gomponents-htmx"
	. "maragu.dev/gomponents/html"
)

func RefreshEl(title string, emails []models.Email) Node {
	return Group{
		Div(
			ID("main-content"),
			Div(
				Class("bg-white border-b border-gray-200 p-4"),
				Div(
					Class("flex items-center space-x-4"),
					Text(title),
				),
			),
			Div(
				Class("flex-1 overflow-y-auto bg-white"),
				ListEl(emails),
			),
		),
	}
}
