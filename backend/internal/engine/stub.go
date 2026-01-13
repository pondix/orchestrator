package engine

type Stub struct{}

func NewStub() *Stub {
	return &Stub{}
}

func (s *Stub) Clusters() ([]Cluster, error) {
	return []Cluster{}, nil
}

func (s *Stub) Instances() ([]Instance, error) {
	return []Instance{}, nil
}

func (s *Stub) Topology() (Topology, error) {
	return Topology{Clusters: []Cluster{}, Instances: []Instance{}}, nil
}
