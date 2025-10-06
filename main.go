package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// help message
func usage() {
	exe := "./main"
	fmt.Fprintf(os.Stderr, "Usage: %s <input.csv> <output.jl>\n", exe)
	fmt.Fprintf(os.Stderr, "Example: %s housesInput.csv housesOutput.jl\n", exe)
}

// CSV to JSON types
func parseCell(s string) any {
	// nulls
	if s == "" {
		return nil
	}

	// booleans
	switch strings.ToLower(s) {
	case "true":
		return true
	case "false":
		return false
	}

	// integer
	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		return i
	}

	// float
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}
	return s
}

func convertCSVtoJSON(inPath, outPath string) error {
	inFile, err := os.Open(inPath)
	if err != nil {
		return fmt.Errorf("opening input: %w", err)
	}
	defer inFile.Close()

	// reads csv file
	reader := csv.NewReader(inFile)
	reader.TrimLeadingSpace = true

	// read header row (col names)
	header, err := reader.Read()
	if err != nil {
		return fmt.Errorf("reading header: %w", err)
	}
	// create new error message for if theres no info in header
	if len(header) == 0 {
		return errors.New("empty header row")
	}
	for i, h := range header {
		header[i] = strings.TrimSpace(h)
	}

	// Prepare output
	outFile, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("creating output: %w", err)
	}

	// batches writes
	writer := bufio.NewWriter(outFile)
	defer writer.Flush() // sends anything in the buffer to the file
	defer outFile.Close()

	// turn rows into JSOn
	rowNum := 1 // everything past header
	// loop through rows
	for {
		record, err := reader.Read() // csv row
		// no file
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("reading row %d: %w", rowNum+1, err)
		}
		rowNum++

		// build go map
		obj := make(map[string]any, len(header)) // colname, any type val
		// loop through columns of that row
		for i := range header {
			var val string       // "" default
			if i < len(record) { // only read record[i] if it exists
				val = strings.TrimSpace(record[i]) // take care of trailing whitespace
			}
			// stores parsed json value in colname in obj map
			obj[header[i]] = parseCell(val)
		}

		// 1 JSON obj per line in the same order as the header
		// go maps don't preserve order so need to print each pair
		// {
		if _, err := writer.WriteString("{"); err != nil {
			return err
		}
		for i, key := range header { // loop the cols
			if i > 0 {
				if err := writer.WriteByte(','); err != nil { // commas between pairs
					return err
				}
			}
			// key:
			if _, err := fmt.Fprintf(writer, "%q:", key); err != nil {
				return err
			}

			// conv go type to json type of val
			value, err := json.Marshal(obj[key])
			if err != nil {
				return fmt.Errorf("marshal value for key %q: %w", key, err)
			}
			// write val
			if _, err := writer.Write(value); err != nil {
				return err
			}
		}
		// }\n
		if _, err := writer.WriteString("}\n"); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	if len(os.Args) != 3 {
		usage()
		os.Exit(2)
	}
	inPath := os.Args[1]
	outPath := os.Args[2]

	if err := convertCSVtoJSON(inPath, outPath); err != nil {
		// error handling
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
