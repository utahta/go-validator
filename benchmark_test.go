package validator

import (
	"testing"
)

func BenchmarkValidateVarSuccess(b *testing.B) {
	v := New()

	s := "1"

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		v.ValidateVar(&s, "len(1)")
	}
}

func BenchmarkValidateStructSuccess(b *testing.B) {
	v := New()

	type Foo struct {
		StringValue string `valid:"len(5|10)"`
		IntValue    int    `valid:"len(5|10)"`
	}

	validFoo := &Foo{StringValue: "Foobar", IntValue: 7}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		v.ValidateStruct(validFoo)
	}
}
