package badge

import "fmt"

type Item struct {
	Name       string
	Badge      string
	IsSection  bool
	IsSelected bool
}

func (i Item) Title() string {
	if i.IsSelected {
		return fmt.Sprintf("%s\t%s", "-->", i.Name)
	}
	return i.Name
}
func (i Item) Description() string { return i.Badge }
func (i Item) FilterValue() string { return i.Name }
