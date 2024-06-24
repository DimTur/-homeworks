package main

import "fmt"

func summaryTable(database Database) {
	studentsData := make(map[int]Students)
	objectsData := make(map[int]Objects)
	subjectsGrade := make(map[int]map[int][]int)

	for _, student := range database.Students {
		studentsData[student.Id] = student
	}

	for _, object := range database.Objects {
		objectsData[object.Id] = object
		subjectsGrade[object.Id] = make(map[int][]int)
	}

	for _, result := range database.Results {
		student := studentsData[result.Student_id]
		object := objectsData[result.Object_id]

		grades := subjectsGrade[object.Id][student.Grade]
		grades = append(grades, result.Result)
		subjectsGrade[object.Id][student.Grade] = grades
	}

	for subjectId, subjectGrade := range subjectsGrade {
		fmt.Printf("________________\n")
		fmt.Printf("%-13s | Mean\n", objectsData[subjectId].Name)
		fmt.Printf("________________\n")

		var totalSum int
		var totalCount int

		for grade, grades := range subjectGrade {
			sum := 0
			for _, gradeValue := range grades {
				sum += gradeValue
			}
			mean := float64(sum) / float64(len(grades))
			totalSum += sum
			totalCount += len(grades)

			fmt.Printf("%-7d grade | %.1f\n", grade, mean)
		}

		totalMean := float64(totalSum) / float64(totalCount)
		fmt.Printf("________________\n")
		fmt.Printf("mean          | %.1f\n", totalMean)
		fmt.Printf("________________\n")
	}
}
