package mixt

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/fjukstad/kvik/dataset"
	"github.com/fjukstad/kvik/utils"
)

var workerAddr string
var d *dataset.Dataset

func Init() error {
	ip := "localhost"
	port := ":8888"
	filename := "scripts/script.r"
	var err error
	d, workerAddr, err = dataset.Init(ip, port, filename)
	return err
}

func Add(a, b int) ([]string, error) {
	command := "add(3,2)"
	resp, err := d.Call(command)
	if err != nil {
		return []string{}, err
	}

	response := utils.PrepareResponse(resp)
	return response, nil
}

func Heatmap(tissue, module string) (string, error) {
	command := "heatmap(\"" + tissue + "\",\"" + module + "\")"
	resp, err := d.Call(command)
	if err != nil {
		return "", err
	}

	response := utils.PrepareResponse(resp)
	url, err := getUrl(response[0])

	if err != nil {
		return "", err
	}

	url = strings.TrimSuffix(url, ".png")

	return url, nil
}

func getUrl(ending string) (string, error) {
	newWorker := strings.Replace(workerAddr, "tcp", "http", -1)

	split := strings.Split(newWorker, ":")
	id, err := strconv.Atoi(split[2])

	if err != nil {
		return "", err
	}

	httpPort := id + 1

	httpURL := split[0] + ":" + split[1] + ":" + strconv.Itoa(httpPort)
	baseURL := httpURL + "/"
	url := baseURL + ending
	return url, nil
}

func GetGenes() ([]string, error) {
	command := "getAllGenes()"
	resp, err := d.Call(command)
	if err != nil {
		return []string{""}, err
	}

	response := utils.PrepareResponse(resp)
	url, err := getUrl(response[0])

	if err != nil {
		return []string{}, err
	}

	listResp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return []string{}, err
	}

	reader := csv.NewReader(listResp.Body)

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

	fmt.Println(url)

	return genes, nil
}

func GetAllModules(gene string) ([]string, error) {
	command := "getAllModules(\"" + gene + "\")"

	resp, err := d.Call(command)
	if err != nil {
		return []string{""}, err
	}
	response := utils.PrepareResponse(resp)
	for i, r := range response {
		response[i] = strings.Trim(r, "\"")
	}
	return response, nil
}

func GetTissues() ([]string, error) {
	command := "getTissues()"
	resp, err := d.Call(command)
	if err != nil {
		return []string{""}, err
	}
	response := utils.PrepareResponse(resp)
	for i, r := range response {
		response[i] = strings.Trim(r, "\"")
	}
	return response, nil
}

type Module struct {
	Name       string
	HeatmapUrl string
	Genes      []Gene
	Signatures []Signature
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

func GetModules(tissue string) ([]Module, error) {
	command := "getModules(\"" + tissue + "\")"
	resp, err := d.Call(command)
	if err != nil {
		return []Module{}, err
	}
	var modules []Module
	response := utils.PrepareResponse(resp)
	for _, r := range response {
		name := strings.Trim(r, "\"")
		modules = append(modules, Module{name, "", nil, nil})
	}
	return modules, nil
}

func GetModule(name string, tissue string) (Module, error) {

	url, err := Heatmap(tissue, name)
	if err != nil {
		return Module{}, err
	}

	genes, err := GetGeneList(name, tissue)
	if err != nil {
		return Module{}, err
	}

	module := Module{name, url, genes, nil}
	return module, nil

}

func GetGeneList(module, tissue string) ([]Gene, error) {
	command := "getGeneList(\"" + tissue + "\",\"" + module + "\")"
	resp, err := d.Call(command)
	if err != nil {
		fmt.Println(err)
		return []Gene{}, err
	}
	response := utils.PrepareResponse(resp)
	url, _ := getUrl(response[0])

	listResp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return []Gene{}, err
	}

	reader := csv.NewReader(listResp.Body)

	var genes []Gene
	line := 0
	for {
		record, err := reader.Read()

		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return []Gene{}, nil
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

	return genes, nil
}
