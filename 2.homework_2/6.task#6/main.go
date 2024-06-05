package main

func main() {
	dbUrl := "../dz3.json"
	dataBase := getDb(dbUrl)

	subjects := []string{"Math", "Biology", "Geography"}

	for _, sub := range subjects {
		summaryTable(dataBase, sub, func(subject string, filteredGrades map[int][]int) {
			printSummary(subject, filteredGrades)
		})
	}
}
