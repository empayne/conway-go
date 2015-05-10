package main

import (
	"fmt"
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

func (ui *Ui) Init(BoardSize int) {
	termbox.Init()

	x, y := termbox.Size()
	if BoardSize > x || BoardSize > y {
		msg := fmt.Sprintf("Error! Board size is larger than terminal size (%d x %d)", x, y)
		panic(msg)
	}

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
	var cellColour termbox.Attribute

	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	for row := 0; row < BoardSize; row++ {
		for col := 0; col < BoardSize; col++ {
			if currentGrid[row][col] {
				cellColour = termbox.ColorWhite
			} else {
				cellColour = termbox.ColorBlack
			}
			termbox.SetCell(row, col, UICell, termbox.ColorDefault, cellColour)
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
