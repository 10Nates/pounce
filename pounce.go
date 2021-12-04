package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

const (
	fps           = 60
	playbox       = 15
	difficultymod = 40
	starttime     = 60
	buffertime    = 4
	lostmsg       = "YOU LOST!"
)

type snake struct {
	pos1    int
	pos2    int
	size    int
	lastpos [2]int
}

type food struct {
	pos1  int
	pos2  int
	speed int
}

var (
	difficulty  = 1
	debugwindow = false
	screen      string
	width       int
	height      int
	mx          int
	my          int
	movebox     [2]int
	player      snake
	critter     food
	timer       float64
	lvltime     float64
	diffs       [5]string
)

func endScreen() {
	//create playbox
	play := [playbox][playbox]string{}

	//render playbox
	var items string
	for i := range play {
		var row string
		for o := range play[i] {
			if play[i][o] == "" {
				play[i][o] = "  "
			}
			row += play[i][o]
		}
		if i == playbox/2 {
			row = row[:len(row)-(len(row)/2+(len(lostmsg)/2))] + lostmsg
		}
		items += row + "\n"
	}

	t := widgets.NewParagraph()
	t.Text = items
	t.Title = "POUNCE"
	t.BorderStyle = ui.NewStyle(ui.ColorGreen)
	mw := (width / 2) - ((playbox*2 + 2) / 2)
	mh := (height / 2) - ((playbox + 2) / 2)
	mwe := mw + (playbox * 2) + 2
	mhe := mh + playbox + 2
	t.SetRect(mw, mh, mwe, mhe)

	s := widgets.NewParagraph()
	s.Text = "Score: " + fmt.Sprint(player.size) + "\nTimer: 0" + "\nDifficulty: " + fmt.Sprint(diffs[difficulty-1])
	s.BorderStyle = ui.NewStyle(ui.ColorGreen)
	s.SetRect(mw, mhe, mwe, mhe+5)

	ui.Render(t, s)
}

func newCritter() {
	rand.Seed(time.Now().UnixNano())

	x := rand.Intn(playbox)
	y := rand.Intn(playbox)
	critter.pos1 = x
	critter.pos2 = y
	critter.speed = player.size
}

func moveCritter() {
	if rand.Float64() > lvltime/starttime && rand.Intn(fps+1) == fps {

		movetype := rand.Intn(4)
		switch movetype {
		case 0:
			if critter.pos1 != 0 {
				critter.pos1 -= 1
			}
		case 1:
			if critter.pos1 != playbox-1 {
				critter.pos1 += 1
			}
		case 2:
			if critter.pos2 != 0 {
				critter.pos2 -= 1
			}
		case 3:
			if critter.pos2 != playbox-1 {
				critter.pos2 += 1
			}
		}

	}
}

func ingame() {
	p := widgets.NewParagraph()
	if debugwindow {
		p.Title = "Debug"
		p.Text += "Tickrate: " + fmt.Sprint(fps)
		p.Text += "\nCritter Move Chance: " + fmt.Sprint(1-lvltime/starttime)
		p.Text += "\nLevel Time: " + fmt.Sprint(lvltime) + " | Multiple: " + fmt.Sprint(float64(difficulty)/100)
		p.SetRect(0, 0, 50, 5)
	}

	//create playbox
	play := [playbox][playbox]string{}

	//playbox tick processing
	timer = timer - 1.0/float64(fps)
	moveCritter()

	if player.pos1 == critter.pos1 && player.pos2 == critter.pos2 {
		player.size += 1
		lvltime *= (1 - float64(difficulty)/difficultymod)
		timer = lvltime + buffertime
		newCritter()
	}

	//playbox modification
	play[movebox[0]][movebox[1]] = "\\033[34m[]\\033[39m"
	play[player.pos1][player.pos2] = "	\\033[31mOO\\033[39m"
	play[critter.pos1][critter.pos2] = "\\033[33m\"<\\033[39m"

	//render playbox
	var items string
	for i := range play {
		var row string
		for o := range play[i] {
			if play[i][o] == "" {
				play[i][o] = "  "
			}
			row += play[i][o]
		}
		items += row + "\n"
	}
	t := widgets.NewParagraph()
	t.Text = items
	t.Title = "POUNCE"
	t.BorderStyle = ui.NewStyle(ui.ColorGreen)
	mw := (width / 2) - ((playbox*2 + 2) / 2)
	mh := (height / 2) - ((playbox + 2) / 2)
	mwe := mw + (playbox * 2) + 2
	mhe := mh + playbox + 2
	t.SetRect(mw, mh, mwe, mhe)

	s := widgets.NewParagraph()
	s.Text = "Score: " + fmt.Sprint(player.size) + "\nTimer: " + fmt.Sprint(int(timer)) + "\nDifficulty: " + fmt.Sprint(diffs[difficulty-1])
	s.BorderStyle = ui.NewStyle(ui.ColorGreen)
	s.SetRect(mw, mhe, mwe, mhe+5)

	if timer > 0.0 {
		ui.Render(p, t, s)
	} else {
		critter = food{}
		endScreen()
	}
}

