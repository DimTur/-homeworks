package main

import "fmt"

func summaryTable(database Database, subject string) {
	studentsData := make(map[int]Students)
	objectsData := make(map[int]Objects)
	subjectsGrade := make(map[string]map[int][]int)

	for _, student := range database.Students {
		studentsData[student.Id] = student
	}

	for _, object := range database.Objects {
		objectsData[object.Id] = object
		subjectsGrade[object.Name] = make(map[int][]int)
	}

	for _, result := range database.Results {
		student := studentsData[result.Student_id]
		object := objectsData[result.Object_id]
		grades := subjectsGrade[object.Name][student.Grade]
		grades = append(grades, result.Result)
		subjectsGrade[object.Name][student.Grade] = grades
	}

	filteredGrades := Filter(subject, subjectsGrade)

	var allGrades []int

	fmt.Printf("____________________\n")
	fmt.Printf("%-13s | Mean\n", subject)
	fmt.Printf("____________________\n")

	for grade, grades := range filteredGrades {

		mean := Mean(grades)
		allGrades = append(allGrades, grades...)

		fmt.Printf("%-7d grade | %.1f\n", grade, mean)
	}
	totalMean := Mean(allGrades)
	fmt.Printf("____________________\n")
	fmt.Printf("mean          | %.1f\n", totalMean)
	fmt.Printf("____________________\n")
}
