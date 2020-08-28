package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/stianeikeland/go-rpio"
)

const cycleLength = 16
const pmwClockFrequency = 1000 * cycleLength
const timeSpace = 3000 // in ms

var onlyOneSignalHandler = make(chan struct{})
var shutdownSignals = []os.Signal{os.Interrupt, syscall.SIGTERM}

func setupSignalHandler() (stopCh <-chan struct{}) {
	close(onlyOneSignalHandler) // panics when called twice
	stop := make(chan struct{})
	c := make(chan os.Signal, 1)
	signal.Notify(c, shutdownSignals...)
	go func() {
		<-c
		rpio.StopPwm()
		rpio.Close()
		fmt.Println("exit")
		os.Exit(1)
	}()
	return stop
}

func stoper(stopCh <-chan struct{}) {
	go func() {
		<-stopCh
	}()
}

func main() {
	stoper(setupSignalHandler())

	err := rpio.Open()
	if err != nil {
		os.Exit(1)
	}
	defer func() {
		fmt.Println("exit")
		rpio.StopPwm()
		rpio.Close()
	}()

	// for red led
	pinR := rpio.Pin(12)
	pinR.Mode(rpio.Pwm)
	pinR.Freq(pmwClockFrequency)
	pinR.DutyCycle(0, cycleLength)

	// for green led
	pinG := rpio.Pin(13)
	pinG.Mode(rpio.Pwm)
	pinG.Freq(pmwClockFrequency)
	pinG.DutyCycle(0, cycleLength)

	// for blue led
	pinB := rpio.Pin(19)
	pinB.Mode(rpio.Pwm)
	pinB.Freq(pmwClockFrequency)
	pinB.DutyCycle(0, cycleLength)

	for {
		fmt.Println("#F00") // in fact, it lights green
		pinR.DutyCycle(16, cycleLength)
		pinG.DutyCycle(0, cycleLength)
		pinB.DutyCycle(0, cycleLength)
		time.Sleep(time.Millisecond * timeSpace)

		fmt.Println("#0F0") // in fact, it lights none
		pinR.DutyCycle(0, cycleLength)
		pinG.DutyCycle(16, cycleLength)
		pinB.DutyCycle(0, cycleLength)
		time.Sleep(time.Millisecond * timeSpace)

		fmt.Println("#00F") // in fact, it lights red and blue
		pinR.DutyCycle(0, cycleLength)
		pinG.DutyCycle(0, cycleLength)
		pinB.DutyCycle(16, cycleLength)
		time.Sleep(time.Millisecond * timeSpace)

		fmt.Println("#FF0") // in fact, it lights green
		pinR.DutyCycle(16, cycleLength)
		pinG.DutyCycle(16, cycleLength)
		pinB.DutyCycle(0, cycleLength)
		time.Sleep(time.Millisecond * timeSpace)

		fmt.Println("#F0F") // in fact, it lights white
		pinR.DutyCycle(16, cycleLength)
		pinG.DutyCycle(0, cycleLength)
		pinB.DutyCycle(16, cycleLength)
		time.Sleep(time.Millisecond * timeSpace)

		fmt.Println("#0FF") // in fact, it lights red and blue
		pinR.DutyCycle(0, cycleLength)
		pinG.DutyCycle(16, cycleLength)
		pinB.DutyCycle(16, cycleLength)
		time.Sleep(time.Millisecond * timeSpace)

		fmt.Println("#FFF") // in fact, it lights white
		pinR.DutyCycle(16, cycleLength)
		pinG.DutyCycle(16, cycleLength)
		pinB.DutyCycle(16, cycleLength)
		time.Sleep(time.Millisecond * timeSpace)
	}
}
