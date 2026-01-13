package engine

type Cluster struct {
	Name string `json:"name"`
}

type Instance struct {
	Key  string `json:"key"`
	Role string `json:"role"`
}

type Topology struct {
	Clusters  []Cluster  `json:"clusters"`
	Instances []Instance `json:"instances"`
}

type Engine interface {
	Clusters() ([]Cluster, error)
	Instances() ([]Instance, error)
	Topology() (Topology, error)
}
