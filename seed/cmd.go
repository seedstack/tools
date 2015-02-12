package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"path/filepath"
	"flag"
	"io/ioutil"
)

// Transformation is a strutucture representating a set 
// of procedure to apply on a source code directory
type Transformation struct {
	Pre []string
	Proc []Procedure
}

// Procedure is a function call with a method name and 
// its parameters
type Procedure struct {
	Name string
	Params []string
}

func main() {
	flag.Parse()
	
	tdfPath, errFilePath := filepath.Abs("./tdf.yml")
	if errFilePath != nil {
		fmt.Println("Error constructing the file path.", errFilePath)
	}

    if flag.Arg(0) != "" {
		tdfPath, errFilePath = filepath.Abs(flag.Arg(0))
		if errFilePath != nil {
			fmt.Println("Error constructing the file path.", errFilePath)
		}
	}
	
	transf := parseTdf(tdfPath)
	path := "../test"
	processFiles(walkthroughDir(path), transf)
}

func parseTdf(tdfPath string) []Transformation {
	dat, err := ioutil.ReadFile(tdfPath)
	if err != nil {
		panic(err)
	}
	var transf []Transformation
	err = yaml.Unmarshal(dat, &transf)
	if err != nil {
		panic(err)
	}
//	fmt.Printf("Value: %#v\n", transf)
	fmt.Println("Parsed the tranformation description file.")
	return transf
}
