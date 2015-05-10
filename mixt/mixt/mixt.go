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
	//filename := "scripts/script.r"
	// d, workerAddr, err = dataset.RequestNewWorker(ip, port, filename)
	//workerAddr = "tcp://localhost:5000"
	//fmt.Println("connecting to worker at", workerAddr)
	//d, err = dataset.ConnectToRunningWorker(workerAddr)
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

func GetAllModules(gene string) ([]string, error) {
	resp, err := komp.Rpc("mixt/R/getAllModules", "", "json")
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
	Name       string
	HeatmapUrl string
	Genes      []Gene
	Signatures []Signature
	Url        string
}

type Gene struct {
	Name        string
	Correlation float64
	K           float64
	Kin         float64
	Updown      string
}

type Signature struct {
	Name       string
	Size       int
	PValue     float64
	Common     []string
	PValueUp   float64
	ComnonUp   []string
	PValueDown float64
	CommonDown []string
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
		return nil, err
	}

	var modules []Module
	for _, moduleName := range moduleNames {
		m, err := GetModule(moduleName, tissue)
		if err != nil {
			fmt.Println("cannot get module", moduleName, err)
			return nil, err
		}
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
		return Module{}, err
	}

	genes, url, err := GetGeneList(name, tissue)
	if err != nil {
		return Module{}, err
	}

	module := Module{name, heatmapUrl, genes, nil, url}
	return module, nil

}

func GetGeneList(module, tissue string) (genes []Gene, url string, err error) {
	resp, err := komp.Rpc("mixt/R/getGeneList", `{"tissue": `+"\""+tissue+"\""+`,
					"module":`+"\""+module+"\""+`}`, "csv")

	if err != nil {
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

	return genes, url, nil
}
