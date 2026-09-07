package fixture

type finalizable struct{}

func (finalizable) FinalizeAt() {}

func bad() {
	var f finalizable
	f.FinalizeAt()
}
