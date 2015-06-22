package mixt

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/fjukstad/kvik/kompute"
)

var komp *kompute.Kompute

func Init(addr, username, password string) {
	komp = kompute.NewKompute(addr, username, password)
	return
}

func Heatmap(tissue, module string) (string, error) {
	session, err := komp.Call("mixt/R/heatmap", `{"tissue": `+"\""+tissue+"\""+`,
					"module":`+"\""+module+"\""+`}`)

	if err != nil {
		fmt.Println("Heatmapper")
		fmt.Println(session, err)
		return "", err
	}

	plotUrl := session.Key

	return plotUrl, nil
}

func GetGenes() ([]string, error) {
	resp, err := komp.Rpc("mixt/R/getAllGenes", "", "csv")
	if err != nil {
		fmt.Println("error:", err)
		return []string{""}, err
	}

	body := strings.NewReader(resp)
	reader := csv.NewReader(body)
	var genes []string
	line := 0
	for {
		record, err := reader.Read()

		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return []string{}, nil
		}

		if line == 0 {
			line += 1
			continue
		}

		name := record[0]
		genes = append(genes, name)
	}

	return genes, nil
}

func GetCommonGenes(tissue, module, geneset, status string) ([]string, error) {
	if status == "" {
		status = "updn.common"
	}

	args := `{"tissue": ` + "\"" + tissue + "\"" + `,
					"module":` + "\"" + module + "\"" + `,
					"geneset":` + "\"" + geneset + "\"" + `,
					"status": ` + "\"" + status + "\"" + `}`

	resp, err := komp.Rpc("mixt/R/getCommonGenes", args, "json")
	if err != nil {
		fmt.Println("Could not get common genes", err)
		return []string{}, err
	}

	geneNames := make([]string, 0)
	err = json.Unmarshal([]byte(resp), &geneNames)
	if err != nil {
		fmt.Println("Could not unmarshal gene names", err)
		return []string{}, err
	}
	return geneNames, nil

}

func GetAllModuleNames(gene string) ([]string, error) {
	resp, err := komp.Rpc("mixt/R/getAllModules", `{"gene": `+"\""+gene+"\""+`}`, "json")
	if err != nil {
		fmt.Println("error:", err)
		return []string{""}, err
	}

	moduleNames := make([]string, 0)
	err = json.Unmarshal([]byte(resp), &moduleNames)
	if err != nil {
		fmt.Println("cannot unmarshal json response", err)
		return nil, err
	}

	return moduleNames, nil
}

func GetTissues() ([]string, error) {
	resp, err := komp.Rpc("mixt/R/getAllTissues", "", "json")
	if err != nil {
		fmt.Println("error:", err)
		return []string{""}, err
	}

	tissues := make([]string, 0)
	err = json.Unmarshal([]byte(resp), &tissues)
	if err != nil {
		fmt.Println("cannot unmarshal json response", err)
		return nil, err
	}

	return tissues, nil
}

type Module struct {
	Name             string
	Tissue           string
	HeatmapUrl       string
	Genes            []Gene
	EnrichmentScores EnrichmentScores
	GOTerms          []GOTerm
	Url              string
}

type Gene struct {
	Name        string
	Correlation string
	K           float64
	Kin         float64
	Updown      string
}

type Response struct {
	Item string
}

