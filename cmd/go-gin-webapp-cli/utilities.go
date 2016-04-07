package main

import (
	"fmt"
	"strconv"
	"strings"
)

//GetStringFromArgs returns a string from the cli arguments.
func GetStringFromArgs(arguments map[string]interface{}, name, _default string) string {
	in := arguments[name]
	if in == nil {
		return _default
	}
	return in.(string)
}

//GetIntFromStr returs an integer from a given string
func GetIntFromStr(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		fmt.Printf("ERR: int conversion failed for %s: err: %s", s, err)
	}
	return int(i)
}

//GetBoolFromStr gets a boolean from a string
func GetBoolFromStr(s string) bool {
	b, err := strconv.ParseBool(s)
	if err != nil {
		fmt.Printf("ERR: bool conversion failed for %s: err: %s", s, err)
	}
	return b
}

//GetIntFromArgs gets an integer from arguments by name
func GetIntFromArgs(arguments map[string]interface{}, name string, _default int) int {
	if arguments[name] == nil {
		return _default
	}
	valStr := arguments[name].(string)
	return GetIntFromStr(valStr)
}

//GetBoolFromArgs gets a boolean from arguments by name
func GetBoolFromArgs(arguments map[string]interface{}, name string, _default bool) bool {
	if arguments[name] == nil {
		return _default
	}
	valStr := arguments[name].(string)
	return GetBoolFromStr(valStr)
}

//ConvertMaps creates a map from a string representation of the maps itself.
//TODO: document or replace with something more
func ConvertMaps(input string) map[string]string {
	var labels = make(map[string]string)
	if input != "" {
		labelsArray := strings.Split(input, " ")
		for i := 0; i < len(labelsArray); i++ {
			if strings.Contains(labelsArray[i], "=") { //just to ensure it's actually a k,v pair
				kvArr := strings.Split(labelsArray[i], "=")
				for i := 0; i < len(kvArr); i++ {
					if i%2 == 0 {
						labels[kvArr[i]] = ""
					} else {
						labels[kvArr[i-1]] = kvArr[i]
					}
				}
			}

		}
	}
	return labels
}
