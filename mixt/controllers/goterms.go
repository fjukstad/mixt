package controllers

import (
	"fmt"

	"bitbucket.org/vdumeaux/mixt/mixt/mixt"
)

type GOTerm struct {
	Name      string
	TissueSet map[string][]GOTermScore
}

type GOTermScore struct {
	Name   string
	PValue float64
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

	fmt.Println("GO terms", GOTermNames)
	//var goterms []GOTerm

	for _, GOTermName := range GOTermNames {
		fmt.Println(GOTermName)
		for _, tissue := range tissues {
			if tissue == "nblood" {
				continue
			}
			mixt.GetGOScoresForTissue(tissue, GOTermName)
		}
	}

	return []GOTerm{}, nil
}
