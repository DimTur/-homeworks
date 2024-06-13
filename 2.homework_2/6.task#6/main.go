package main

func main() {
	dbUrl := "../dz3.json"
	dataBase := getDb(dbUrl)

	var subjects []string
	for _, object := range dataBase.Objects {
		subjects = append(subjects, object.Name)
	}

	for _, sub := range subjects {
		summaryTable(dataBase, sub, func(subject string, filteredGrades map[int][]int) {
			printSummary(subject, filteredGrades)
		})
	}
}
