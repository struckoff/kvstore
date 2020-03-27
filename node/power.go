package node

type Power struct {
	p float64
}

func (p Power) Get() float64 {
	return p.p
}

func NewPower(p float64) Power {
	return Power{p}
}
