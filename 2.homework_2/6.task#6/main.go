package main

func main() {
	dbUrl := "../dz3.json"
	dataBase := getDb(dbUrl)

	subject1 := "Math"
	subject2 := "Biology"
	subject3 := "Geography"

	summaryTable(dataBase, subject1)
	summaryTable(dataBase, subject2)
	summaryTable(dataBase, subject3)
}
