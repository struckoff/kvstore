package nodes

type Capacity interface {
	Get() (float64, error)
}
