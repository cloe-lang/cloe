package compile

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompile(t *testing.T) {
	f, err := ioutil.TempFile("", "foo.tisp")
	assert.Nil(t, err)

	f.WriteString(`(write "Hello, world!")`)

	err = f.Close()
	assert.Nil(t, err)

	es, err := Compile(f.Name())

	assert.Nil(t, err)
	assert.Equal(t, 1, len(es))
}

func TestCompileSourceOfInvalidSyntax(t *testing.T) {
	f, err := ioutil.TempFile("", "foo.tisp")
	assert.Nil(t, err)

	f.WriteString(`(write "Hello, world!"`)

	err = f.Close()
	assert.Nil(t, err)

	_, err = Compile(f.Name())

	assert.NotNil(t, err)
}

func TestCompileWithInvalidPath(t *testing.T) {
	_, err := Compile("I'm the invalid path.")
	assert.NotNil(t, err)
}

func TestCompileStdin(t *testing.T) {
	es, err := Compile("")
	assert.Nil(t, err)
	assert.Zero(t, len(es))
}
