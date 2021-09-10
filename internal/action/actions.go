package action

import (
	"sort"
)

type Actions []*Action

func (a Actions) Len() int {
	return len(a)
}

func (a Actions) Less(i, j int) bool {
	return a[i].action.Index < a[j].action.Index
}

func (a Actions) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a Actions) Sort() []*Action {
	sort.Sort(a)
	return a
}
