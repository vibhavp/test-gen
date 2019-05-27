package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/vibhavp/test-gen/gen"
	"github.com/vibhavp/test-gen/step"
	yaml "gopkg.in/yaml.v3"
	"io/ioutil"
)

var fileName = flag.String("file", "", "test file name")

func init() {
	flag.Parse()
}

func main() {
	if *fileName == "" {
		flag.Usage()
		os.Exit(128)
	}

	file, err := os.Open(*fileName)
	if err != nil {
		log.Fatal(err)
	}

	decoder := yaml.NewDecoder(file)
	test := struct{ Test step.Test }{}

	if err = decoder.Decode(&test); err != nil {
		log.Fatal(err)
	}

	template, err := ioutil.ReadFile("template.py")
	if err != nil {
		log.Fatalf("Error while reading template.py: %v", template)
	}

	sw := bytes.NewBuffer(template)
	sw.WriteString(fmt.Sprintf("# Generated from %s\n", *fileName))
	ctxt := gen.NewGenContext(sw)
	if err = ctxt.GenTest(&test.Test); err != nil {
		log.Fatal(err)
	}

	fmt.Print(sw.String())
}
