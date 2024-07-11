package main

func EvalSequence(matrix [][]int, userAnswer []int) int {

	// validation
	maxGrade := calMaxGrade(matrix)
	userGrade := calcUserGrade(matrix, userAnswer)

	percent := userGrade * 100 / maxGrade

	return percent
}

func calMaxGrade(matrix [][]int) int {
	return -1
}

func calcUserGrade(matrix [][]int, userAnswer []int) int {
	return 0
}
