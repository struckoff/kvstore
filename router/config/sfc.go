package config

type SFCConfig struct {
	////Amount of space filling curve dimensions
	//Dimensions uint64 `envconfig:"KVROUTER_SFC_DIMENSIONS"`
	////Size of space filling curve
	//Size uint64 `envconfig:"KVROUTER_SFC_SIZE"`
	////Space filling curve type
	//Curve CurveType `envconfig:"KVROUTER_SFC_CURVE"`

	Dimensions uint64
	//Size of space filling curve
	Size uint64
	//Space filling curve type
	Curve CurveType
}
