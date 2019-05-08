package main

import (
	"flag"
	"fmt"
	"os"

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
			fmt.Println(fmt.Sprintf("%v [%v] (%v lines)", gpiochip.Name, gpiochip.Label, gpiochip.NumLines))
		}
	}
}

func printHelp() {
	fmt.Println("Usage: ", os.Args[0], "[OPTIONS]")
	fmt.Println("List all GPIO chips, print their labels and number of GPIO lines.")
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
