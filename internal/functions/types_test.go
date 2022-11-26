// Copyright 2022 Heath Stewart.
// Licensed under the MIT License. See LICENSE.txt in the project root for license information.

package functions

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFromDefaultValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		value   any
		want    string
		wantErr bool
	}{
		{
			name:  "empty string",
			value: "",
			want:  "",
		},
		{
			name:  "string",
			value: "value",
			want:  "value",
		},
		{
			name:  "int",
			value: 1,
			want:  "1",
		},
		{
			name:  "boolean (true)",
			value: true,
			want:  "true",
		},
		{
			name:  "boolean (false)",
			value: false,
			want:  "",
		},
		{
			name:    "unsupported",
			value:   time.Now,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := fromDefaultValue(tt.value)
			if tt.wantErr {
				assert.Errorf(t, err, "unsupported type")
				return
			}

			if !assert.NoError(t, err) {
				return
			}

			assert.Equal(t, tt.want, v.String())
		})
	}
}

func TestIntValue_Format(t *testing.T) {
	t.Parallel()

	v := intValue{1}
	got, ok := v.Format("2")
	assert.True(t, ok)
	assert.Equal(t, "2", got)

	_, ok = v.Format("invalid")
	assert.False(t, ok)
}

func TestIntValue_String(t *testing.T) {
	t.Parallel()

	v := intValue{1}
	assert.Equal(t, "1", v.String())
}

func TestBoolValue_Display(t *testing.T) {
	t.Parallel()

	v := boolValue{true}
	assert.Equal(t, "Y/n", v.Display())

	v.defaultValue = false
	assert.Equal(t, "y/N", v.Display())
}

func TestBoolValue_Format(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		defaultValue bool
		value        string
		want         string
		isInvalid    bool
	}{
		{
			name:         "empty (true)",
			defaultValue: true,
			want:         "true",
		},
		{
			name:         "empty (false)",
			defaultValue: false,
			want:         "",
		},
		{
			name:  "y",
			value: "y",
			want:  "true",
		},
		{
			name:  "yes",
			value: "yes",
			want:  "true",
		},
		{
			name:  "true",
			value: "true",
			want:  "true",
		},
		{
			name:  "n",
			value: "n",
		},
		{
			name:  "no",
			value: "no",
		},
		{
			name:  "false",
			value: "false",
		},
		{
			name:      "invalid",
			value:     "invalid",
			isInvalid: true,
		},
	}

	for _, tt := range tests {
		v := boolValue{tt.defaultValue}
		t.Run(tt.name, func(t *testing.T) {
			got, ok := v.Format(tt.value)
			if !assert.NotEqual(t, tt.isInvalid, ok) {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestBoolValue_String(t *testing.T) {
	t.Parallel()

	v := boolValue{true}
	assert.Equal(t, "true", v.String())

	v.defaultValue = false
	assert.Equal(t, "", v.String())
}
