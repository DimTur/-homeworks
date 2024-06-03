package main

import (
	"fmt"
	"sort"
)

type Candidate struct {
	Name  string
	Votes int
}

func countVotes(s []string) map[string]int {
	result := make(map[string]int)
	for _, i := range s {
		result[i] += 1
	}
	return result
}

func sortVotes(votes map[string]int) []Candidate {
	candidates := make([]Candidate, 0, len(votes))
	for name, voteCount := range votes {
		candidates = append(candidates, Candidate{Name: name, Votes: voteCount})
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Votes > candidates[j].Votes
	})

	return candidates
}

func main() {
	request := []string{"Ann", "Kate", "Peter", "Kate", "Ann", "Ann", "Helen"}
	fmt.Println(request)
	voteCounter := countVotes(request)
	fmt.Println(sortVotes(voteCounter))
}
