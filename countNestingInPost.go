package main

func countNesting(mathPath string) int{
	nestingNum := 0

	for i := 0; i < len(mathPath) ; i++ {
		if mathPath[i] == '.' {
			nestingNum++
		}
	}
	return  nestingNum + 1
}