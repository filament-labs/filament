package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	// Read JSON
	data, err := os.ReadFile("constants.json")
	if err != nil {
		panic(err)
	}

	var constants map[string]string
	if err := json.Unmarshal(data, &constants); err != nil {
		panic(err)
	}

	// --- Generate Go constants ---
	goFile, err := os.Create("../../internal/app/constants_gen.go")
	if err != nil {
		panic(err)
	}
	defer goFile.Close()

	fmt.Fprintln(goFile, "package app")
	fmt.Fprintln(goFile, "const (")

	for key, value := range constants {
		constName := key
		fmt.Fprintf(goFile, "\t%s = \"%s\"\n", constName, value)
	}

	fmt.Fprintln(goFile, ")")

	// --- Generate TypeScript constants ---
	tsFilePath := "../../web/src/util/constants.ts" // adjust path
	tsFile, err := os.Create(tsFilePath)
	if err != nil {
		panic(err)
	}
	defer tsFile.Close()

	fmt.Fprintln(tsFile, "export const CONSTANTS = {")
	for key, value := range constants {
		fmt.Fprintf(tsFile, "  %s: '%s',\n", key, value)
	}
	fmt.Fprintln(tsFile, "}")
}
