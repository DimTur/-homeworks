package main

func main() {
	dbUrl := "../dz3.json"
	dataBase := getDb(dbUrl)

	studentsCache := writeStudentsCache(dbUrl)
	objectsCache := writeObjectsCache(dbUrl)

	joinTable(dataBase, studentsCache, objectsCache)
}
