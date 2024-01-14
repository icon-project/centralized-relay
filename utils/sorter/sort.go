package sorter

import "sort"

// GenericSorter is a generic sorter for slices of any comparable type.
type GenericSorter[T any] struct {
	Items []T
	By    func(p1, p2 T) bool
}

// Len is part of sorter.Interface.
func (gs *GenericSorter[T]) Len() int {
	return len(gs.Items)
}

// Swap is part of sorter.Interface.
func (gs *GenericSorter[T]) Swap(i, j int) {
	gs.Items[i], gs.Items[j] = gs.Items[j], gs.Items[i]
}

// Less is part of sorter.Interface. It calls the "by" closure in the sorter.
func (gs *GenericSorter[T]) Less(i, j int) bool {
	return gs.By(gs.Items[i], gs.Items[j])
}

// Sort sorts the argument slice according to the provided "less" function.
func Sort[T any](items []T, by func(p1, p2 T) bool) {
	gs := &GenericSorter[T]{
		Items: items,
		By:    by,
	}
	sort.Sort(gs)
}
