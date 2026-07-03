package models

type CreateUserInput struct {
	Username        string `form:"username" validate:"required,min=1,max=100" label:"login"`
	Password        string `form:"password" validate:"required,min=1" label:"password"`
	ConfirmPassword string `form:"confirm_password" validate:"required,eqfield=Password" label:"confirm password"`
}

type UpdateUserInput struct {
	Username        string `form:"username" validate:"required,min=1,max=100" label:"login"`
	Password        string `form:"password" validate:"omitempty,min=1" label:"password"`
	ConfirmPassword string `form:"confirm_password" validate:"eqfield=Password" label:"confirm password"`
}
