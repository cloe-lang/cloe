package compile

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/cloe-lang/cloe/src/lib/consts"
	"github.com/stretchr/testify/assert"
)

func TestCompile(t *testing.T) {
	for _, s := range []string{
		`(print "Hello, world!")`,
		`(import "http") (print (http.get "http://httpbin.org"))`,
	} {
		f, err := ioutil.TempFile("", "")
		assert.Nil(t, err)

		f.WriteString(s)

		err = f.Close()
		assert.Nil(t, err)

		es, err := Compile(f.Name())

		assert.Nil(t, err)
		assert.Equal(t, 1, len(es))
	}
}

func TestCompileSourceOfInvalidSyntax(t *testing.T) {
	f, err := ioutil.TempFile("", "")
	assert.Nil(t, err)

	f.WriteString(`(print "Hello, world!"`)

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

func TestCompileWithSubModule(t *testing.T) {
	m := createModuleScript(t)

	f, err := ioutil.TempFile("", "")
	assert.Nil(t, err)

	f.WriteString(fmt.Sprintf(`(import "%v") (print (%v.hello "John"))`, m, path.Base(m)))

	err = f.Close()
	assert.Nil(t, err)

	es, err := Compile(f.Name())

	assert.Nil(t, err)
	assert.Equal(t, 1, len(es))
}

func createModuleScript(t *testing.T) string {
	f, err := ioutil.TempFile("", "module")
	assert.Nil(t, err)

	f.WriteString(`(def (hello name) (merge "Hello, " name "!"))`)

	err = f.Close()
	assert.Nil(t, err)

	err = os.Rename(f.Name(), f.Name()+consts.FileExtension)
	assert.Nil(t, err)

	return f.Name()
}
