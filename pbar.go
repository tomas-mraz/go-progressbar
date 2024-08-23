package pbar

import (
	"fmt"
	"github.com/tomas-mraz/go-ansi"
	"math"
	"strings"
	"sync"
	"time"
)

type ProgressBar struct {
	startTime      time.Time
	terminalHeight int
	terminalWidth  int
	screenHeight   int
	screenWidth    int
	total          int
	progress       int
	progressBarY   int
	doneStr        string
	emptyStr       string
	lock           sync.Mutex
}

func New(max int) *ProgressBar {
	tx, ty := ansi.GetTerminalSize()
	sx, sy := ansi.GetScreenSize()

	bar := ProgressBar{
		startTime:      time.Now(),
		terminalHeight: tx,
		terminalWidth:  ty,
		screenHeight:   sx,
		screenWidth:    sy,
		progressBarY:   ansi.GetBottomScrollIndex(),
		doneStr:        "#",
		emptyStr:       ".",
		total:          max,
	}

	//TODO show cursor in case app termination
	ansi.CursorHide()
	return &bar
}

func (p *ProgressBar) End() {
	// make enter on the end of app output if necessary
	fmt.Println()

	// zapamatovat pozici nového záznamu
	x, y := ansi.GetCursorPosition()

	// return to previous location of bar and erase it
	ansi.CursorAbsolute(0, p.progressBarY)
	ansi.EraseInLine(3)

	// vrátit se na původní pozici
	ansi.CursorAbsolute(x, y)
	ansi.CursorShow()
}

func (p *ProgressBar) Add(added int) {
	p.progress = p.progress + added

	// make enter on the end of app output if necessary
	fmt.Println()

	// zapamatovat pozici nového záznamu
	x, y := ansi.GetCursorPosition()

	// check last line for future NEW LINE
	lastLine := ansi.IsLastLine()

	// return to previous location of bar and erase it
	ansi.CursorAbsolute(0, p.progressBarY)
	ansi.EraseInLine(3)

	if lastLine {
		// výpis je na posledním řádku > posouvá se okno
		b1 := ansi.GetBottomScrollIndex()
		fmt.Println() // nejde CursorDown když je na konci screen bufferu
		b2 := ansi.GetBottomScrollIndex()
		if b1 == b2 { // detekce konce screen bufferu
			// souřadnice couvají protože se odsouvá obsah screen bufferu
			y--
		}
	} else {
		ansi.LastLine()
	}

	// nová pozice progressBaru
	_, p.progressBarY = ansi.GetCursorPosition() // save Y of progressbar

	// format and print bar line
	consumedTime := time.Since(p.startTime)
	totalTime := time.Duration(consumedTime.Seconds()*float64(p.total)/float64(p.progress)) * time.Second
	remainingTime := totalTime - consumedTime
	prefix := fmt.Sprintf("Progress: [%3d%%]", p.progress*100/p.total)
	suffix := fmt.Sprintf("(%s|%s)", formatDuration(consumedTime), formatDuration(remainingTime))
	counter := fmt.Sprintf("%d/%d", p.progress, p.total)
	barWidth := int(math.Abs(float64(p.terminalHeight - (len(prefix) + len(counter) + len(suffix) + 6))))
	barDone := int(float64(barWidth) * float64(p.progress) / float64(p.total)) // Calculate the bar done length
	done := strings.Repeat(p.doneStr, barDone)                                 // Fill the bar with done string
	empty := strings.Repeat(p.emptyStr, barWidth-barDone)                      // Fill the bar with todo string
	bar := fmt.Sprintf("[%s%s]", done, empty)
	fmt.Printf("%s %s %s %s", prefix, bar, counter, suffix)

	// vrátit se na původní pozici
	ansi.CursorAbsolute(x, y)
}

func formatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	var str string
	if h > 0 {
		str += fmt.Sprintf("%dh", h)
	}
	if m > 0 {
		str += fmt.Sprintf("%dm", m)
	}
	str += fmt.Sprintf("%ds", s)
	return str
}
