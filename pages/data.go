package pages

import "fyne.io/fyne/v2"

// Tutorial defines the data structure for a tutorial
type Page struct {
	Title, Intro string
	View         func(w fyne.Window) fyne.CanvasObject
}

var (
	// Tutorials defines the metadata for each tutorial
	Pages = map[string]Page{
		"welcome": {"Welcome", "", welcomeScreen},
		"log": {"日志", "", logScreen},
	}

	// TutorialIndex  defines how our tutorials should be laid out in the index tree
	PageIndex = map[string][]string{
		"":            {"welcome","log"},
		"collections": {"list", "table", "tree"},
		"containers":  {"apptabs", "border", "box", "center", "grid", "split", "scroll"},
		"widgets":     {"accordion", "button", "card", "entry", "form", "input", "text", "toolbar", "progress"},
	}
)
