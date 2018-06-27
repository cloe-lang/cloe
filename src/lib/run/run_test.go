package run

import (
	"errors"
	"os"
	"sync"
	"testing"

	"github.com/cloe-lang/cloe/src/lib/builtins"
	"github.com/cloe-lang/cloe/src/lib/compile"
	"github.com/cloe-lang/cloe/src/lib/core"
	"github.com/stretchr/testify/assert"
)

func TestRunWithNoThunk(t *testing.T) {
	Run(nil)
}

func TestRunWithOneThunk(t *testing.T) {
	Run([]compile.Effect{compile.NewEffect(core.PApp(builtins.Print, core.NewNumber(42)), false)})
}

func TestRunWithThunks(t *testing.T) {
	o := compile.NewEffect(core.PApp(builtins.Print, core.NewNumber(42)), false)
	Run([]compile.Effect{o, o, o, o, o, o, o, o})
}

func TestRunWithExpandedList(t *testing.T) {
	Run([]compile.Effect{compile.NewEffect(
		core.NewList(core.PApp(builtins.Print, core.True), core.PApp(builtins.Print, core.False)),
		true)})
}

func TestEvalEffectListFail(t *testing.T) {
	defer func() {
		assert.NotNil(t, recover())
	}()

	wg := sync.WaitGroup{}
	wg.Add(1)
	evalEffectList(core.DummyError, &wg, func(err error) { panic(err) })
}

func TestRunEffectFail(t *testing.T) {
	defer func() {
		assert.NotNil(t, recover())
	}()

	wg := sync.WaitGroup{}
	sem <- true
	wg.Add(1)
	runEffect(core.Nil, &wg, func(err error) { panic(err) })
}

func TestFail(t *testing.T) {
	s := 0

	fail := failWithExit(func(i int) { s = i })

	fail(errors.New("foo"))

	assert.Equal(t, 1, s)
}

func TestFailPanic(t *testing.T) {
	defer func() {
		assert.NotNil(t, recover())
	}()

	os.Stderr.Close()
	fail := failWithExit(func(int) {})
	fail(errors.New("foo"))
}
