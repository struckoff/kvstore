package nodes

type Capacity struct {
	c float64
}

func (c Capacity) Get() float64 {
	return c.c
}

func NewCapacity(c float64) Capacity {
	return Capacity{c}
}
