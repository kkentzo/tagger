package watchers

type PathSet struct {
	paths map[string]bool
	// TODO: synchronize access using a mutex
}

func NewPathSet(paths []string) *PathSet {
	// convert array to map
	pathsM := make(map[string]bool)
	for _, name := range paths {
		pathsM[name] = true
	}

	return &PathSet{
		paths: pathsM,
	}
}

func (ps *PathSet) Has(path string) bool {
	return ps.paths[path]
}

func (ps *PathSet) Add(path string) {
	ps.paths[path] = true
}

func (ps *PathSet) Remove(path string) {
	delete(ps.paths, path)
}
