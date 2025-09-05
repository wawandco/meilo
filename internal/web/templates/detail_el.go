package templates

import (
	"fmt"
	"time"

	lucide "github.com/eduardolat/gomponents-lucide"
	"github.com/wawandco/meilo/internal/models"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func DetailEl(e *models.Email) Node {
	return Group{
		Div(
			Class("flex-1 flex flex-col"),
			Header(
				Class("bg-white border-b border-gray-200 p-4 flex justify-between items-center"),
				Div(
					H2(
						Class("text-2xl font-semibold text-gray-900"),
						Text(e.Subject),
					),
					Div(
						Class("mt-1 text-sm text-gray-500"),
						Text("From"),
						Span(
							Class("font-medium text-gray-700"),
							Text(fmt.Sprintf("<%v>", e.Sender)),
						),
						Text("•"),
						Time(
							DateTime(e.ReceivedAt.Format("2006-01-02T15:04")),
							Text(e.ReceivedAt.Format(time.RFC822)),
						),
					),
				),
				Div(
					Class("flex space-x-3"),
					Button(
						Class("p-2 text-gray-500 hover:text-gray-700 hover:bg-gray-100 rounded-md"),
						lucide.Trash2(),
					),
				),
			),
			// Body
			Main(
				Class("flex-1 overflow-y-auto p-6 bg-white"),
				Raw(e.Body),
			),
		),
	}
}
