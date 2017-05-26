package main

import (
	"os"
	"testing"

	"github.com/mattn/go-shellwords"
	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	ss, err := shellwords.NewParser().Parse(os.Getenv("ARGS"))
	assert.Equal(t, nil, err)

	os.Args = append(os.Args[:1], ss...)
	main()
}
