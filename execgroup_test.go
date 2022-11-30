package execgroup

import (
	"fmt"
	"testing"
	"time"
)

const delay = 200 * time.Millisecond

func TestExecGroup(t *testing.T) {
	var a, b, c int

	var eg ExecGroup

	eg.Do(func() {
		a = 1
		b = 1
		c = 1
		time.Sleep(delay)
		a++
	})

	eg.Do(func() {
		a = 1
		b = 1
		c = 1
		time.Sleep(delay)
		b++
	})

	eg.Do(func() {
		a = 1
		b = 1
		c = 1
		time.Sleep(delay)
		c++
	})

	eg.Wait()

	if a != 2 && b != 2 && c != 2 {
		t.Errorf("Something went wrong: a=%d, b=%d, c=%d", a, b, c)
	}
}

func TestExecGroup_recursive(t *testing.T) {
	var a, b, c int
	var eg ExecGroup

	eg.Do(func() {
		a++
		eg.Do(func() {
			time.Sleep(delay)
			b++
			eg.Do(func() {
				time.Sleep(delay)
				c++
			})
			b++
		})
		a++
	})

	eg.Wait()

	if a != 2 && b != 2 && c != 1 {
		t.Errorf("Something went wrong: a=%d, b=%d, c=%d", a, b, c)
	}
}

func TestExecGroup_panic(t *testing.T) {
	var eg ExecGroup
	eg.Do(func() {
		panic(fmt.Errorf("Fmt Error"))
	})
	eg.Do(func() {
		time.Sleep(100 * time.Millisecond)
		panic("string error")
	})
	err := eg.Wait().(MultiError)

	if err.Len() != 2 {
		t.Errorf("expected 2 errors, saw %d", err.Len())
	}

	e0 := err.Errors()[0].Error()
	a0 := "Fmt Error"
	if e0 != a0 {
		t.Errorf("wanted '%s' was '%s'", a0, e0)
	}

	e1 := err.Errors()[1].Error()
	a1 := "panic in execgroup: string error"
	if e1 != a1 {
		t.Errorf("wanted '%s' was '%s'", a1, e1)
	}
}
