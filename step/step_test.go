package step_test

import (
	"testing"
	"os"
	yaml "gopkg.in/yaml.v3"
	"github.com/vibhavp/test-gen/step"
	"github.com/k0kubun/pp"
)

func TestSimple(t *testing.T) {
	f, err := os.Open("test.yml")
	if err != nil {
		t.Fatal(err)
	}

	d := yaml.NewDecoder(f)

	s := &struct {Test step.Test}{step.Test{}}
	if err != d.Decode(s) {
		t.Fatal(err)
	}

	pp.Println(*s)
}
