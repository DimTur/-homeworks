package main

import "fmt"

func joinTable(database Database) {
	studentsData := make(map[int]Students)
	objectsData := make(map[int]Objects)

	for _, student := range database.Students {
		studentsData[student.Id] = student
	}

	for _, object := range database.Objects {
		objectsData[object.Id] = object
	}

	fmt.Println("____________________________________________")
	fmt.Println("Student name  | Grade | Object    |   Result")
	fmt.Println("____________________________________________")

	for _, result := range database.Results {
		student := studentsData[result.Student_id]
		object := objectsData[result.Object_id]
		fmt.Printf("%-15s | %-5d | %-10s | %-7d\n", student.Name, student.Grade, object.Name, result.Result)

	}
	fmt.Println("____________________________________________")
}
