package core

type Tenant struct {
	Organization string
	Project      string
}

func (t *Tenant) ID() string {
	return t.Organization + t.Project
}

func NewTenant() Tenant {
	return Tenant{
		Organization: "org",
		Project:      "proj",
	}
}
