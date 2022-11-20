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
