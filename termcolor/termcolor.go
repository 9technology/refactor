// Copyright (c) 2012, Vincent Rischmann
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package termcolor

import (
	"fmt"
	"os"
)

type Attribute int
type Color int
type Background int

var (
	Bold      Attribute = 1
	Dark      Attribute = 2
	Reverse   Attribute = 7
	Underline Attribute = 4

	Blue    Color = 34
	Cyan    Color = 36
	Green   Color = 32
	Grey    Color = 30
	Magenta Color = 35
	Red     Color = 31
	White   Color = 37
	Yellow  Color = 33

	BgBlue    Background = 44
	BgCyan    Background = 46
	BgGreen   Background = 42
	BgGrey    Background = 40
	BgMagenta Background = 45
	BgRed     Background = 41
	BgWhite   Background = 47
	BgYellow  Background = 43
)

const Reset = "\033[0m"

func ColoredWithBackground(text string, color Color, background Background, attributes ...Attribute) (ret string) {
	if os.Getenv("ANSI_COLORS_DISABLED") == "" {
		fmtStr := "\033[%dm"

		if background != -1 {
			ret += fmt.Sprintf(fmtStr, background)
		}

		for _, attr := range attributes {
			ret += fmt.Sprintf(fmtStr, attr)
		}

		ret += fmt.Sprintf(fmtStr, color) + text + Reset
	}

	return
}

func Colored(text string, color Color, attributes ...Attribute) string {
	return ColoredWithBackground(text, color, -1, attributes...)
}