func GetModules(tissue string) ([]Module, error) {
	resp, err := komp.Rpc("mixt/R/getModules", `{"tissue" : `+"\""+tissue+"\""+`}`, "json")
	if err != nil {
		fmt.Println("error:", err)
		return nil, err
	}

	moduleNames := make([]string, 0)
	err = json.Unmarshal([]byte(resp), &moduleNames)
	if err != nil {
		fmt.Println("cannot unmarshal json response", err)
		fmt.Println(resp)
		return nil, err

	}

	// removing the grey module
	moduleNames = moduleNames[1:len(moduleNames)]

	var modules []Module

	resChan := make(chan Module)

	for i, _ := range moduleNames {
		go func(i int) {
			/*
				m, err := GetModule(moduleNames[i], tissue)

				if err != nil {
					fmt.Println("cannot get module", moduleNames[i], err)
					resChan <- Module{}
					return
				}
			*/
			m := Module{moduleNames[i], tissue, "", nil, EnrichmentScores{}, []GOTerm{}, ""}
			resChan <- m
		}(i)
	}

	for range moduleNames {
		m := <-resChan
		m.Tissue = tissue
		modules = append(modules, m)
	}

	sort.Sort(ByName(modules))

	return modules, nil
}

type ByName []Module

func (a ByName) Len() int {
	return len(a)
}
func (a ByName) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByName) Less(i, j int) bool {
	return a[i].Name < a[j].Name
}

func GetModule(name string, tissue string) (Module, error) {

	if name == "grey" {
		return Module{}, nil
	}

	heatmapUrl, err := Heatmap(tissue, name)
	if err != nil {
		fmt.Println("heatmap")
		return Module{}, err
	}

	genes, url, err := GetGeneList(name, tissue)
	if err != nil {
		fmt.Println("ghenelist")
		return Module{}, err
	}

	scores, err := GetEnrichmentScores(name, tissue)
	if err != nil {
		fmt.Println("Could not get enrichment scores")
		return Module{}, err
	}

	goterms, err := GetGOTerms(name, tissue, "")
	if err != nil {
		fmt.Println("Could not get goterms", err)
		return Module{}, err
	}

	module := Module{name, tissue, heatmapUrl, genes, scores, goterms, url}
	return module, nil
}

func GetGeneList(module, tissue string) (genes []Gene, url string,
	err error) {
	var args string
	args = `{"tissue": ` + "\"" + tissue + "\"" + `,
					"module":` + "\"" + module + "\"" + `}`

	session, err := komp.Call("mixt/R/getGeneList", args)

	if err != nil {
		fmt.Println("Call failed one time, trying again. ")
		fmt.Println(session, err)

		session, err = komp.Call("mixt/R/getGeneList", args)
		if err != nil {
			fmt.Println("Call failed for the second time :( ")
			return nil, "", err
		}
	}

	resp, err := session.GetResult(komp, "csv")
	if err != nil {
		fmt.Println("GETRESULT", session.Url, session.Result)
		fmt.Println(resp, err)

		return nil, "", err
	}

	body := strings.NewReader(resp)
	reader := csv.NewReader(body)

	line := 0
	for {
		record, err := reader.Read()

		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return []Gene{}, "", nil
		}

		if line == 0 {
			line += 1
			continue
		}

		name := record[0]
		var updown string
		if record[1] == "-1" {
			updown = "down"
		} else {
			updown = "up"
		}

		cor, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			fmt.Println("Could not parse float :( ", err)
			return []Gene{}, "", err
		}

		c := fmt.Sprintf("%.4g", cor)

		g := Gene{Name: name,
			Correlation: c,
			K:           0,
			Kin:         0,
			Updown:      updown}

		genes = append(genes, g)
	}

	fmt.Println("WE GOT", len(genes), "genes")

	return genes, session.Key, nil
}

type Score struct {
	Set          string  `json:"sig.set"`
	Name         string  `json:"_row"`
	Size         int     `json:"sig.size,string"`
	UpDownCommon int     `json:"updn.common,string"`
	UpDownPvalue float64 `json:"updn.pval,string"`
	UpCommon     int     `json:"up.common,string"`
	UpPvalue     float64 `json:"up.pval,string"`
	DownCommon   int     `json:"dn.common,string"`
	DownPvalue   float64 `json:"dn.pval,string"`
	Tissue       string  `json:"tissue,omitempty"`
}

type EnrichmentScores struct {
	Sets map[string][]Set
}

