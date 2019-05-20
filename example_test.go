package validator_test

import (
	"context"
	"fmt"

	"github.com/utahta/go-validator"
)

func ExampleValidateStruct_simple() {
	type user struct {
		Name  string `valid:"required,alphanum"`
		Age   uint   `valid:"required,len(0|116)"`
		Email string `valid:"optional,email"`
	}

	v := validator.New()
	err := v.ValidateStruct(&user{
		Name:  "gopher",
		Age:   9,
		Email: "",
	})
	fmt.Println(err)

	err = v.ValidateStruct(&user{
		Name:  "_",
		Age:   200,
		Email: "invalid",
	})
	fmt.Println(err)

	// Output:
	// <nil>
	// Name: '_' does validate as 'alphanum';Age: '200' does validate as 'len(0|116)';Email: 'invalid' does validate as 'email'
}

func ExampleValidateStruct_setFunc() {
	type content struct {
		Type string `valid:"contentType(image/jpeg|image/png|image/gif)"`
	}

	v := validator.New(
		validator.WithFunc("contentType", func(_ context.Context, f validator.Field, opt validator.FuncOption) (bool, error) {
			v := f.Value().String()
			for _, param := range opt.TagParams {
				if v == param {
					return true, nil
				}
			}
			return false, nil
		}),
	)
	err := v.ValidateStruct(&content{
		Type: "image/jpeg",
	})
	fmt.Println(err)

	err = v.ValidateStruct(&content{
		Type: "image/bmp",
	})
	fmt.Println(err)

	// Output:
	// <nil>
	// Type: 'image/bmp' does validate as 'contentType(image/jpeg|image/png|image/gif)'
}

func ExampleValidateStruct_or() {
	type user struct {
		ID string `valid:"or(alpha|numeric)"`
	}

	v := validator.New()
	err := v.ValidateStruct(&user{
		ID: "abc",
	})
	fmt.Println(err)

	err = v.ValidateStruct(&user{
		ID: "123",
	})
	fmt.Println(err)

	err = v.ValidateStruct(&user{
		ID: "abc123",
	})
	fmt.Println(err)

	// Output:
	// <nil>
	// <nil>
	// ID: 'abc123' does validate as 'or(alpha|numeric)'
}
