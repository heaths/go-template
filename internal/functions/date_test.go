// Copyright 2022 Heath Stewart.
// Licensed under the MIT License. See LICENSE.txt in the project root for license information.

package functions

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDateFunc(t *testing.T) {
	t.Parallel()

	d := DateFunc()
	assert.Equal(t, d.t.UTC(), d.t)
}

func TestDate_Year(t *testing.T) {
	t.Parallel()

	d := DateFunc()
	assert.Equal(t, d.t.Year(), d.Year())
}

func TestDate_Local(t *testing.T) {
	t.Parallel()

	d := DateFunc().Local()
	assert.Equal(t, d.t.Local(), d.t)
}

func TestDate_Format(t *testing.T) {
	t.Parallel()

	d := DateFunc()
	assert.Equal(t, fmt.Sprint(d.t.Year()), d.Format("2006"))
}
