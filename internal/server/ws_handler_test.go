package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	fnName, arguments, err := parseMessage("fn1|a=1|b=")
	assert.Nil(t, err)
	assert.Equal(t, "fn1", fnName)
	assert.Equal(t, "1", arguments["a"])
	assert.Equal(t, "", arguments["b"])
}
