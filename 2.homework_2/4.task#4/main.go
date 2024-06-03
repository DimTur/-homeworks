package main

func main() {
	dbUrl := "../dz3.json"
	dataBase := getDb(dbUrl)
	summaryTable(dataBase)
}
