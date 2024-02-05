package repository

import (
	"context"
	"embed"
	"html/template"
	"path"
)

const (
	viewsDir  = "templates/*"
	extension = "*.gohtml"
)

type PageDrawer[T TemplateProvider] interface {
	DrawPage(context.Context) (T, error)
}

func NewPageDrawer() PageDrawer[TemplateProvider] {
	return &pageDrawer{}
}

//go:embed templates/login/*
var templateLogin embed.FS

type TemplateProvider = *template.Template

type pageDrawer struct {
	TemplateProvider
}

func (d *pageDrawer) DrawPage(ctx context.Context) (TemplateProvider, error) {
	return d.ParseFS(templateLogin, path.Join(viewsDir, extension))
}
