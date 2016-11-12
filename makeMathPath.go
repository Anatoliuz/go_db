package main
import "strconv"


const sizeOfPath int = 3

func intCapacity(num  int) int {
	size := 0
	for num > 0 {
		num = num / 10
		size++
	}
	return size
}

func  makeMathPathBetweenDots(number int) string           {
	var mathPath string;
	for i := sizeOfPath - intCapacity(number); i > 0 ; i--  {
		mathPath += "0"
	}
	numStr := strconv.Itoa(number)
	mathPath += numStr
	return mathPath;
}