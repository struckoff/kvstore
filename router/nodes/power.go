package nodes

//Power represents node weigh
type Power struct {
	p float64
}

//Get returns power value
func (p Power) Get() float64 {
	return p.p
}

//NewPower returns power instance
func NewPower(p float64) Power {
	return Power{p}
}
