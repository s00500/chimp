package store

import "testing"

type TestStruct struct {
	Integer int
}

func (t TestStruct) Initialize() TestStruct {
	t.Integer = 1
	return t
}

func Test_ReadWhileMutate(t *testing.T) {
	m := Lockable[TestStruct]{}

	m.Mutate(func(state *TestStruct) {
		_, drop := m.ReadRef()
		defer drop()
	})

}
