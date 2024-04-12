package log

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"
)

const (
	AnsiNoColor    = "\033[0m"
	AnsiFaint      = "\033[2m"
	AnsiResetFaint = "\033[22m"
	AnsiRedBG      = "\033[41m"
	AnsiGreenBG    = "\033[42m"
	AnsiYellowBG   = "\033[43m"
	AnsiBlueBG     = "\033[44m"
	AnsiCyanBG     = "\033[46m"
	AnsiRed        = "\033[31m"
	AnsiGreen      = "\033[32m"
	AnsiYellow     = "\033[33m"
	AnsiBlue       = "\033[34m"
	AnsiCyan       = "\033[36m"
)

func ColorizeDuration(d time.Duration) string {
	var color string
	switch {
	case d < 1*time.Millisecond:
		color = AnsiGreen
	case d < 3*time.Millisecond:
		color = AnsiYellow
	default:
		color = AnsiRed
	}

	return fmt.Sprintf("%s %-12v %s", color, d.String(), AnsiNoColor)
}

// status code does not need padding because all of them are 3 digits
func ColorizeStatus(status int) string {
	switch {
	case status > 199 && status < 300:
		return AnsiGreenBG + " " + strconv.Itoa(status) + " " + AnsiNoColor
	case status > 299 && status < 400:
		return AnsiBlueBG + " " + strconv.Itoa(status) + " " + AnsiNoColor
	case status > 399 && status < 500:
		return AnsiYellowBG + " " + strconv.Itoa(status) + " " + AnsiNoColor
	default:
		return AnsiRedBG + " " + strconv.Itoa(status) + " " + AnsiNoColor
	}
}

func NewLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))
}
