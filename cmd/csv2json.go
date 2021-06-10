package cmd

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(convert)
}

var convert = &cobra.Command{
	Use:     "convert",
	Aliases: []string{"con"},
	Short:   "To convert file from CSV to JSON",
	Run: func(cmd *cobra.Command, args []string) {
		fileData, err := getFileData(args, cmd)

		if err != nil {
			exitGracefully(err)
		}

		if _, err := checkIfValidFile(fileData.filepath); err != nil {
			exitGracefully(err)
		}
		writerChannel := make(chan map[string]string)
		done := make(chan bool)

		go processCsvFile(fileData, writerChannel)
		go writeJSONFile(fileData.filepath, writerChannel, done, fileData.pretty)

		<-done
	},
}

func getFileData(args []string, cmd *cobra.Command) (inputFile, error) {
	if len(args) < 1 {
		return inputFile{}, errors.New("A filepath argument is required")
	}

	separator, _ := cmd.Flags().GetString("separator")
	pretty, _ := cmd.Flags().GetBool("pretty")

	fileLocation := args[0]

	if !(separator == "comma" || separator == "semicolon") {
		return inputFile{}, errors.New("Only comma or semicolon separators are allowed")
	}

	return inputFile{fileLocation, separator, pretty}, nil
}

func checkIfValidFile(filename string) (bool, error) {
	//check if file is CSV
	if fileExtention := filepath.Ext(filename); fileExtention != ".csv" {
		return false, fmt.Errorf("File %s does not exist", filename)
	}

	//check if file does exist
	if _, err := os.Stat(filename); err != nil && os.IsNotExist(err) {
		return false, fmt.Errorf("File %s does not exist", filename)
	}
	return true, nil
}

func processLine(headers []string, dataList []string) (map[string]string, error) {
	if len(dataList) != len(headers) {
		return nil, errors.New("Line doesn't match headers format. Skipping")
	}

	recordMap := make(map[string]string)

	for i, name := range headers {
		recordMap[name] = dataList[i]
	}

	return recordMap, nil
}

func processCsvFile(fileData inputFile, writerChannel chan<- map[string]string) {
	file, err := os.Open(fileData.filepath)

	check(err)

	defer file.Close()

	//Get Headers
	var headers, line []string
	reader := csv.NewReader(file)

	if fileData.separator == "semicolon" {
		reader.Comma = ';'
	}

	headers, err = reader.Read()
	check(err)

	for {
		line, err = reader.Read()

		if err == io.EOF {
			close(writerChannel)
			break
		} else if err != nil {
			exitGracefully(err)
		}

		record, err := processLine(headers, line)

		if err != nil {
			fmt.Printf("Line: %sError: %s\n", line, err)
			continue
		}

		writerChannel <- record
	}
}

func createStringWriter(csvPath string) func(string, bool) {
	jsonDir := filepath.Dir(csvPath)
	jsonName := fmt.Sprintf("%s.json", strings.TrimSuffix(filepath.Base(csvPath), ".csv"))
	finalLocation := fmt.Sprintf("%s/%s", jsonDir, jsonName)

	f, err := os.Create(finalLocation)
	check(err)

	return func(data string, close bool) {
		_, err := f.WriteString(data)
		check(err)

		if close {
			f.Close()
		}
	}
}

func getJSONFunc(pretty bool) (func(map[string]string) string, string) {
	var jsonFunc func(map[string]string) string
	var breakLine string
	if pretty {
		breakLine = "\n"
		jsonFunc = func(record map[string]string) string {
			jsonData, _ := json.MarshalIndent(record, "   ", "   ")
			return "   " + string(jsonData)
		}
	} else {
		breakLine = ""
		jsonFunc = func(record map[string]string) string {
			jsonData, _ := json.Marshal(record)
			return string(jsonData)
		}
	}
	return jsonFunc, breakLine
}

func writeJSONFile(csvPath string, writerChannel <-chan map[string]string, done chan<- bool, pretty bool) {
	writeString := createStringWriter(csvPath)
	jsonFunc, breakLine := getJSONFunc(pretty)

	fmt.Println("Writing JSON file...")

	writeString("["+breakLine, false)
	first := true
	for {
		record, more := <-writerChannel
		if more {
			if !first {
				writeString(","+breakLine, false)
			} else {
				first = false
			}

			jsonData := jsonFunc(record)
			writeString(jsonData, false)
		} else {
			writeString(breakLine+"]", true)
			fmt.Println("Completed!")
			done <- true
			break
		}
	}
}

func exitGracefully(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}

func check(e error) {
	if e != nil {
		exitGracefully(e)
	}
}