type Set struct {
	SetName       string
	SignatureName string
	Size          int
	UpDownCommon  int
	UpDownPvalue  float64
	UpCommon      int
	UpPvalue      float64
	DownCommon    int
	DownPvalue    float64
}

func GetEnrichmentScores(module, tissue string) (enrichment EnrichmentScores, err error) {

	args := `{"tissue": ` + "\"" + tissue + "\"" + `,
					"module":` + "\"" + module + "\"" + `}`

	resp, err := komp.Rpc("mixt/R/getEnrichmentScores", args, "json")

	if err != nil {
		fmt.Println("Could not get enrich")
		return EnrichmentScores{}, err
	}

	res := []byte(resp)

	var scores []Score
	err = json.Unmarshal(res, &scores)
	if err != nil {
		fmt.Println("Unm", err)
		fmt.Println(res)
		return EnrichmentScores{}, err
	}

	sets := make(map[string][]Set)

	for _, s := range scores {
		name := s.Set
		name = strings.Replace(name, ".", "", -1)
		sets[name] = append(sets[name], Set{name, s.Name, s.Size, s.UpDownCommon, s.UpDownPvalue,
			s.UpCommon, s.UpPvalue, s.DownCommon, s.DownPvalue})

	}

	enrichment = EnrichmentScores{sets}
	return enrichment, nil

}

func GetEnrichmentScore(module, tissue, geneset string) (Score, error) {

	args := `{"tissue": ` + "\"" + tissue + "\"" + `,
	"module":` + "\"" + module + "\"" + `, 
	"geneset":` + "\"" + geneset + "\"" + `}`

	resp, err := komp.Rpc("mixt/R/getEnrichmentScores", args, "json")

	if err != nil {
		return Score{}, err
	}

	res := []byte(resp)

	var score []Score
	err = json.Unmarshal(res, &score)
	if err != nil {
		fmt.Println(err)
		return Score{}, err
	}
	return score[0], nil

}

func GetGeneSetNames() ([]string, error) {

	resp, err := komp.Rpc("mixt/R/getGeneSetNames", "", "json")
	if err != nil {
		fmt.Println("Could not get gene set names :(", err)
		return []string{}, err
	}

	res := []byte(resp)

	var names []string
	err = json.Unmarshal(res, &names)
	return names, err

}

type ModuleScores struct {
	_ []map[string]Score
}

func GetEnrichmentForTissue(tissue, geneset string) ([]Score, error) {

	args := `{"tissue": ` + "\"" + tissue + "\"" + `,
	"geneset":` + "\"" + geneset + "\"" + `}`

	resp, err := komp.Rpc("mixt/R/getEnrichmentForTissue", args, "json")

	if err != nil {
		return []Score{}, err
	}

	res := []byte(resp)

	var modulescores []Score
	err = json.Unmarshal(res, &modulescores)
	return modulescores, err

}

type GOTerm struct {
	GOId          string `json:"GO.ID"`
	Term          string
	Annotated     int
	Significant   int
	Expected      float64
	ClassicFisher string `json:"classicFisher"`
}

func GetGOTerms(module, tissue, terms string) ([]GOTerm, error) {

	var args string
	if len(terms) < 1 {
		args = `{"tissue": ` + "\"" + tissue + "\"" + `,
		"module":` + "\"" + module + "\"" + `}`
	} else {

		args = `{"tissue": ` + "\"" + tissue + "\"" + `,
		"module":` + "\"" + module + "\"" + `, 
		"terms":` + "\"" + terms + "\"" + `}`
	}

	resp, err := komp.Rpc("mixt/R/getGOTerms", args, "json")

	if err != nil {
		return []GOTerm{}, err
	}

	res := []byte(resp)

	var goterms []GOTerm
	err = json.Unmarshal(res, &goterms)
	return goterms, err

}

func Get(key, filetype string) ([]byte, error) {
	return komp.Get(key, filetype)
}
