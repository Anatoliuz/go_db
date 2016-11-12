package main


type LastPathASC []Post
type LastPathDESC []Post


func (a LastPathASC) Len() int           { return len(a) }
func (a LastPathASC) Swap(i, j int) {
		a[i], a[j] = a[j], a[i]
}
func (a LastPathASC) Less(i, j int) bool {
	return a[i].LastPath < a[j].LastPath
}
func (a LastPathDESC) Len() int           { return len(a) }
func (a LastPathDESC) Swap(i, j int) {
		a[i], a[j] = a[j], a[i]
}
func (a LastPathDESC) Less(i, j int) bool {
	return a[i].LastPath > a[j].LastPath
}