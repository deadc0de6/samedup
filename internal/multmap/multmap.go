/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2023, deadc0de6
*/

package multmap

// MultMap a map of key-value
type MultMap[K comparable] map[K][]interface{}

// Add add a new entry in the multmap
func (s MultMap[K]) Add(key K, values ...interface{}) []interface{} {
	_, ok := s[key]
	if !ok {
		s[key] = values
	} else {
		s[key] = append(s[key], values...)
	}
	return s[key]
}

// GetAllValues return all values
func (s MultMap[K]) GetAllValues() []interface{} {
	var all []interface{}
	for _, v := range s {
		all = append(all, v...)
	}
	return all
}

// New new multmap
func New[K comparable]() MultMap[K] {
	return MultMap[K]{}
}
