package generator

import (
	"strings"
)

type CSV struct {
	columnSeparator string
	dataIterator    chan []string
	builder         strings.Builder
}

func NewCSVGenerator(columnSeparator string, iter chan []string) Generator {
	return &CSV{
		columnSeparator: columnSeparator,
		dataIterator:    iter,
	}
}

func (csv *CSV) Generate() (err error) {
	for row := range csv.dataIterator {
		csv.buildRow(row)
		csv.builder.WriteString("\n")
	}
	return
}

func (csv *CSV) buildRow(data []string) (err error) {
	for index, column := range data {
		csv.builder.WriteString(column)
		if index == len(data)-1 {
			return
		}
		csv.builder.WriteString(csv.columnSeparator)
	}
	return
}

func (csv CSV) String() string {
	return csv.builder.String()
}
