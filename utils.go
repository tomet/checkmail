package main

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/go-ini/ini"
)

//--------------------------------------------------------------------------------
// Config-file
//--------------------------------------------------------------------------------

func loadConfigFile(filename string, required bool) *Config {
	c := new(Config)
	
	filename = expandTilde(filename)
	
	debug("Load config-file %q", filename)
	
	cfg, err := ini.InsensitiveLoad(filename)
	if err != nil {
		if os.IsNotExist(err) && required == false {
			debug("Default config-file not found: %q", filename)
			return c
		}
		failure("Failed to load config-file %q: %s", filename, err)
	}
	
	err = cfg.MapTo(c)
	if err != nil {
		failure("Error in config-file %q: %s", filename, err)
	}
	
	return c
}

//--------------------------------------------------------------------------------
// Paths
//--------------------------------------------------------------------------------

func expandTilde(p string) string {
	homeDir, err:= os.UserHomeDir()
	if err != nil {
		homeDir = "/"
	}
	if p == "~" {
		return homeDir
	} else if strings.HasPrefix(p, "~/") {
		return filepath.Join(homeDir, p[2:])
	}
	return p
}

//--------------------------------------------------------------------------------
// Strings
//--------------------------------------------------------------------------------

func utf8Len(s string) int {
	return utf8.RuneCountInString(s)
}

func padLeft(s string, width int) string {
	l := utf8Len(s)

	if l > width {
		return trimToLen(s, width)
	}
	
	runes := []rune(s)

	for l < width {
		runes = append(runes, ' ')
		l++
	}

	return string(runes)
}

func trimToLen(s string, maxLen int) string {
	l := utf8Len(s)

	if l > maxLen && maxLen > 5 {
		runes := []rune(s)
		s = string(runes[:maxLen-3])
		s += "..."
	}

	return s
}

//--------------------------------------------------------------------------------
// Terminal
//--------------------------------------------------------------------------------

func terminalWidth() int {
	env := os.Getenv("COLUMNS")
	w, err := strconv.Atoi(env)
	if err != nil || w < 20 || w > 200 {
		return 80
	}
	return w
}

//--------------------------------------------------------------------------------
// Generic min/max
//--------------------------------------------------------------------------------

type Integer interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64
}

type Float interface {
	float32 | float64
}

type Ordered interface {
	byte | rune | Integer | Float | ~string
}

func maxOf[T Ordered](a, b T) T {
	if b > a {
		return b
	}
	return a
}

func minOf[T Ordered](a, b T) T {
	if b < a {
		return b
	}
	return a
}
