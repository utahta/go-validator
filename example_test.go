package validator_test

import (
	"fmt"

	"github.com/utahta/go-validator"
)

func ExampleValidateStruct_simpleSuccess() {
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

	// Output:
	// <nil>
}

func ExampleValidateStruct_simpleFailure() {
	type user struct {
		Name  string `valid:"required,alphanum"`
		Age   uint   `valid:"required,len(0|116)"`
		Email string `valid:"optional,email"`
	}

	v := validator.New()
	err := v.ValidateStruct(&user{
		Name:  "_",
		Age:   200,
		Email: "invalid",
	})

	fmt.Println(err)

	// Output:
	// Name: '_' does validate as 'alphanum';Age: '200' does validate as 'len(0|116)';Email: 'invalid' does validate as 'email'
}
