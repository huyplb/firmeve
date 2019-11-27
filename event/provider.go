package event

import (
	"github.com/firmeve/firmeve"
	"github.com/firmeve/firmeve/container"
)

type Provider struct {
	firmeve.BaseProvider
}

func (p *Provider) Name() string {
	return `event`
}

func (p *Provider) Register() {
	p.Firmeve.Bind(`event`, New(), container.WithShare(true))
}

func (p *Provider) Boot() {

}
