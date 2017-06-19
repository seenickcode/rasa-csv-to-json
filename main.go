package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
)

type RasaOutput struct {
	Data *RasaNLUData `json:"rasa_nlu_data"`
}

type RasaNLUData struct {
	CommonExamples []*Intent `json:"common_examples"`
}

type Intent struct {
	Intent   string    `json:"intent"`
	Text     string    `json:"text"`
	Entities []*Entity `json:"entities"`
}

type Entity struct {
	Name  string `json:"entity"`
	Start int    `json:"start"`
	End   int    `json:"end"`
	Value string `json:"value"`
}

const (
	EntityStartingCol = 2
	EntitySlots       = 3
	EntityNumCols     = 4
)

func main() {

	// read all of the records in CSV in to an slice
	csvFile, err := os.Open("file.csv")
	if err != nil {
		log.Fatal(err)
	}
	lines, err := csv.NewReader(bufio.NewReader(csvFile)).ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// create an empty slice to receive the data
	intents := []*Intent{}
	fmt.Println("\nOriginal data in CSV format\n")

	// loop through CSV records, create a slice and append it to slice array
	for i, line := range lines {
		if i == 0 {
			// skip header line
			continue
		}
		entities := []*Entity{}

		// collect entities
		for slotNdx := 0; slotNdx < EntitySlots; slotNdx++ {
			entityNdx := EntityStartingCol + (slotNdx * EntityNumCols)
			if len(line) > entityNdx && len(line[entityNdx]) > 0 {
				start, err := strconv.Atoi(line[entityNdx+1])
				if err != nil {
					log.Fatal(err)
				}
				end, err := strconv.Atoi(line[entityNdx+2])
				if err != nil {
					log.Fatal(err)
				}
				e := &Entity{
					Name:  line[entityNdx+0],
					Start: start,
					End:   end,
					Value: line[entityNdx+3],
				}
				entities = append(entities, e)
			}

		}
		i := &Intent{
			Intent:   line[0],
			Text:     line[1],
			Entities: entities,
		}
		fmt.Println(i)
		intents = append(intents, i)
	}

	output := &RasaOutput{
		Data: &RasaNLUData{
			CommonExamples: intents,
		},
	}

	fmt.Println("\nOriginal data in JSON format\n")

	// print the reformatted struct as JSON
	// cs := spew.ConfigState{
	// 	Indent:                  "\t",
	// 	DisableMethods:          true,
	// 	DisablePointerMethods:   true,
	// 	DisablePointerAddresses: true,
	// 	SpewKeys:                false,
	// }
	// cs.Dump(output)

	res1a, _ := json.MarshalIndent(output, "", "  ")
	fmt.Println("\nOriginal data in JSON format\n")
	fmt.Printf("%s\n", res1a)

	outf, err := os.Create("out.json")
	defer outf.Close()
	n2, err := outf.Write(res1a)
	fmt.Printf("wrote %d bytes\n", n2)

}
