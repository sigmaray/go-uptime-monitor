package models

type CreateUserInput struct {
	Username        string `validate:"required,min=1,max=100" label:"login"`
	Password        string `validate:"required,min=1" label:"password"`
	ConfirmPassword string `validate:"required,eqfield=Password" label:"confirm password"`
}
