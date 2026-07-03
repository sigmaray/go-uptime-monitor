package handlers

import (
	"html/template"

	"gorm.io/gorm"
)

type Handler struct {
	DB        *gorm.DB
	Templates *template.Template
}

func NewHandler(db *gorm.DB, tmpl *template.Template) *Handler {
	return &Handler{DB: db, Templates: tmpl}
}
