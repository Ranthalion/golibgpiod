package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ranthalion/golibgpiod"
)

func main() {

	showHelp := flag.Bool("h", false, "Display help message")
	showVersion := flag.Bool("v", false, "Display the version")
	flag.Parse()

	if *showHelp {
		printHelp()
		return
	} else if *showVersion {
		printVersion()
		return
	} else {
		for _, gpiochip := range golibgpiod.GetGpioChips() {
			fmt.Println(gpiochip.Name, "-", gpiochip.NumLines, "lines:")
			gpiochip.GetLines()
			for _, line := range gpiochip.Lines {
				name := strings.Trim(line.Name, string([]byte{0}))

				if name == "" {
					name = "unnamed"
				} else {
					name = fmt.Sprintf("\"%v\"", name)
				}

				consumer := strings.Trim(line.Consumer, string([]byte{0}))
				if consumer == "" {
					consumer = "unused"
				} else {
					consumer = fmt.Sprintf("\"%v\"", consumer)
				}

				direction := "output"

				if line.Direction == golibgpiod.GPIOD_LINE_DIRECTION_INPUT {
					direction = "input"
				}

				activeState := "active-low"
				if line.ActiveState == golibgpiod.GPIOD_LINE_ACTIVE_STATE_HIGH {
					activeState = "active-high"
				}

				flags := ""
				if line.Used {
					flags = "[used"
				}
				if line.OpenDrain {
					if len(flags) == 0 {
						flags = "["
					} else {
						flags += " "
					}
					flags += "open-drain"
				}
				if line.OpenSource {
					if len(flags) == 0 {
						flags = "["
					} else {
						flags += " "
					}
					flags += "open-source"
				}

				if len(flags) > 0 {
					flags += "]"
				}

				fmt.Println(fmt.Sprintf("\tline %3v: %12v %12v %7v %12v %v", line.Offset, name, consumer, direction, activeState, flags))
			}
			gpiochip.Close()
		}
	}
}

func printHelp() {
	fmt.Println("Usage: ", os.Args[0], "[OPTIONS] <gpiochip1> ...")
	fmt.Println("Print information about all lines of the specified GPIO chip(s) (or all gpiochips if none are specified).")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  -h:\tdisplay this message and exit")
	fmt.Println("  -v:\tdisplay the version and exit")
}

func printVersion() {
	fmt.Println("Not Versioned")
	fmt.Println("Copyright (C) 2017-2018 Bartosz Golaszewski")
	fmt.Println("License: LGPLv2.1")
	fmt.Println("This is free software: you are free to change and redistribute it.")
	fmt.Println("There is NO WARRANTY, to the extent permitted by law.")
}
