package templates

import (
	"fmt"
	"time"

	"github.com/wawandco/meilo/internal/models"
	. "maragu.dev/gomponents"
	hx "maragu.dev/gomponents-htmx"
	. "maragu.dev/gomponents/html"
)

func ListEl(emails []models.Email) Node {
	return Group{
		If(len(emails) > 0,
			Map(emails, func(e models.Email) Node {
				return Div(
					Class("border-b border-gray-100 hover:bg-gray-50 cursor-pointer"),
					hx.Get(fmt.Sprintf("/details/%v", e.ID)),
					hx.Swap("innerHTML"),
					hx.Trigger("click"),
					hx.Target("#main-content"),
					Div(
						Class("flex items-center p-4"),
						Div(
							Class("flex-1 min-w-0"),
							Div(
								Class("flex items-center justify-between mb-1"),
								Div(
									Class("flex items-center space-x-3"),
									Span(
										Class("font-semibold text-gray-900 text-sm"),
										Text(fmt.Sprintf("From: %s", e.Sender)),
									),
								),
								Span(
									Class("text-xs text-gray-500 whitespace-nowrap"),
									Text(e.ReceivedAt.Format(time.RFC822)),
								),
							),
							Div(
								Class("text-sm text-gray-600 truncate mb-1"),
								Raw(e.Body),
							),
							Div(
								Class("text-sm text-gray-600 truncate mb-1"),
								Text("To: "+e.RawRecipients()),
							),
						),
					),
				)
			}),
		),
		If(len(emails) == 0,
			Div(
				Class("flex items-center p-4"),
				Text("No emails found"),
			),
		),
	}
}
