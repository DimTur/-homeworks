package main

import (
	"fmt"
	"log"
)

func joinTable(database Database, studetsCache Cache[int, Students], objectsCache Cache[int, Objects]) {
	fmt.Println("____________________________________________")
	fmt.Println("Student name    | Grade | Object     | Result")
	fmt.Println("____________________________________________")

	for _, result := range database.Results {
		student, ok := studetsCache.Get(result.Student_id)
		if !ok {
			log.Fatalf("Student with ID %d not found", result.Student_id)
		}
		object, ok := objectsCache.Get(result.Object_id)
		if !ok {
			log.Fatalf("Student with ID %d not found", result.Object_id)
		}

		fmt.Printf("%-15s | %-5d | %-10s | %-7d\n", student.Name, student.Grade, object.Name, result.Result)

	}
	fmt.Println("____________________________________________")
}
