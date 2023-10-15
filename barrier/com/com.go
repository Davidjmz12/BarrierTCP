package com

import (
	"bufio"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"log"
	"os"
)

func CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func PrintMsgLog(color func(format string, a ...interface{}) string, args ...interface{}) {
	log.Println(color(fmt.Sprintln(args...)))
}

func ParsePeers(path string) (lines []string) {
	file, err := os.Open(path)
	CheckError(err)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func AddOneVector(vector []int, me int) []int {
	result := make([]int, len(vector))
	for i := 0; i < len(vector); i++ {
		if i == me-1 {
			result[i] = vector[i] + 1
		} else {
			result[i] = vector[i]
		}
	}
	return result
}

func AllTrue(slice []bool) bool {
	for _, value := range slice {
		if !value {
			return false
		}
	}
	return true
}

func TryDecode(obj interface{}, ret interface{}) bool {
	decoder, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:      &ret,
		ErrorUnused: true, // Reports unused fields as an error
		ErrorUnset:  true,
	})

	err := decoder.Decode(obj)

	return err == nil
}
