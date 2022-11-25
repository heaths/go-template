// Copyright 2022 Heath Stewart.
// Licensed under the MIT License. See LICENSE.txt in the project root for license information.

package template

import (
	"testing"

	"github.com/heaths/go-template/internal/processor"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

func TestWithLanguage(t *testing.T) {
	p := new(processor.Processor)
	WithLanguage(language.English)(p)
	p.Initialize()

	assert.Equal(t, language.English, *p.Language)
}

func TestWithDelims(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		leftDelim  string
		rightDelim string
		wantPanic  bool
	}{
		{
			name: "neither",
		},
		{
			name:       "both",
			leftDelim:  "<%",
			rightDelim: "%>",
		},
		{
			name:      "left",
			leftDelim: "<%",
			wantPanic: true,
		},
		{
			name:       "right",
			rightDelim: "%>",
			wantPanic:  true,
		},
	}

	p := new(processor.Processor)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				assert.Panics(t, func() {
					WithDelims(tt.leftDelim, tt.rightDelim)
				})
				return
			}

			WithDelims(tt.leftDelim, tt.rightDelim)(p)
			assert.Equal(t, tt.leftDelim, p.LeftDelim)
			assert.Equal(t, tt.rightDelim, p.RightDelim)
		})
	}
}
