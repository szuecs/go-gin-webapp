package main

import "testing"

func TestGetStringFromArgs(t *testing.T) {
	var args map[string]interface{}
	args = make(map[string]interface{}, 1)
	expected := "string"
	args["test"] = expected
	value := GetStringFromArgs(args, "test", "alternate")
	if value != expected {
		t.FailNow()
	}
	value = GetStringFromArgs(args, "np", "alternate")
	if value != "alternate" {
		t.FailNow()
	}
}

func TestGetIntFromStr(t *testing.T) {
	value := "1"
	expected := 1
	integer := GetIntFromStr(value)
	if integer != expected {
		t.FailNow()
	}
}

func TestGetBoolFromStr(t *testing.T) {
	value := "true"
	actual := GetBoolFromStr(value)
	if actual != true {
		t.FailNow()
	}
}

func TestGetIntFromArgs(t *testing.T) {
	var args map[string]interface{}
	args = make(map[string]interface{}, 1)
	expected := "1"
	args["test"] = expected
	value := GetIntFromArgs(args, "test", 0)
	if value != 1 {
		t.FailNow()
	}
	value = GetIntFromArgs(args, "np", 0)
	if value != 0 {
		t.FailNow()
	}
}

func TestGetBoolFromArgs(t *testing.T) {
	var args map[string]interface{}
	args = make(map[string]interface{}, 1)
	expected := "true"
	args["test"] = expected
	value := GetBoolFromArgs(args, "test", false)
	if value != true {
		t.FailNow()
	}
	value = GetBoolFromArgs(args, "np", true)
	if value != true {
		t.FailNow()
	}
}

func TestConvertMaps(t *testing.T) {
	input := "A=false B='foo'"
	aMap := ConvertMaps(input)
	if aMap["A"] != "false" || aMap["B"] != "'foo'" {
		t.FailNow()
	}
}
