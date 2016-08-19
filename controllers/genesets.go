package controllers

import (
	"fmt"
	"net/http"

	"github.com/fjukstad/kvik/gsea"
	"github.com/gorilla/mux"

	"bitbucket.org/vdumeaux/mixt/mixt"
)

type GeneSet struct {
	Name string
	// tissue -> enrichment scores for each module
	TissueSets map[string][]ModuleEnrichment
}

type ModuleEnrichment struct {
	Name     string
	Pvalue   float64
	Abstract string
}

func SetResults(searchTerms []string) ([]GeneSet, error) {
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
	setResults := make(chan GeneSet, 10)

	// fetch enrichment scores for the genesets in the
	for _, geneSetName := range geneSetNames {
		go func(geneSetName string) {
			ts := make(map[string][]ModuleEnrichment)
			for _, tissue := range tissues {
				var mes []ModuleEnrichment
				sc, err := mixt.GetEnrichmentForTissue(tissue, geneSetName)
				if err != nil {
					fmt.Println("Could not get enrichment for tissue", err)
					setResults <- GeneSet{}
					return
				}

				for _, s := range sc {
					abstract, _ := gsea.Abstract(s.Name)
					mes = append(mes, ModuleEnrichment{s.Name, s.UpDownPvalue, abstract})
				}
				ts[tissue] = mes
			}
			setResults <- GeneSet{geneSetName, ts}
			//gs = append(gs, GeneSet{geneSetName, ts})
		}(geneSetName)
	}

	for range geneSetNames {
		geneset := <-setResults
		if geneset.Name == "" {
			fmt.Println("Could not get enrichment for tissue", err)
			return []GeneSet{}, err
		}
		gs = append(gs, geneset)
	}

	return gs, nil
}

func GeneSetAbstractHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	geneset := vars["geneset"]
	abstract, _ := gsea.Abstract(geneset)
	if len(abstract) <= 2 {
		briefDescription, _ := gsea.BriefDescription(geneset)
		w.Write([]byte(briefDescription))
		return
	}
	w.Write([]byte(abstract))
}
