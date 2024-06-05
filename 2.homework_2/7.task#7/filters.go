package main

func Filter(database Database) (map[int]Students, map[int]Objects, []Results) {
	const excellentGrade = 5
	studentsData := make(map[int]Students)
	objectsData := make(map[int]Objects)
	excellentResults := []Results{}

	for _, student := range database.Students {
		studentsData[student.Id] = student
	}

	for _, object := range database.Objects {
		objectsData[object.Id] = object
	}

	for _, res := range database.Results {
		if res.Result == excellentGrade {
			excellentResults = append(excellentResults, res)
		}
	}

	return studentsData, objectsData, excellentResults
}
