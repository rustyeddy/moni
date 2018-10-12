package store

// Index is a map of paths to Objects
type index map[string]*Object

func (i index) Find(idx string) (obj *Object) {
	return i[idx]
}

func (i index) Exists(idx string) bool {
	if obj := i.Find(idx); obj == nil {
		return false
	}
	return true
}

func (i index) Get(idx string) *Object {
	return i.Find(idx)
}

func (i index) Set(idx string, obj *Object) {
	i[idx] = obj
}

func (i index) Len() int {
	return len(i)
}

func (i index) NameObjects() (names []string, objects []*Object) {
	count := len(i)
	if count < 1 {
		return nil, nil
	}
	names = make([]string, count)
	objects = make([]*Object, count)
	for n, obj := range i {
		names = append(names, n)
		objects = append(objects, obj)
	}
	return names, objects
}
