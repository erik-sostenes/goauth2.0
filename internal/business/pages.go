// Package business is the business layer that defines the business rules
package business

import (
	"bytes"
	"context"
	"math/rand"

	"github.com/erik-sostenes/auth-api/internal/repository"

	"golang.org/x/oauth2"
)

type PageProvider interface {
	ProvidePage(context.Context) (*bytes.Buffer, string, error)
}

func NewPageProvider(oauthConfig *oauth2.Config, drawer repository.PageDrawer[repository.TemplateProvider]) PageProvider {
	if oauthConfig == nil {
		panic("missing oauthConfig dependency")
	}

	if drawer == nil {
		panic("missing drawer dependency")
	}

	return &pageProvider{
		oauthConfig: oauthConfig,
		drawer:      drawer,
	}
}

type pageProvider struct {
	oauthConfig *oauth2.Config
	drawer      repository.PageDrawer[repository.TemplateProvider]
}

func (p *pageProvider) ProvidePage(ctx context.Context) (page *bytes.Buffer, state string, err error) {
	state = p.generateState()

	template, err := p.drawer.DrawPage(ctx)
	if err != nil {
		return
	}

	url := p.oauthConfig.AuthCodeURL(state)

	page = bytes.NewBufferString("")

	err = template.Execute(page, map[string]any{"Url": url})
	if err != nil {
		return
	}

	return page, state, nil
}

func (p pageProvider) generateState() string {
	const (
		charset      = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-._~"
		stringLength = 64
	)

	state := make([]rune, stringLength)

	for i := range state {
		state[i] = (rune)(charset[rand.Intn(len(charset))])
	}

	return string(state)
}
