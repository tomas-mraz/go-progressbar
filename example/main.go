package main

import (
	"fmt"
	"github.com/tomas-mraz/go-ansi"
	"go-progressbar/xbar"
	"log"
	"os"
	"syscall"
	"time"
	"unsafe"
)

/*
func progress(p int) {
	// zapamatovat pozici
	x, y := ansi.GetCursorPosition()
	log.Println(x)
	log.Println(y)

	// přejít na poslední řádek a nakreslit
	ansi.CursorAbsolute(0, 20)
	fmt.Print("tady je progress bar")

	// vrátit se na původní pozici
	ansi.CursorAbsolute(x, y)
}*/

type consoleScreenBufferInfo struct {
	size              coord
	cursorPosition    coord
	attributes        word
	window            smallRect
	maximumWindowSize coord
}
type coord struct {
	x Short
	y Short
}
type Short int16
type word uint16
type smallRect struct {
	left   Short
	top    Short
	right  Short
	bottom Short
}

var (
	kernel32                       = syscall.NewLazyDLL("kernel32.dll")
	procGetConsoleScreenBufferInfo = kernel32.NewProc("GetConsoleScreenBufferInfo")
)

func main() {
	f, err := os.OpenFile("main.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	tx, ty := ansi.GetTerminalSize()
	log.Printf("terminal %d x %d", tx, ty)
	sx, sy := ansi.GetScreenSize()
	log.Printf("screen %d x %d", sx, sy)

	handle := syscall.Handle(os.Stdout.Fd())
	var csbi consoleScreenBufferInfo
	procGetConsoleScreenBufferInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&csbi)))
	log.Printf("bottom: %d", int(csbi.window.bottom))

	bar := xbar.New(50)

	/*
		fmt.Println("A")
		bar.Add(1)
		time.Sleep(10 * time.Second)

		fmt.Println("B")
		bar.Add(1)
		time.Sleep(10 * time.Second)

		fmt.Println("C")
		bar.Add(1)
		time.Sleep(10 * time.Second)
	*/

	for x := 0; x < 50; x++ {
		fmt.Println(x)
		bar.Add(1)
		time.Sleep(3 * time.Second)
	}

	bar.End()

	//fmt.Println("AHOJ")
	//fmt.Println("X")
	/*
		area := cursor.NewArea()
		fmt.Print("AAA")
		time.Sleep(3 * time.Second)

		area.Move(5, 0)
		fmt.Print("B")
		time.Sleep(3 * time.Second)

		area.Up(1)
		fmt.Print("C")
		time.Sleep(3 * time.Second)
	*/
	//area.StartOfLineUp(1)
	//fmt.Print("D")
	//	area.Up(2)
	//	fmt.Print("E")

	/*	time.Sleep(1 * time.Second)
		area.StartOfLine()
		area.Move(8, -1)
		fmt.Print("3. Appended row")
	*/
}
