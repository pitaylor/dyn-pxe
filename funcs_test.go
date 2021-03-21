package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShortHash(t *testing.T) {
	assert.Equal(t, "e3b0c44", ShortHash(""))
	assert.Equal(t, "ade0997", ShortHash("XYZ"))
}
