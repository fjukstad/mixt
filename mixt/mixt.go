package mixt

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/fjukstad/kvik/r"
	"github.com/pkg/errors"
)

var R r.Client
var pkg = "mixtbt"

func Init(addr, username, password string) {
	R = r.Client{addr, username, password}
	return
}

func Heatmap(tissue, module string) (string, error) {
	fun := "heatmap"
	args := "tissue='" + tissue + "', module='" + module + "'"
	return plot(pkg, fun, args)

}

func HeatmapReOrder(tissue, module, orderByTissue, orderByModule, cohort string) (string, error) {
	args := "tissue='" + tissue + "', module='" + module + "', orderByModule='" + orderByModule + "', orderByTissue='" + orderByTissue + "', cohort.name='" + cohort + "'"
	fun := "cohort_heatmap"
	return plot(pkg, fun, args)
}

func CohortHeatmap(tissue, module, cohort string) (string, error) {
	fun := "cohort_heatmap"
	args := "tissue='" + tissue + "', module='" + module + "', cohort.name='" + cohort + "'"
	return plot(pkg, fun, args)
}

func CohortBoxplot(module, orderByTissue, orderByModule, cohort string) (string, error) {
	fun := "cohort_boxplot"
	args := "blood.module='" + module + "', orderByTissue='" + orderByTissue + "', orderByModule='" + orderByModule + "', cohort='" + cohort + "'"
	return plot(pkg, fun, args)
}

func CohortScatterplot(tissueA, tissueB, moduleA, moduleB, cohort string) (string, error) {
	fun := "cohort_scatterplot"
	args := "x.tissue='" + tissueA + "', y.tissue='" + tissueB + "', x.module='" + moduleA + "', y.module='" + moduleB + "', cohort.name='" + cohort + "'"
	return plot(pkg, fun, args)
}

func plot(pkg, fun, args string) (string, error) {

	key, err := R.Call(pkg, fun, args)

	if err != nil {
		fmt.Println("Could not plot :( ")
		fmt.Println(key, err)
		return "", err
	}

	return key, nil
}

