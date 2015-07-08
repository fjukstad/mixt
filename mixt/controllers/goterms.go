package controllers

import (
	"fmt"

	"bitbucket.org/vdumeaux/mixt/mixt/mixt"
)

type GOTerm struct {
	Name      string
	Id        string
	TissueSet map[string][]mixt.GOTerm
}

func GOTermResults(searchTerms []string) ([]GOTerm, error) {
	tissues, err := mixt.GetTissues()
	if err != nil {
		fmt.Println("getting tissues went bad:", err)
		return []GOTerm{}, err
	}

	var GOTermNames []string
	for _, term := range searchTerms {
		hits, err := SearchForGOTerms(term)
		if err != nil {
			return []GOTerm{}, err
		}
		GOTermNames = append(GOTermNames, hits...)
	}

	var goterms []GOTerm

	for _, GOTermName := range GOTermNames {
		set := make(map[string][]mixt.GOTerm)
		var id string
		for _, tissue := range tissues {
			if tissue == "nblood" {
				continue
			}
			score, err := mixt.GetGOScoresForTissue(tissue, GOTermName)
			if err != nil {
				return []GOTerm{}, err
			}
			set[tissue] = score
			//GOTermScore{score.Module, score.ClassicFisher,
			//	score.Weight01Fisher}
			id = score[0].GOId
		}
		goterms = append(goterms, GOTerm{GOTermName, id, set})
	}

	return goterms, nil
}
