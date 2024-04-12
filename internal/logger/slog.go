package logger

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

func NewLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))
}

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
		return AnsiGreen + " " + strconv.Itoa(status) + " " + AnsiNoColor
	case status > 299 && status < 400:
		return AnsiBlue + " " + strconv.Itoa(status) + " " + AnsiNoColor
	case status > 399 && status < 500:
		return AnsiYellow + " " + strconv.Itoa(status) + " " + AnsiNoColor
	default:
		return AnsiRed + " " + strconv.Itoa(status) + " " + AnsiNoColor
	}
}

func HumanReadableBytes(n int) string {
	switch {
	case n < 1024:
		return fmt.Sprintf("%d B", n)
	case n >= 1024 && n < 1048576:
		return fmt.Sprintf("%.2f KB", float64(n)/1024.0)
	default:
		return fmt.Sprintf("%.2f GB", float64(n)/1048576.0)
	}
}
