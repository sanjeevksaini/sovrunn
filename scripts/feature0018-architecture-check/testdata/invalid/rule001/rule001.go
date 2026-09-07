package fixture

type tx struct{}

func (tx) Commit() {}

func bad() {
	var t tx
	t.Commit()
}
