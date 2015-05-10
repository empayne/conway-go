package main

import (
	"github.com/nsf/termbox-go"
	"time"
)

const UIDelay = 100 * time.Millisecond
const UICell = ' '

/*
	Display-related definitions / functions:
*/

type Ui struct {
	eventQueue chan termbox.Event
}

func (ui *Ui) Init() {
	termbox.Init()
	termbox.SetInputMode(termbox.InputEsc)
	ui.eventQueue = make(chan termbox.Event)

	go func() {
		for {
			ui.eventQueue <- termbox.PollEvent()
		}
	}()
}

func (ui *Ui) Destroy() {
	ui.eventQueue = nil
	termbox.Close()
}

func (ui *Ui) PrintGrid(currentGrid [][]bool, BoardSize int) {
	termbox.Clear(termbox.ColorBlack, termbox.ColorBlack)

	for row := 0; row < BoardSize; row++ {
		for col := 0; col < BoardSize; col++ {
			if currentGrid[row][col] {
				termbox.SetCell(row, col, UICell, termbox.ColorDefault, termbox.ColorWhite)
			} else {
				termbox.SetCell(row, col, UICell, termbox.ColorDefault, termbox.ColorBlack)
			}
		}
	}
}

func (ui *Ui) Update() bool {
	select {
	case <-ui.eventQueue:
		return true
	default:
		termbox.Flush()
		time.Sleep(UIDelay)
		return false
	}
}
