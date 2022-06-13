package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	// parse args
	args := os.Args[1:]
	if len(args) < 1 || len(args) > 2 {
		fmt.Println("this command expects 1 or 2 arguments: csvFile [jsonFile]")
		os.Exit(1)
	}
	csvFileName := args[0]
	var jsonFileName string
	if len(args) == 2 {
		jsonFileName = args[1]
	} else {
		fileParts := strings.Split(csvFileName, ".")
		fileParts = fileParts[:len(fileParts)-1]
		jsonFileName = strings.Join(fileParts, ".") + ".json"
	}
	fmt.Printf("reading from %s\n", csvFileName)
	fmt.Printf("writing to %s\n", jsonFileName)

	// read scanner
	csvFile, err := os.Open(csvFileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer csvFile.Close()
	csvReader := csv.NewReader(csvFile)

	// writer
	jsonFile, err := os.Create(jsonFileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer jsonFile.Close()
	jsonWriter := bufio.NewWriter(jsonFile)
	defer jsonWriter.Flush()

	read := true
	var headers []string
	line := 0
	for read {
		record, err := csvReader.Read()
		if record == nil && err == io.EOF {
			read = false
			break
		} else if err != nil {
			fmt.Printf("error reading csv: %s\n", err.Error())
			os.Exit(1)
		} else if record == nil {
			fmt.Printf("record on line %d came back nil? no error\n", line)
			os.Exit(1)
		}

		if headers == nil || len(headers) == 0 {
			// TODO: this doesn't play well with colliding headers (e.g. if "user name" and "user_name" are both in headers)
			headers = []string{}
			for _, header := range record {
				headers = append(headers, strings.ReplaceAll(header, " ", "_"))
			}
			continue
		}

		m, err := makeMap(headers, record)
		if err != nil {
			fmt.Printf("error making map of line %d: %s\n", line, err.Error())
			os.Exit(1)
		}
		jsonData, err := json.Marshal(m)
		if err != nil {
			fmt.Printf("error marshaling json of line %d: %s\n", line, err.Error())
		}
		jsonWriter.WriteString(string(jsonData) + "\n")
		jsonWriter.Flush()
	}
}

func makeMap(headers, data []string) (map[string]string, error) {
	if len(headers) != len(data) {
		return nil, fmt.Errorf("mismatched lengths: %d headers and %d data", len(headers), len(data))
	}

	m := map[string]string{}
	for i, header := range headers {
		m[header] = data[i]
	}

	return m, nil
}
