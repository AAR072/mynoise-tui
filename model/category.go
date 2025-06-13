package model

type categoryItem struct{ name string }

func (c categoryItem) Title() string       { return c.name }
func (c categoryItem) Description() string { return "Category" }
func (c categoryItem) FilterValue() string { return c.name }
