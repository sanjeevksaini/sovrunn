package fixture

type Port struct{}

func (Port) EvaluatePrivileged() {}

func bad() {
	var p Port
	p.EvaluatePrivileged()
}
