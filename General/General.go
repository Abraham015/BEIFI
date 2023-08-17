package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type General struct{}

func (g *General) NumberCities(file *os.File) int {
	cities := 0
	name := file.Name()
	var number string

	for _, c := range name {
		switch c {
		case '1', '2', '3', '4', '5', '6', '7', '8', '9', '0':
			number += string(c)
		}
	}

	cities, _ = strconv.Atoi(number)
	return cities
}

func (g *General) TypeProblem(file *os.File) string {
	var t string

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			fmt.Println(err)
			break
		}

		if strings.TrimSpace(line) == "TYPE : TSP" {
			t = "TSP"
			break
		} else if strings.TrimSpace(line) == "TYPE: ATSP" {
			t = "ATSP"
			break
		}

		if err == io.EOF {
			break
		}
	}

	return t
}

