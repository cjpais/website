package cjtime

type Year struct {
	String string
	Int int
	Months []*Month
}

type Month struct {
	Year *Year
	String string
	Int int
	Days []*Day
}

type Day struct {
	Month *Month
	String string
	Int int
	Posts []*Post
}
