package validator

import (
	"testing"
)

type (
	I interface {
		Foo() string
	}

	Impl struct {
		F string `valid:"len(3)"`
	}

	SubTest struct {
		Test string `valid:"required"`
	}

	TestString struct {
		BlankTag  string `valid:""`
		Required  string `valid:"required"`
		Len       string `valid:"len(10)"`
		Min       string `valid:"min(1)"`
		Max       string `valid:"max(10)"`
		MinMax    string `valid:"min(1),max(10)"`
		Lt        string `valid:"max(9)"`
		Lte       string `valid:"max(10)"`
		Gt        string `valid:"min(11)"`
		Gte       string `valid:"min(10)"`
		OmitEmpty string `valid:"optional,min(1),max(10)"`
		Sub       *SubTest
		SubIgnore *SubTest `valid:"-"`
		Anonymous struct {
			A string `valid:"required"`
		}
		Iface I
	}
)

func (i *Impl) Foo() string {
	return i.F
}

func BenchmarkValidateVarSuccess(b *testing.B) {
	v := New()

	s := "1"

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if err := v.ValidateVar(&s, "len(1)"); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkValidateVarParallelSuccess(b *testing.B) {
	v := New()

	s := "1"

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if err := v.ValidateVar(&s, "len(1)"); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkValidateStructSuccess(b *testing.B) {
	v := New()

	type Foo struct {
		StringValue string `valid:"len(5|10)"`
		IntValue    int    `valid:"len(5|10)"`
	}

	s := &Foo{StringValue: "Foobar", IntValue: 7}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if err := v.ValidateStruct(s); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkValidateStructParallelSuccess(b *testing.B) {
	v := New()

	type Foo struct {
		StringValue string `valid:"len(5|10)"`
		IntValue    int    `valid:"len(5|10)"`
	}

	s := &Foo{StringValue: "Foobar", IntValue: 7}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if err := v.ValidateStruct(s); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkValidateStructComplexSuccess(b *testing.B) {
	v := New()

	s := &TestString{
		Required:  "12345",
		Len:       "1234567890",
		Min:       "1",
		Max:       "1234567890",
		MinMax:    "12345",
		Lt:        "123456789",
		Lte:       "1234567890",
		Gt:        "12345678901",
		Gte:       "1234567890",
		OmitEmpty: "",
		Sub: &SubTest{
			Test: "1",
		},
		SubIgnore: &SubTest{
			Test: "",
		},
		Anonymous: struct {
			A string `valid:"required"`
		}{
			A: "1",
		},
		Iface: &Impl{
			F: "123",
		},
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if err := v.ValidateStruct(s); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkValidateStructComplexParallelSuccess(b *testing.B) {
	v := New()

	s := &TestString{
		Required:  "12345",
		Len:       "1234567890",
		Min:       "1",
		Max:       "1234567890",
		MinMax:    "12345",
		Lt:        "123456789",
		Lte:       "1234567890",
		Gt:        "12345678901",
		Gte:       "1234567890",
		OmitEmpty: "",
		Sub: &SubTest{
			Test: "1",
		},
		SubIgnore: &SubTest{
			Test: "",
		},
		Anonymous: struct {
			A string `valid:"required"`
		}{
			A: "1",
		},
		Iface: &Impl{
			F: "123",
		},
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if err := v.ValidateStruct(s); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkValidateVarFailure(b *testing.B) {
	v := New()

	s := "12"

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if err := v.ValidateVar(&s, "len(1)"); err == nil {
			b.Fatal("want invalid argument error, but got nil")
		}
	}
}

func BenchmarkValidateStructComplexFailure(b *testing.B) {
	v := New()

	s := &TestString{
		Required:  "",
		Len:       "",
		Min:       "",
		Max:       "12345678901",
		MinMax:    "",
		Lt:        "1234567890",
		Lte:       "12345678901",
		Gt:        "1",
		Gte:       "1",
		OmitEmpty: "12345678901",
		Sub: &SubTest{
			Test: "",
		},
		SubIgnore: &SubTest{
			Test: "",
		},
		Anonymous: struct {
			A string `valid:"required"`
		}{
			A: "",
		},
		Iface: &Impl{
			F: "12",
		},
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if err := v.ValidateStruct(s); err == nil {
			b.Fatal("want validation errors, but got nil")
		}
	}
}

func BenchmarkValidateStructComplexParallelFailure(b *testing.B) {
	v := New()

	s := &TestString{
		Required:  "",
		Len:       "",
		Min:       "",
		Max:       "12345678901",
		MinMax:    "",
		Lt:        "1234567890",
		Lte:       "12345678901",
		Gt:        "1",
		Gte:       "1",
		OmitEmpty: "12345678901",
		Sub: &SubTest{
			Test: "",
		},
		SubIgnore: &SubTest{
			Test: "",
		},
		Anonymous: struct {
			A string `valid:"required"`
		}{
			A: "",
		},
		Iface: &Impl{
			F: "12",
		},
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if err := v.ValidateStruct(s); err == nil {
				b.Fatal("want validation errors, but got nil")
			}
		}
	})
}
