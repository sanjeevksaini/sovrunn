package grantintent

type GrantPort struct{}

func (GrantPort) Submit() {}

func allowed() {
	var p GrantPort
	p.Submit()
}
