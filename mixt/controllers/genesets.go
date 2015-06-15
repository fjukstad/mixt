package controllers

import (
	"fmt"

	"bitbucket.org/vdumeaux/mixt/mixt/mixt"
)

type GeneSet struct {
	Name string
	// tissue -> enrichment scores for each module
	TissueSets map[string][]ModuleEnrichment
}

type ModuleEnrichment struct {
	Name   string
	Pvalue float64
}

func SetResults(searchTerms []string) ([]GeneSet, error) {
	fmt.Println(searchTerms)

	// get all tissues
	tissues, err := mixt.GetTissues()
	if err != nil {
		fmt.Println("getting tissues went bad:", err)
		return []GeneSet{}, err
	}

	// fetch all gene sets that matched the search
	var geneSetNames []string
	for _, term := range searchTerms {
		hits, err := SearchForGeneSet(term)
		if err != nil {
			return []GeneSet{}, err
		}
		geneSetNames = append(geneSetNames, hits...)
	}

	var gs []GeneSet

	// fetch enrichment scores for the genesets in the
	for _, geneSetName := range geneSetNames {
		ts := make(map[string][]ModuleEnrichment)
		for _, tissue := range tissues {
			var mes []ModuleEnrichment
			sc, err := mixt.GetEnrichmentForTissue(tissue, geneSetName)
			if err != nil {
				fmt.Println("Could not get enrichment for tissue", err)
				return []GeneSet{}, err
			}

			for _, s := range sc {
				mes = append(mes, ModuleEnrichment{s.Name, s.UpDownPvalue})
				fmt.Println(s.Name, s.UpDownPvalue)
			}
			fmt.Println(tissue)
			ts[tissue] = mes
		}
		gs = append(gs, GeneSet{geneSetName, ts})
	}

	return gs, nil
}
