package main

type Cache[K comparable, V any] struct {
	m map[K]V
}

func (c *Cache[K, V]) Init() {
	c.m = make(map[K]V)
}

func (c *Cache[K, V]) Set(key K, value V) {
	c.m[key] = value
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	k, ok := c.m[key]
	return k, ok
}

func writeStudentsCache(dbUrl string) Cache[int, Students] {
	dataBase := getDb(dbUrl)

	studentsCache := Cache[int, Students]{}
	studentsCache.Init()

	for _, student := range dataBase.Students {
		studentsCache.Set(student.Id, student)
	}

	return studentsCache
}

func writeObjectsCache(dbUrl string) Cache[int, Objects] {
	dataBase := getDb(dbUrl)

	objectsCache := Cache[int, Objects]{}
	objectsCache.Init()

	for _, object := range dataBase.Objects {
		objectsCache.Set(object.Id, object)
	}
	return objectsCache
}
