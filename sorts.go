package main

type FirstPathASC []Post
type FirstPathDESC []Post

type LastPathASC []Post
type LastPathDESC []Post

func (a FirstPathASC) Len() int           { return len(a) }
func (a FirstPathASC) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a FirstPathASC) Less(i, j int) bool { return a[i].FirstPath < a[j].FirstPath}

func (a FirstPathDESC) Len() int           { return len(a) }
func (a FirstPathDESC) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a FirstPathDESC) Less(i, j int) bool { return a[i].FirstPath > a[j].FirstPath}

func (a LastPathASC) Len() int           { return len(a) }
func (a LastPathASC) Swap(i, j int) {
	if a[i].FirstPath == a[j].FirstPath {
		a[i], a[j] = a[j], a[i]
	}
}
func (a LastPathASC) Less(i, j int) bool {
	return a[i].LastPath < a[j].LastPath
}
func (a LastPathDESC) Len() int           { return len(a) }
func (a LastPathDESC) Swap(i, j int) {
	if a[i].FirstPath == a[j].FirstPath {
		a[i], a[j] = a[j], a[i]
	}
}
func (a LastPathDESC) Less(i, j int) bool {
	return a[i].LastPath > a[j].LastPath
}

//func sort_tree(posts []Post, str string) {
//	if str == "ASC" {
//		sort.Sort(FirstPathASC(posts))
//		sort.Sort(LastPathASC(posts))
//	} else {
//		sort.Sort(FirstPathDESC(posts))
//		sort.Sort(LastPathDESC(posts))
//	}
//}