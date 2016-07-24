package main

import (
	"fmt"
	"time"

	"github.com/robbiev/ui"
)

func main() {
	err := ui.Main(func() {
		radio := ui.NewRadioButtons()
		radio.OnSelected(func(*ui.RadioButtons) {
			fmt.Println(radio.Selected())
		})

		button := ui.NewButton("New")
		text := ui.NewMultilineEntry()

		grid := ui.NewGrid()
		grid.SetPadded(true)
		grid.Append(
			radio,
			0,             // left
			0,             // top
			1,             // xspan
			2,             // yspan
			false,         // hexpand
			ui.AlignStart, // uialign
			false,         // vexpand
			ui.AlignStart, // valign
		)
		grid.Append(
			button,
			1,             // left
			0,             // top
			1,             // xspan
			1,             // yspan
			false,         // hexpand
			ui.AlignStart, // uialign
			false,         // vexpand
			ui.AlignStart, // valign
		)
		grid.Append(
			text,
			1,            // left
			1,            // top
			1,            // xspan
			1,            // yspan
			true,         // hexpand
			ui.AlignFill, // uialign
			true,         // vexpand
			ui.AlignFill, // valign
		)

		window := ui.NewWindow("errnote", 1024, 768, true)
		window.SetMargined(true)
		window.SetChild(grid)

		button.OnClicked(func(*ui.Button) {
			radio.Append(time.Now().Format(time.Stamp))
		})

		window.OnClosing(func(window *ui.Window) bool {
			// when using grid this appears to be necessary
			window.SetChild(nil)
			ui.Quit()
			return true
		})

		window.Show()
	})

	if err != nil {
		panic(err)
	}
}
