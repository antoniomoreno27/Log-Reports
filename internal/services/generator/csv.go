package generator

import (
	"bufio"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/antoniomoreno27/logs-LS/internal/services/logger"
)

type CSV struct {
	name            string
	path            string
	columnSeparator string
	wg              *sync.WaitGroup
	dataIterator    chan []string
	builder         strings.Builder
}

const (
	fileExtention = ".csv"
)

func NewCSVGenerator(filename, path, columnSeparator string, wg *sync.WaitGroup, iter chan []string) Generator {
	return &CSV{
		name:            filename,
		path:            path,
		wg:              wg,
		columnSeparator: columnSeparator,
		dataIterator:    iter,
	}
}

func (csv *CSV) buildFile() (err error) {
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

func (csv *CSV) Generate() error {
	defer csv.wg.Done()
	filename := csv.name + fileExtention
	path := path.Join(csv.path, filename)
	logger.Warnf("Creating file %s", csv.path+"_"+csv.name)
	os.Remove(path)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	buffer := bufio.NewWriterSize(file, 4096)
	logger.Warnf("Writing data into %s", csv.name+fileExtention)
	err = csv.buildFile()
	if err != nil {
		return err
	}
	buffer.WriteString(csv.String())
	err = buffer.Flush()
	if err != nil {
		return err
	}
	logger.Warnf("Finished writing file %s", csv.name+fileExtention)
	return nil
}