func GetGenes() ([]string, error) {

	key, err := R.Call(pkg, "getAllGenes", "")
	if err != nil {
		return []string{}, err
	}

	resp, err := R.Get(key, "csv")
	if err != nil {
		fmt.Println("error:", err)
		return []string{""}, err
	}

	body := bytes.NewReader(resp)
	reader := csv.NewReader(body)
	var genes []string
	line := 0
	for {
		record, err := reader.Read()

		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return []string{}, err
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

	fun := "getCommonGenes"
	args := "tissue='" + tissue + "', module='" + module + "', geneset='" + geneset + "', status='" + status + "'"

	key, err := R.Call(pkg, fun, args)
	if err != nil {
		return []string{}, err
	}

	resp, err := R.Get(key, "json")
	if err != nil {
		fmt.Println("Could not get common genes", err)
		return []string{}, err
	}

	geneNames := make([]string, 0)
	err = json.Unmarshal([]byte(resp), &geneNames)
	if err != nil {
		return []string{}, errors.Wrap(err, "get common genes unmarshal")
	}
	return geneNames, nil

}

func GetAllModuleNames(gene string) ([]string, error) {
	fun := "getAllModules"
	args := "gene='" + gene + "'"

	key, err := R.Call(pkg, fun, args)
	if err != nil {
		return []string{}, err
	}

	resp, err := R.Get(key, "json")
	if err != nil {
		fmt.Println("error:", err)
		return []string{""}, err
	}

	moduleNames := make([]string, 0)
	err = json.Unmarshal([]byte(resp), &moduleNames)
	if err != nil {
		return nil, errors.Wrap(err, "getallmodulenames unmarshal")
	}

	return moduleNames, nil
}

func GetTissues() ([]string, error) {

	fun := "getAllTissues"
	args := ""

	key, err := R.Call(pkg, fun, args)
	if err != nil {
		return []string{}, err
	}

	resp, err := R.Get(key, "json")
	if err != nil {
		return []string{""}, errors.Wrap(err, "get tissues r.get failed")
	}

	tissues := make([]string, 0)
	err = json.Unmarshal([]byte(resp), &tissues)
	if err != nil {
		return nil, errors.Wrap(err, "gettissues unmarshal")
	}

	return tissues, nil
}

type Module struct {
	Name                  string
	Tissue                string
	HeatmapUrl            string
	AlternativeHeatmapUrl string
	Genes                 []Gene
	EnrichmentScores      EnrichmentScores
	GOTerms               []GOTerm
	Url                   string
	ScatterplotUrl        string
	BoxplotUrl            string
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

	fun := "getModules"
	args := "tissue='" + tissue + "'"

	key, err := R.Call(pkg, fun, args)
	if err != nil {
		return nil, err
	}

	resp, err := R.Get(key, "json")
	if err != nil {
		fmt.Println("error:", err)
		return nil, err
	}

	moduleNames := make([]string, 0)
	err = json.Unmarshal([]byte(resp), &moduleNames)
	if err != nil {
		return nil, errors.Wrap(err, "get modules unmarshal error"+string(resp))
	}

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
			m := Module{moduleNames[i], tissue, "", "", nil, EnrichmentScores{}, []GOTerm{}, "", "", ""}
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

func GetModule(name, tissue, cohort string) (Module, error) {

	if name == "grey" {
		return Module{}, nil
	}

	heatmapUrl, err := CohortHeatmap(tissue, name, cohort)
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

	goterms, err := GetGOTerms(name, tissue, []string{})
	if err != nil {
		fmt.Println("Could not get goterms", err)
		return Module{}, err
	}

	alternativeHeatmapUrl := ""
	if tissue == "blood" {
		alternativeHeatmapUrl, err = CohortHeatmap("bnblood", name, cohort)
		if err != nil {
			fmt.Println("Couldt not get bnblood heatmap", err)
			return Module{}, err
		}
	}

	cohortBoxplot, err := CohortBoxplot(name, tissue, name, cohort)
	if err != nil {
		fmt.Println("Could not generate boxplot", err)
		//return Module{}, err
	}

	module := Module{name, tissue, heatmapUrl, alternativeHeatmapUrl, genes, scores, goterms, url, "", cohortBoxplot}
	return module, nil
}

func GeneListCSV(modules []string, tissue string) ([]byte, error) {
	header := "module, genes \n"
	lines := header
	sep := " "
	for _, module := range modules {
		genes, _, err := GetGeneList(module, tissue)
		if err != nil {
			return nil, err
		}
		line := module + ", "
		for _, gene := range genes {
			line = line + gene.Name + sep
		}
		lines = lines + line + "\n"
	}

	fmt.Println(lines)
	return []byte(lines), nil
}

func GetGeneList(module, tissue string) (genes []Gene, url string,
	err error) {
	fun := "getGeneList"
	args := "tissue='" + tissue + "', module='" + module + "'"

	key, err := R.Call(pkg, fun, args)
	if err != nil {
		return nil, "", err
	}

	resp, err := R.Get(key, "csv")
	if err != nil {
		return nil, "", err
	}

	body := bytes.NewReader(resp)
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

	return genes, key, nil

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

	fun := "getEnrichmentScores"
	args := "tissue='" + tissue + "', module='" + module + "'"

	key, err := R.Call(pkg, fun, args)
	if err != nil {
		return EnrichmentScores{}, err
	}

	resp, err := R.Get(key, "json")
	if err != nil {
		fmt.Println("Could not get enrich")
		return EnrichmentScores{}, err
	}

	res := []byte(resp)

	var scores []Score
	err = json.Unmarshal(res, &scores)
	if err != nil {
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

	fun := "getEnrichmentScores"
	args := "tissue='" + tissue + "', module='" + module + "', geneset='" + geneset + "'"

	key, err := R.Call(pkg, fun, args)
	if err != nil {
		return Score{}, err
	}

	resp, err := R.Get(key, "json")
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
	return GetSlice(pkg, "getGeneSetNames", "")
}

func GetGOTermNames() ([]string, error) {
	return GetSlice(pkg, "getGOTermNames", "")
}

func GetSlice(pkg, fun, args string) ([]string, error) {
	key, err := R.Call(pkg, fun, args)
	if err != nil {
		return []string{}, err
	}

	resp, err := R.Get(key, "json")
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

	fun := "getEnrichmentForTissue"
	args := "tissue='" + tissue + "', geneset='" + geneset + "'"

	key, err := R.Call(pkg, fun, args)
	if err != nil {
		return []Score{}, err
	}

	resp, err := R.Get(key, "json")
	if err != nil {
		return []Score{}, err
	}

	res := []byte(resp)

	var modulescores []Score
	err = json.Unmarshal(res, &modulescores)
	if err != nil {
		fmt.Println("Enrichment for tissue:.", err)
		return []Score{}, err
	}
	return modulescores, err

}

type GOTerm struct {
	GOId           string `json:"GO.ID"`
	Term           string
	Annotated      int
	Significant    int
	Expected       float64
	ClassicFisher  string `json:"classicFisher"`
	Weight01Fisher string `json:"weight01Fisher"`
	Rank           int    `json:"Rank in weight01Fisher"`
	Module         string `json:"module"`
}

func GetGOTerms(module, tissue string, terms []string) ([]GOTerm, error) {
	fun := "getGOTerms"
	args := "tissue='" + tissue + "', module='" + module + "'"
	if len(terms) > 1 {

		var fmtterms []string
		for i, _ := range terms {
			fmtterms = append(fmtterms, "\""+terms[i]+"\"")
		}

		goTermNames := "["
		goTermNames += strings.Join(fmtterms, ", ")
		goTermNames += "]"

		args = args + ", terms='" + goTermNames + "'"
	}

	key, err := R.Call(pkg, fun, args)
	if err != nil {
		return []GOTerm{}, err
	}

	resp, err := R.Get(key, "json")
	if err != nil {
		return []GOTerm{}, err
	}

	res := []byte(resp)

	var goterms []GOTerm
	err = json.Unmarshal(res, &goterms)
	return goterms, err

}

func GetGOScoresForTissue(tissue, goterm string) ([]GOTerm, error) {

	fun := "getGOScoresForTissue"
	args := "tissue='" + tissue + "', term='" + goterm + "'"

	key, err := R.Call(pkg, fun, args)
	if err != nil {
		return []GOTerm{}, err
	}

	resp, err := R.Get(key, "json")
	if err != nil {
		return []GOTerm{}, err
	}

	res := []byte(resp)

	var scores []GOTerm
	err = json.Unmarshal(res, &scores)
	return scores, err

}

type UserScore struct {
	PValue float64  `json:"p_values"`
	Module string   `json:"module"`
	Common []string `json:"common"`
}

func UserEnrichmentScores(tissue string, genelist []string) ([]UserScore, error) {

	genes := parseGeneList(genelist)

	fun := "userEnrichmentScores"
	args := "tissue='" + tissue + "', genelist=" + genes

	key, err := R.Call(pkg, fun, args)
	if err != nil {
		return []UserScore{}, err
	}

	resp, err := R.Get(key, "json")
	if err != nil {
		fmt.Println("Could not calculate er scores for user list")
		return []UserScore{}, err
	}

	res := []byte(resp)

	var scores []UserScore
	err = json.Unmarshal(res, &scores)
	return scores, err
}

func parseGeneList(genelist []string) string {

	var fmtgenelist []string
	for i, _ := range genelist {
		fmtgenelist = append(fmtgenelist, "\""+genelist[i]+"\"")
	}

	genes := "c("
	genes += strings.Join(fmtgenelist, ", ")
	genes += ")"

	return genes
}

func GetCommonGOTermGenes(module, tissue, id string) ([]string, error) {
	fun := "getCommonGOTermGenes"
	args := "tissue='" + tissue + "', module='" + module + "', gotermID='" + id + "'"

	return GetSlice(pkg, fun, args)

}

func GetCommonUserERGenes(module, tissue string, genelist []string) ([]string, error) {

	genes := parseGeneList(genelist)

	fun := "commonEnrichmentScoreGenes"
	args := "tissue='" + tissue + "', module='" + module + "', genelist='" + genes + "'"

	return GetSlice(pkg, fun, args)
}

func EigengeneCorrelation(tissueA, tissueB string) ([]byte, error) {
	args := "tissueA='" + tissueA + "', tissueB='" + tissueB + "'"
	return analysis("eigengeneCorrelation", args)
}

func ModuleHypergeometricTest(tissueA, tissueB string) ([]byte, error) {
	args := "tissueA='" + tissueA + "', tissueB='" + tissueB + "'"
	return analysis("moduleHypergeometricTest", args)
}

func ROITest(tissueA, tissueB string) ([]byte, error) {
	args := "tissueA='" + tissueA + "', tissueB='" + tissueB + "'"
	return analysis("roiTest", args)
}

func PatientRankCorrelation(tissueA, tissueB string) ([]byte, error) {
	args := "tissueA='" + tissueA + "', tissueB='" + tissueB + "'"
	return analysis("patientRankCorrelation", args)
}

func ClinicalEigengene(tissue string) ([]byte, error) {
	args := "tissue='" + tissue + "'"
	return analysis("eigengeneClinicalRelation", args)
}

func ClinicalROI(tissue string) ([]byte, error) {
	args := "tissue='" + tissue + "'"
	return analysis("roiClinicalRelation", args)
}

func ClinicalRanksum(tissue, cohort string) ([]byte, error) {
	args := "tissue='" + tissue + "'" + ", cohort='" + cohort + "'"
	return analysis("clinicalRanksum", args)
}

func PatientRankSum(tissueA, tissueB, cohort string) ([]byte, error) {
	args := "tissueA='" + tissueA + "', tissueB='" + tissueB + "',cohort='" + cohort + "'"
	return analysis("patientRankSum", args)
}

func GeneOverlapTest(tissueA, tissueB string) ([]byte, error) {
	args := "tissueA='" + tissueA + "', tissueB='" + tissueB + "'"
	return analysis("geneOverlapTest", args)
}

func analysis(fun, args string) ([]byte, error) {
	key, err := R.Call(pkg, fun, args)
	if err != nil {
		fmt.Println("Could not run analysis:", err)
		return nil, err
	}

	return R.Get(key, "csv")
}

type Analyses struct {
	Ranksum []float64 `json:"ranksum"`
	Overlap []float64 `json:"overlap"`
	Common  []string  `json:"common"`
}

func ModuleComparisonAnalyses(tissueA, tissueB, moduleA, moduleB string) (Analyses, error) {

	fun := "comparisonAnalyses"
	args := "tissueA='" + tissueA + "', tissueB='" + tissueB + "', moduleA='" + moduleA + "', moduleB='" + moduleB + "'"
	key, err := R.Call(pkg, fun, args)
	if err != nil {
		return Analyses{}, err
	}

	resp, err := R.Get(key, "json")
	if err != nil {
		return Analyses{}, err
	}
	res := []byte(resp)

	var analyses Analyses
	err = json.Unmarshal(res, &analyses)
	return analyses, err
}

func GetTOMGraph(tissue, component, format string) ([]byte, error) {

	if tissue == "bnblood" {
		fmt.Println("TOM graph not available for bnblood")
		return nil, errors.New("TOM graph not available for bnblood")
	}

	if component == "nodes" {
		component = "Nodes"
	} else {
		component = "Edges"
	}

	fun := "getTOMGraph" + component
	args := "tissue='" + tissue + "'"

	key, err := R.Call(pkg, fun, args)
	if err != nil {
		fmt.Println("Could not get TOM graph nodes")
		return nil, err
	}
	return R.Get(key, format)

}

func Get(key, filetype string) ([]byte, error) {
	return R.Get(key, filetype)
}

func GetCohorts() ([]string, error) {
	fun := "getCohorts"
	args := ""

	key, err := R.Call(pkg, fun, args)
	if err != nil {
		fmt.Println("Could not get cohorts", err)
		return []string{}, errors.Wrap(err, "get cohort r call fail"+pkg+fun+args)
	}

	resp, err := R.Get(key, "json")

	if err != nil {
		return []string{}, errors.Wrap(err, "could not get cohort"+key)
	}

	cohorts := []string{}

	err = json.Unmarshal(resp, &cohorts)
	if err != nil {
		return []string{}, errors.Wrap(err, "get cohort unmarshal error"+string(resp))
	}
	return cohorts, nil

}
