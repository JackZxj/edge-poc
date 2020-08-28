package main

import (
	"fmt"
	"os"
	"time"

	"github.com/stianeikeland/go-rpio"
)

const cycleLength = 32
const pmwClockFrequency = 50 * cycleLength
const timeSpace = 50 // in ms

func main() {
	err := rpio.Open()
	if err != nil {
		os.Exit(1)
	}
	defer rpio.Close()

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

	fmt.Println("red") // but in fact, it lights green
	for i := uint32(0); i <= 32; i++ {
		pinR.DutyCycle(i, cycleLength)
		time.Sleep(time.Millisecond * timeSpace)
	}
	pinR.DutyCycle(0, cycleLength)

	fmt.Println("green") // but in fact, it lights red and blue
	for i := uint32(0); i <= 32; i++ {
		pinG.DutyCycle(i, cycleLength)
		time.Sleep(time.Millisecond * timeSpace)
	}
	pinG.DutyCycle(0, cycleLength)

	fmt.Println("blue") // but in fact, it lights red and blue
	for i := uint32(0); i <= 32; i++ {
		pinB.DutyCycle(i, cycleLength)
		time.Sleep(time.Millisecond * timeSpace)
	}
	pinB.DutyCycle(0, cycleLength)
	fmt.Println("finish")
	rpio.Close()
}
