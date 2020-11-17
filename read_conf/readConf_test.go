package readConf

import (
	"bufio"
	"os"
	"testing"
)

func TestParseSecName(t *testing.T) {
	var cfg config
	var e error
	cfg.init()
	secName := defaultSection
	var line []byte
	line = []byte("[paths]")
	secName, e = parseSecName(line, &cfg)

	ex1 := "paths"

	var ex2 error

	if secName != ex1 {
		t.Errorf("ex1 wrong, '%s'", secName)
	}
	if e != ex2 {
		t.Errorf("ex2 wrong, '%s'", e)
	}

}

func TestParsekey(t *testing.T) {
	line := "app_mode = development"
	keyName, offset, err := parseKeyName(string(line))
	ex1, ex2 := "app_mode", 10
	var ex3 error
	if keyName != ex1 {
		t.Errorf("ex1 wrong")
	}
	if offset != ex2 {
		t.Errorf("ex2 wrong")
	}
	if err != ex3 {
		t.Errorf("ex3 wrong")
	}
}

func TestParseValue(t *testing.T) {
	line := "app_mode = development"
	value, err := parseValue(string(line[10:]))
	ex1 := "development"
	var ex2 error
	if value != ex1 {
		t.Errorf("ex1 wrong")
	}
	if err != ex2 {
		t.Errorf("ex2 wrong")
	}

}

func TestParse(t *testing.T) {
	var e error
	f, err := os.Open("./testdata/test.ini")
	if e != err {
		t.Errorf("open file filed")
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	_, erro := parse(buf)
	if erro != e {
		t.Errorf("parse filed")
	}

}
