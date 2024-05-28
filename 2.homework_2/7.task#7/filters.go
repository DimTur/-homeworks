package main

func Filter(r int, results []Results) []Results {
	filtered := []Results{}

	for _, res := range results {
		if res.Result == r {
			filtered = append(filtered, res)
		}
	}

	return filtered
}