func move(key string) {
	if screen == "menu" {
		switch key {
		case "w":
			if difficulty != 5 && movebox[0] == 1 {
				difficulty += 1
			}
		case "s":
			if difficulty != 1 && movebox[0] == 1 {
				difficulty -= 1
			}
		case "a":
			if movebox[0] != 0 {
				movebox[0] -= 1
			}
		case "d":
			if movebox[0] != 1 {
				movebox[0] += 1
			}
		case "<Space>", "g":
			if movebox[0] == 0 {
				wipe()
				screen = "game"
			}
		}
		return
	}
	switch key {
	case "w":
		if movebox[0] != 0 {
			movebox[0] -= 1
		}
	case "s":
		if movebox[0] != playbox-1 {
			movebox[0] += 1
		}
	case "a":
		if movebox[1] != 0 {
			movebox[1] -= 1
		}
	case "d":
		if movebox[1] != playbox-1 {
			movebox[1] += 1
		}
	case "<Space>", "g":
		player.lastpos = [2]int{player.pos1, player.pos2}
		player.pos1 = movebox[0]
		player.pos2 = movebox[1]

		if !(player.pos1 == critter.pos1 && player.pos2 == critter.pos2) {
			if player.size > 0 {
				player.size -= 1
			} else {
				timer = 0.0
			}
		}
	}
}

func wipe() {
	timer = starttime
	lvltime = starttime
	movebox = [2]int{0, 0}
	player = snake{}
	critter = food{}

	ui.Clear()
}

func events() {
	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(time.Second / fps).C
	for {
		select {
		case e := <-uiEvents:
			switch e.ID { // event string/identifier
			case "q", "<C-c>": // press 'q' or 'C-c' to quit
				return
			case "<MouseLeft>", "<MouseRight>":
				payload := e.Payload.(ui.Mouse)
				mx, my = payload.X, payload.Y
				move("MOUSECLICK")
			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				width, height = payload.Width, payload.Height
				ui.Clear()
			case "m", "<Escape>":
				if screen == "menu" {
					return
				}
				wipe()
				screen = "menu"
			}
			switch e.Type {
			case ui.KeyboardEvent: // handle all key presses
				move(e.ID)
			}
		// use Go's built-in tickers for updating and drawing data
		case <-ticker:
			render()
		}
	}
}

func menu() {
	sizew := 60
	sizeh := 28

	t := widgets.NewParagraph()
	t.Text = `
	 _|_|_|      _|_|    _|    _|  _|      _|    _|_|_|  _|_|_|_|
	 _|    _|  _|    _|  _|    _|  _|_|    _|  _|        _|
	 _|_|_|    _|    _|  _|    _|  _|  _|  _|  _|        _|_|_|
	 _|        _|    _|  _|    _|  _|    _|_|  _|        _|
	 _|          _|_|      _|_|    _|      _|    _|_|_|  _|_|_|_|
	
	
	 Snake except it's nothing like snake and it's a completely 
	 different game altogether!

	 Catch your prey as they get progressively more evasive and 
	 your metabolism gets progressively faster. Make sure to
	 stay precise! It takes energy to pounce.

	 Move selection:     WASD
	 Pounce:             Space or G
	 Menu:               ESC
	 Exit:               Ctrl+C or Q
	`
	t.Title = "POUNCE"
	t.BorderStyle = ui.NewStyle(ui.ColorGreen)
	mw := (width / 2) - ((sizew + 2) / 2)
	mh := (height / 2) - ((sizeh + 2) / 2)
	mwe := mw + (sizew + 2) + 2
	mhe := mh + sizeh + 2
	t.SetRect(mw, mh, mwe, mhe)

	s := widgets.NewParagraph()
	s.Text = "\n    Play    \n"
	s.BorderStyle = ui.NewStyle(ui.ColorGreen)
	s.SetRect(mw+5, mhe-7, mw+19, mhe-2)

	d := widgets.NewParagraph()
	const dw = 30
	diff := "Difficulty: " + diffs[difficulty-1]
	var dummy = [dw]string{}
	var row string
	for range dummy {
		row += " "
	}
	d.Text = "\n" + row[:len(row)-((len(row)/2)+(len(diff)/2+1))] + diff
	d.BorderStyle = ui.NewStyle(ui.ColorGreen)
	d.SetRect(mwe-dw-6, mhe-7, mwe-5, mhe-2)

	if movebox[0] == 0 {
		s.BorderStyle = ui.NewStyle(ui.ColorBlue)
	}
	if movebox[0] == 1 {
		d.BorderStyle = ui.NewStyle(ui.ColorBlue)
	}
	ui.Render(t, s, d)
}

func render() {
	switch screen {
	case "menu":
		menu()
	case "game":
		ingame()
	}
}

func main() {
	//setup things that can't be constants but should be
	rand.Seed(time.Now().UnixNano())
	diffs = [5]string{"Easy", "Medium", "Tough", "Hard", "Nightmare"}
	if (len(os.Args) > 1) && (os.Args[1] == "debug") {
		debugwindow = true
	}

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	//base variables
	width, height = ui.TerminalDimensions()
	timer = starttime
	lvltime = starttime
	screen = "menu"

	newCritter()
	events()
}
