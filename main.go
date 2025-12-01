package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/robbiev/ui"
)

const (
	defaultDir = "/home/robbie/notes"
)

type note struct {
	title      string
	file       *os.File
	cachedBody *string
	modTime    time.Time
	nameOnDisk string
}

type byModTimeDesc []*note

func (f byModTimeDesc) Len() int           { return len(f) }
func (f byModTimeDesc) Less(i, j int) bool { return f[i].modTime.Before(f[j].modTime) }
func (f byModTimeDesc) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }

func (n *note) Body() (string, error) {
	if n.cachedBody != nil {
		return *n.cachedBody, nil
	}

	b, err := ioutil.ReadAll(n.file)
	if err != nil {
		return "", err
	}

	// +1 for the newline byte
	body := string(b[len(n.title)+1 : len(b)])
	n.cachedBody = &body
	return *n.cachedBody, nil
}

func readNotes(dir string) []*note {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	var notes []*note
	for _, f := range files {
		osFile, err := os.OpenFile(filepath.Join(dir, f.Name()), os.O_RDWR, 0600)
		if err != nil {
			log.Println(err)
			continue
		}

		var title string
		scanner := bufio.NewScanner(osFile)
		if scanner.Scan() {
			title = scanner.Text()
		}
		_, err = osFile.Seek(0, os.SEEK_SET)
		if err != nil {
			log.Println(err)
			continue
		}
		notes = append(notes, &note{
			title:      title,
			file:       osFile,
			modTime:    f.ModTime(),
			nameOnDisk: f.Name(),
		})
	}
	sort.Sort(byModTimeDesc(notes))
	return notes
}

func main() {
	dir := defaultDir
	if len(os.Args) > 1 {
		dir = os.Args[1]
		fmt.Println(dir)
	}

	notes := readNotes(dir)

	err := ui.Main(func() {

		grid, button, destroy := newUI(notes)

		window := ui.NewWindow("errnote", 1024, 768, true)
		window.SetMargined(true)
		window.SetChild(grid)

		var buttonClick func(*ui.Button)
		buttonClick = func(*ui.Button) {
			//radio.Append(time.Now().Format(time.Stamp))
			window.SetChild(nil)
			destroy()

			fileName := fmt.Sprintf("%d", time.Now().UnixNano()/int64(time.Millisecond))
			ioutil.WriteFile(filepath.Join(dir, fileName), []byte(fileName+"\n"), 0600)

			notes = readNotes(dir)
			grid, button, destroy = newUI(notes)
			button.OnClicked(buttonClick)
			window.SetChild(grid)
		}
		button.OnClicked(buttonClick)

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

func newUI(notes []*note) (*ui.Grid, *ui.Button, func()) {
	radio := ui.NewRadioButtons()

	for i := len(notes) - 1; i >= 0; i-- {
		radio.Append(notes[i].title)
	}

	button := ui.NewButton("New")
	text := ui.NewMultilineEntry()
	title := ui.NewEntry()

	title.OnChanged(func(*ui.Entry) {
		note := notes[len(notes)-radio.Selected()-1]
		note.title = title.Text()
		_, _ = note.file.Seek(0, os.SEEK_SET)
		_ = note.file.Truncate(0)
		_, err := note.file.WriteString(note.title + "\n" + *note.cachedBody)
		if err != nil {
			log.Println(err)
		}
	})
	text.OnChanged(func(*ui.MultilineEntry) {
		note := notes[len(notes)-radio.Selected()-1]
		*note.cachedBody = text.Text()
		_, _ = note.file.Seek(0, os.SEEK_SET)
		_ = note.file.Truncate(0)
		_, err := note.file.WriteString(note.title + "\n" + *note.cachedBody)
		if err != nil {
			log.Println(err)
		}
	})

	selector := func(note *note) {
		title.SetText(note.title)
		body, err := note.Body()
		if err != nil {
			log.Println(err)
			return
		}
		text.SetText(body)
	}

	// select the first one
	if len(notes) > 0 {
		selector(notes[len(notes)-1])
	}

	radio.OnSelected(func(*ui.RadioButtons) {
		note := notes[len(notes)-radio.Selected()-1]
		selector(note)
	})

	grid := ui.NewGrid()
	grid.SetPadded(true)
	grid.Append(
		radio,
		0,             // left
		0,             // top
		1,             // xspan
		3,             // yspan
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
		title,
		1,            // left
		1,            // top
		1,            // xspan
		1,            // yspan
		true,         // hexpand
		ui.AlignFill, // uialign
		false,        // vexpand
		ui.AlignFill, // valign
	)
	grid.Append(
		text,
		1,            // left
		2,            // top
		1,            // xspan
		1,            // yspan
		true,         // hexpand
		ui.AlignFill, // uialign
		true,         // vexpand
		ui.AlignFill, // valign
	)
	return grid, button, func() {
		for i := 0; i < 4; i++ {
			grid.Delete(0)
		}
		grid.Destroy()
	}
}
