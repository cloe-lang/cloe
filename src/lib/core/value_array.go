package core

const valueArrayLength = 8

// ValueArray is an array of values.
type ValueArray struct {
	values [valueArrayLength]Value
	start  int
	end    int
}

// NewValueArray creates a new array.
func NewValueArray(vs [8]Value) ValueArray {
	e := -1

	for _, v := range vs {
		if v != nil {
			e++
		}
	}

	return ValueArray{vs, 0, e}
}

// Append appends a value to an array.
func (a *ValueArray) Append(v Value) bool {
	if a.values[valueArrayLength-1] != nil {
		return false
	}

	a.end++
	a.values[a.end] = v

	return true
}

// Next returns a next value inside an array.
func (a *ValueArray) Next() Value {
	if a.start == valueArrayLength {
		return nil
	}

	v := a.values[a.start]
	a.values[a.start] = nil
	a.start++
	return v
}

// Slice returns values in an array as a slice.
func (a ValueArray) Slice() []Value {
	return a.values[a.start : a.end+1]
}

// Empty returns true if an array is empty, or false otherwise.
func (a ValueArray) Empty() bool {
	return a.end < a.start
}

// Merge merges 2 arrays into one.
func (a ValueArray) Merge(aa ValueArray) (ValueArray, []Value) {
	new := NewValueArray([8]Value{})
	vs := append(a.Slice(), aa.Slice()...)
	var start int

	for i, v := range vs {
		if !new.Append(v) {
			start = i
			break
		}
	}

	return new, vs[start:]
}
