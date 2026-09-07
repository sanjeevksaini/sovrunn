package fixture

type GrantPort struct{}

func (GrantPort) Submit() {}

func bad() {
	var p GrantPort
	p.Submit()
}
