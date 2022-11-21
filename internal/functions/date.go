// Copyright 2022 Heath Stewart.
// Licensed under the MIT License. See LICENSE.txt in the project root for license information.

package functions

import (
	"time"
)

func DateFunc() Date {
	return Date{
		t: time.Now().UTC(),
	}
}

type Date struct {
	t time.Time
}

func (d Date) Year() int {
	return d.t.Year()
}

func (d Date) Local() Date {
	return Date{
		t: d.t.Local(),
	}
}

func (d Date) Format(layout string) string {
	return d.t.Format(layout)
}
