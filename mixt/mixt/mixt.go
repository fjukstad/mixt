package mixt

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
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
		fmt.Println(session, err)
		return "", err
	}

	plotUrl := session.Graphics

	// Since we don't have a dns thing for opencpu
	plotUrl = strings.Replace(plotUrl, "opencpu", "docker0.bci.mcgill.ca", -1)

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
	HeatmapUrl       string
	Genes            []Gene
	EnrichmentScores []EnrichmentScore
	Url              string
}

type Gene struct {
	Name        string
	Correlation float64
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

	fmt.Println(moduleNames)

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
			m := Module{moduleNames[i], "", nil, nil, ""}
			resChan <- m
		}(i)
	}

	for range moduleNames {
		m := <-resChan
		modules = append(modules, m)
	}

	return modules, nil
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

	module := Module{name, heatmapUrl, genes, scores, url}
	return module, nil
}

func GetGeneList(module, tissue string) (genes []Gene, url string,
	err error) {
	session, err := komp.Call("mixt/R/getGeneList", `{"tissue": `+"\""+tissue+"\""+`,
					"module":`+"\""+module+"\""+`}`)

	if err != nil {
		fmt.Println("SESSION ERROR")
		fmt.Println(session, err)
		return nil, "", err
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

		g := Gene{Name: name,
			Correlation: 0,
			K:           0,
			Kin:         0,
			Updown:      updown}

		genes = append(genes, g)
	}

	url = session.GetUrl("csv")
	url = strings.Replace(url, "opencpu", "docker0.bci.mcgill.ca", -1)

	return genes, url, nil
}

type EnrichmentScore struct {
	Set          string  `json:"set"`
	Name         string  `json:"_row"`
	Size         int     `json:"sig.size"`
	UpDownCommon int     `json:"updn.common"`
	UpDownPvalue float64 `json:"updn.p"`
	UpCommon     int     `json:"up.common"`
	UpPvalue     float64 `json:"up.p"`
	DownCommon   int     `json:"dn.common"`
	DownPvalue   float64 `json:"dn.p"`
}

func GetEnrichmentScores(module, tissue string) (scores []EnrichmentScore,
	err error) {

	session, err := komp.Call("mixt/R/getEnrichmentScores",
		`{"tissue": `+"\""+tissue+"\""+`,
					"module":`+"\""+module+"\""+`}`)

	if err != nil {
		fmt.Println(session, err)
		return nil, err
	}

	resp, err := session.GetResult(komp, "json")
	if err != nil {
		fmt.Println(resp, err)
		return nil, err
	}

	res := []byte(resp)

	err = json.Unmarshal(res, &scores)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return scores, nil

}
