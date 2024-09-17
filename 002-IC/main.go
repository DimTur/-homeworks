package main

import (
	"errors"
	"fmt"
)

type Validator interface {
	Validate() error
}

type MatrixValidator struct {
	matrix [][]int
}

type UserAnswerValidator struct {
	userAnswer []int
	matrix     [][]int
}

// Function checks if the matrix is empty,
// if it is square, if no loops, and if the matrix is symmetric
func (mv MatrixValidator) Validate() error {
	n := len(mv.matrix)
	if n == 0 {
		return errors.New("matrix is empty")
	}

	hasPath := false

	for i := 0; i < n; i++ {
		if len(mv.matrix[i]) != n {
			return errors.New("matrix is not square")
		}
		if mv.matrix[i][i] != 0 {
			return errors.New("matrix has loop")
		}
		for j := 0; j < n; j++ {
			if mv.matrix[i][j] < 0 {
				return errors.New("matrix contains negative values")
			}
			if mv.matrix[i][j] > 0 {
				hasPath = true
			}
		}
	}

	if !hasPath {
		return errors.New("matrix contains no valid paths")
	}

	return nil
}

// Function validates a user's answer
// by ensuring it does not contain any duplicate values.
func (uav UserAnswerValidator) Validate() error {
	n := len(uav.userAnswer)
	m := len(uav.matrix)
	if n == 0 {
		return nil
	}

	validMap := make(map[int]interface{})

	for i := 0; i < n; i++ {
		if _, exists := validMap[uav.userAnswer[i]]; exists {
			return errors.New("the answer cannot contain duplicate values")
		}
		if uav.userAnswer[i] >= m || uav.userAnswer[i] < 0 {
			return errors.New("the answer cannot contain value out of matrix range")
		}
		validMap[uav.userAnswer[i]] = struct{}{}
	}

	return nil
}

func EvalSequence(matrix [][]int, userAnswer []int) int {
	mv := MatrixValidator{matrix: matrix}
	if err := mv.Validate(); err != nil {
		fmt.Println("Matrix incorrect:", err)
		return -1
	}

	uav := UserAnswerValidator{userAnswer: userAnswer, matrix: matrix}
	if err := uav.Validate(); err != nil {
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
	maxPathWeight := 0

	// calculate max path weight for each vertex
	for i := 0; i < len(matrix); i++ {
		if visited[i] {
			continue
		}
		pathWeight := dFSMaxUtil(matrix, i, visited)

		// update max path weight
		if pathWeight > maxPathWeight {
			maxPathWeight = pathWeight
		}
	}

	return maxPathWeight
}

// DFS helper function to calculate maximum path weight
func dFSMaxUtil(matrix [][]int, vertex int, visited []bool) int {

	visited[vertex] = true

	// Find all vertices adjacent to the current vertex
	maxWeight := 0
	for i := 0; i < len(matrix); i++ {
		if matrix[vertex][i] > 0 && !visited[i] {
			weight := matrix[vertex][i] + dFSMaxUtil(matrix, i, visited)
			if weight > maxWeight {
				maxWeight = weight
			}
		}
	}

	// Store the maximum weight of the path from the current vertex
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
		{0, 0, 0},
		{0, 0, 0},
		{0, 0, 0},
	}
	ua := []int{0, 1, 2}

	fmt.Println(EvalSequence(mtx1, ua))
}
