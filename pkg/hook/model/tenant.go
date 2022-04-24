package model


type Tenant struct {
	Config  Configuration
	Environ map[string]string
}
