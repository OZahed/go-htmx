package main

import (
	"fmt"
	"strconv"
	"time"
)

const (
	NoColor  = "\u001b[0m"
	RedBG    = "\u001b[41m"
	GreenBG  = "\u001b[42m"
	YellowBG = "\u001b[43m"
	BlueBG   = "\u001b[44m"
	CyanBG   = "\u001b[46m"
	Red      = "\u001b[31m"
	Green    = "\u001b[32m"
	Yellow   = "\u001b[33m"
	Blue     = "\u001b[34m"
	Cyan     = "\u001b[36m"
)

func colorizeDuration(d time.Duration) string {
	var color string
	switch {
	case d < 1*time.Millisecond:
		color = Green
	case d < 3*time.Millisecond:
		color = Yellow
	default:
		color = Red
	}

	return fmt.Sprintf("%s %-12v %s", color, d.String(), NoColor)
}

// status code does not need padding because all of them are 3 digits
func colorizeStatus(status int) string {
	switch {
	case status > 199 && status < 300:
		return GreenBG + " " + strconv.Itoa(status) + " " + NoColor
	case status > 299 && status < 400:
		return CyanBG + " " + strconv.Itoa(status) + " " + NoColor
	case status > 399 && status < 500:
		return YellowBG + " " + strconv.Itoa(status) + " " + NoColor
	default:
		return RedBG + " " + strconv.Itoa(status) + " " + NoColor
	}
}
