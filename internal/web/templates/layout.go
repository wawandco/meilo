package templates

import (
	lucide "github.com/eduardolat/gomponents-lucide"
	. "maragu.dev/gomponents"
	hx "maragu.dev/gomponents-htmx"
	. "maragu.dev/gomponents/html"
)

type LayoutProps struct {
	YieldTitle string
	Yield      Node
}

func Layout(props LayoutProps) Node {
	return HTML(
		Lang("en"),
		Head(
			Meta(
				Charset("UTF-8"),
			),
			Meta(
				Name("viewport"),
				Content("width=device-width, initial-scale=1.0"),
			),
			TitleEl(
				Text("Meilo"),
			),
			Script(Src("https://cdn.tailwindcss.com")),
			Script(Src("https://unpkg.com/htmx.org@2.0.6/dist/htmx.min.js")),
			Link(
				Rel("icon"),
				Type("image/x-icon"),
				Attr("sizes", "16x16 32x32"),
				Href("https://wawand.co/images/favicon/favicon.ico"),
			),
		),
		Body(
			Class("bg-gray-50 font-sans"),
			Div(
				Class("flex h-screen"),
				// Sidebar
				Div(
					Class("w-64 bg-white border-r border-gray-200 flex flex-col"),
					// Logo/Header
					Div(
						Class("p-4 border-b border-gray-200"),
						H1(
							Class("text-lg font-semibold text-gray-800"),
							Text("Meilo Dashboard"),
						),
					),
					// Navigation
					Nav(
						Class("flex-1 p-4"),
						Ul(
							Class("space-y-2"),
							Li(
								A(
									Href("/"),
									Class("flex items-center px-3 py-2 text-gray-700 hover:bg-gray-100 rounded-md transition-colors"),
									lucide.Inbox(),
									Span(
										Text(" Inbox"),
									),
								),
							),
						),
					),
				),
				// Main Content
				Div(
					Class("flex-1 flex flex-col"),
					// Header with Search
					Header(
						Class("bg-white border-b border-gray-200 p-4"),
						Div(
							Class("flex items-center justify-between"),
							H2(
								Class("text-xl font-semibold text-gray-800"),
								Text("Inbox"),
							),
							Div(
								Class("flex items-center justify-between"),
								Button(
									hx.Get("/refresh"),
									hx.Target("#main-content"),
									hx.Trigger("click"),
									hx.Swap("innerHTML"),
									Class("p-2 text-gray-500 hover:text-gray-700 hover:bg-gray-100 rounded-md"),
									lucide.RefreshCw(),
								),
								Button(
									Class("p-2 text-gray-500 hover:text-red-700 hover:bg-red-100 rounded-md"),
									lucide.Trash(),
									hx.Delete("/delete-all"),
									hx.Trigger("click"),
									hx.Swap("none"),
									hx.Confirm("⚠️ This action cannot be undone. Delete all emails?"),
								),
							),
						),
					),
					Div(
						ID("main-content"),
						Div(
							Class("bg-white border-b border-gray-200 p-4"),
							Div(
								Class("flex items-center space-x-4"),
								Text(props.YieldTitle),
							),
						),
						Div(
							Class("flex-1 overflow-y-auto bg-white"),
							props.Yield,
						),
					),

				),
			),
		),
	)
}
