package main

import (
	"fmt"
	"net/http"

	"bitbucket.org/vdumeaux/mixt/mixt/mixt"
)

func main() {

	mixt.Init("opencpu", "biopsy@mcgill", "van-mi-ka-al")

	tissues := []string{"blood", "biopsy", "nblood", "bnblood"}

	for _, tissueA := range tissues {

		fmt.Println(tissueA)
		modulesA, err := mixt.GetModules(tissueA)
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, moduleA := range modulesA {
			_, err := http.Get("http://docker0.bci.mcgill.ca:8004/modules/" + tissueA + "/" + moduleA.Name)
			fmt.Println(tissueA, moduleA.Name, err)
		}

		http.Get("http://docker0.bci.mcgill.ca:8004/clinical-comparison/" + tissueA + "/roi")
		http.Get("http://docker0.bci.mcgill.ca:8004/clinical-comparison/" + tissueA + "/eigengene")

		for _, tissueB := range tissues {
			modulesB, err := mixt.GetModules(tissueB)
			_, err = http.Get("http://docker0.bci.mcgill.ca:8004/tissues/" + tissueA + "/" + tissueB)
			fmt.Println(tissueA, tissueB, err)

			http.Get("http://docker0.bci.mcgill.ca:8004/tissues/" + tissueA + "/" + tissueB + "/eigengene")
			http.Get("http://docker0.bci.mcgill.ca:8004/tissues/" + tissueA + "/" + tissueB + "/overlap")
			http.Get("http://docker0.bci.mcgill.ca:8004/tissues/" + tissueA + "/" + tissueB + "/roi")
			http.Get("http://docker0.bci.mcgill.ca:8004/tissues/" + tissueA + "/" + tissueB + "/patientrank")

			for _, moduleA := range modulesA {
				for _, moduleB := range modulesB {
					fmt.Println(tissueA, tissueB, moduleA.Name, moduleB.Name)
					http.Get("http://docker0.bci.mcgill.ca:8004/compare/" + tissueA + "/" + tissueB + "/" + moduleA.Name + "/" + moduleB.Name)
				}
			}

		}
	}
}
