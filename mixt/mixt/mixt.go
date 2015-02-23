package mixt

import (
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

func Hist() (string, error) {
	command := "plt()"
	resp, err := d.Call(command)
	if err != nil {
		return "", err
	}

	response := utils.PrepareResponse(resp)
	newWorker := strings.Replace(workerAddr, "tcp", "http", -1)

	split := strings.Split(newWorker, ":")
	id, err := strconv.Atoi(split[2])

	if err != nil {
		return "", err
	}

	httpPort := id + 1

	httpURL := split[0] + ":" + split[1] + ":" + strconv.Itoa(httpPort)

	baseURL := httpURL + "/"
	url := baseURL + response[0]
	return url, nil
}

func GetGenes() ([]string, error) {
	command := "getGenes()"
	resp, err := d.Call(command)
	if err != nil {
		return []string{""}, err
	}

	response := utils.PrepareResponse(resp)
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

	url, err := Hist()
	if err != nil {
		return Module{}, err
	}

	module := Module{name, url, nil, nil}
	return module, nil

}
