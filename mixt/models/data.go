package models

import (
	"strconv"
	"strings"

	"github.com/fjukstad/kvik/dataset"
	"github.com/fjukstad/kvik/utils"
)

var workerAddr string

func Init() (*dataset.Dataset, error) {
	ip := "localhost"
	port := ":8888"
	filename := "scripts/script.r"
	var err error
	var d *dataset.Dataset
	d, workerAddr, err = dataset.Init(ip, port, filename)
	return d, err
}

func Add(d *dataset.Dataset, a, b int) ([]string, error) {
	command := "add(3,2)"
	resp, err := d.Call(command)
	if err != nil {
		return []string{}, err
	}

	response := utils.PrepareResponse(resp)
	return response, nil
}

func Hist(d *dataset.Dataset) (string, error) {
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

func GetGenes(d *dataset.Dataset) ([]string, error) {
	command := "getGenes()"
	resp, err := d.Call(command)
	if err != nil {
		return []string{""}, err
	}

	response := utils.PrepareResponse(resp)
	return response, nil
}
