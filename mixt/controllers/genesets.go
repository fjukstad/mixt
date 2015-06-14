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

	// get all modules
	modules := make(map[string][]mixt.Module)
	for _, tissue := range tissues {
		m, err := mixt.GetModules(tissue)
		if err != nil {
			fmt.Println("Could not get modules")
			return []GeneSet{}, err
		}
		modules[tissue] = m
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

	ts := make(map[string][]ModuleEnrichment)
	var gs []GeneSet

	// fetch enrichment scores for the genesets in the
	for _, geneSetName := range geneSetNames {
		for tissue, ms := range modules {
			var mes []ModuleEnrichment
			for _, module := range ms {
				sc, err := mixt.GetEnrichmentScore(module.Name, tissue, geneSetName)
				if err != nil {
					fmt.Println("Could not get ER score")
					return []GeneSet{}, err
				}

				mes = append(mes, ModuleEnrichment{module.Name, sc.UpDownPvalue})

			}
			ts[tissue] = mes
		}
		gs = append(gs, GeneSet{geneSetName, ts})
	}

	return gs, nil
}
