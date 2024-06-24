package main

import "fmt"

func joinTable(database Database, filter func(Database) (map[int]Students, map[int]Objects, []Results)) {
	studentsData, objectsData, filteredResults := filter(database)

	fmt.Println("____________________________________________")
	fmt.Println("Student name    | Grade | Object     |   Result")
	fmt.Println("____________________________________________")

	for _, result := range filteredResults {
		student := studentsData[result.Student_id]
		object := objectsData[result.Object_id]
		fmt.Printf("%-15s | %-5d | %-10s | %-7d\n", student.Name, student.Grade, object.Name, result.Result)
	}
	fmt.Println("____________________________________________")
}
