package main

import (
	"errors"
	"fmt"
)

// The mtxValidator function checks if the matrix is empty,
// if it is square, if no loops, and if the matrix is symmetric
func mtxValidator(matrix [][]int) error {
	n := len(matrix)
	if n == 0 {
		return errors.New("matrix is empty")
	}

	for i := 0; i < n; i++ {
		if len(matrix[i]) != n {
			return errors.New("matrix is not square")
		}
		if matrix[i][i] != 0 {
			return errors.New("matrix has loop")
		}
		for j := 0; j < n; j++ {
			if matrix[i][j] != matrix[j][i] {
				return errors.New("matrix is not symmetric")
			}
		}
	}
	return nil
}

// The uAnswerValidator function validates a user's answer
// by ensuring it does not contain any duplicate values.
func uAnswerValidator(userAnswer []int, matrix [][]int) error {
	n := len(userAnswer)
	m := len(matrix)
	if n == 0 {
		return nil
	}

	validMap := make(map[int]bool)
	validMap[userAnswer[0]] = true

	for i := 1; i < n; i++ {
		if validMap[userAnswer[i]] {
			return errors.New("the answer cannot contain duplicate values")
		}
		if userAnswer[i] >= m || userAnswer[i] < 0 {
			return errors.New("the answer cannot contain value out of matrix range")
		}
		validMap[userAnswer[i]] = true
	}

	return nil
}

func EvalSequence(matrix [][]int, userAnswer []int) int {
	if err := mtxValidator(matrix); err != nil {
		fmt.Println("Matrix incorrect:", err)
		return -1
	}
	if err := uAnswerValidator(userAnswer, matrix); err != nil {
		fmt.Println("User answer incorrect:", err)
		return -2
	}

	maxGrade := calMaxGrade(matrix)
	userGrade := calcUserGrade(matrix, userAnswer)

	percent := userGrade * 100 / maxGrade

	return percent
}

// DFS max way in graph
func calMaxGrade(matrix [][]int) int {
	// init array to store max path weight for each vertex
	visited := make([]bool, len(matrix))
	maxGrade := make([]int, len(matrix))

	// calculate max path weight for each vertex
	maxPathWeight := 0
	for i := 0; i < len(matrix); i++ {
		if maxGrade[i] == 0 { // only for unvisited vertices
			dFSMaxUtil(matrix, i, visited, maxGrade)
		}

		// update max path weight
		if maxGrade[i] > maxPathWeight {
			maxPathWeight = maxGrade[i]
		}
	}

	return maxPathWeight
}

// DFS helper function to calculate maximum path weight
func dFSMaxUtil(matrix [][]int, vertex int, visited []bool, maxGrade []int) int {
	// If the maximum path weight has already been calculated, return it
	if maxGrade[vertex] != 0 {
		return maxGrade[vertex]
	}

	visited[vertex] = true

	// Find all vertices adjacent to the current vertex
	maxWeight := 0
	for i := 0; i < len(matrix); i++ {
		if matrix[vertex][i] > 0 && !visited[i] {
			weight := matrix[vertex][i] + dFSMaxUtil(matrix, i, visited, maxGrade)
			if weight > maxWeight {
				maxWeight = weight
			}
		}
	}

	// Store the maximum weight of the path from the current vertex
	maxGrade[vertex] = maxWeight
	visited[vertex] = false // Reset vivsit for other paths
	return maxWeight
}

// Function to calculate path weight for a custom sequence
func calcUserGrade(matrix [][]int, userAnswer []int) int {
	// Check that the custom sequence is not empty and has at least one vertex
	if len(userAnswer) == 0 {
		return 0
	}

	// Walk through the sequence and add up the weights of the edges
	userGrade := 0
	for i := 0; i < len(userAnswer)-1; i++ {
		fromVert := userAnswer[i]
		toVert := userAnswer[i+1]

		userGrade += matrix[fromVert][toVert]
	}

	return userGrade
}

func main() {
	mtx1 := [][]int{
		{0, 2, 3, 0, 0},
		{2, 0, 0, 1, 1},
		{3, 0, 0, 0, 0},
		{0, 1, 0, 0, 0},
		{0, 1, 0, 0, 0},
	}
	ua := []int{4, 1, 2, 5}

	fmt.Println(EvalSequence(mtx1, ua))
}
