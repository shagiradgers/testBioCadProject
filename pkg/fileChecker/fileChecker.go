package fileChecker

import (
	"encoding/csv"
	"fmt"
	"github.com/gocarina/gocsv"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testBaioCadProject/pkg/database"
)

type Checker struct {
	Db *database.Database
	Config
}

// NewChecker create new checker
func NewChecker(db *database.Database, config Config) *Checker {
	initCsvReader()

	return &Checker{
		db,
		config,
	}
}

// initCsvReader add configuration to gocsv reader
func initCsvReader() {
	gocsv.SetCSVReader(
		func(file io.Reader) gocsv.CSVReader {
			r := csv.NewReader(file)
			r.Comma = '\t'
			r.FieldsPerRecord = -1
			return r
		},
	)
}

// writeToDoc write all models with some unitGuid in .doc file
func (c *Checker) writeToDoc(unitGuid string) error {
	path := c.OutputPath + "/" + unitGuid + ".doc"

	models, err := c.Db.GetModels(unitGuid, -1)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write([]byte(``))
	if err != nil {
		return err
	}

	file.Write([]byte(unitGuid))

	for _, model := range models {
		text := fmt.Sprintf("mqtt: %v\nInvid: %v\nMsg_id: %v\nText: %v\nContext: %v\nClass: %v\nLevel: %v\n"+
			"Area: %v\nAddr: %v\nBlock: %v\nData_type: %v\nBit: %v\nInvert_bit:%v",
			model.Mqtt, model.Invid, model.MsgId, model.Text, model.Context, model.Class, model.Level,
			model.Area, model.Addr, model.Block, model.DataType, model.Bit, model.InvertBit)
		text += "\n-----------------------------\n"
		file.Write([]byte(text))
	}

	return nil
}

// CheckFile checking file and add there to db
func (c *Checker) CheckFile(filePath string) error {
	in, err := c.Db.IsFileChecked(filePath)
	if err != nil {
		return err
	}

	// if file was checked
	if in {
		return nil
	}

	dataFile, err := os.OpenFile(filePath, os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer dataFile.Close()

	// check, that file is tsv
	if tmp := strings.Split(dataFile.Name(), "."); tmp[len(tmp)-1] != "tsv" {
		return c.Db.AddCheckedFile(filePath, "file is not tsv")
	}

	models := make([]*database.InputModel, 0)
	uniqueModels := map[string]struct{}{}
	status := "ok"

	err = gocsv.UnmarshalFile(dataFile, &models)
	if err != nil {
		if err == gocsv.ErrEmptyCSVFile {
			status = "file is empty"
		} else {
			status = err.Error()
		}
	}

	err = c.Db.AddCheckedFile(filePath, status)
	if err != nil {
		return err
	}

	for _, model := range models {
		err = c.Db.AddModel(*model)
		if err != nil {
			return err
		}

		if _, ok := uniqueModels[model.UnitGuid]; !ok {
			uniqueModels[model.UnitGuid] = struct{}{}
		}
	}

	for unitGuid, _ := range uniqueModels {
		err = c.writeToDoc(unitGuid)
		if err != nil {
			return err
		}
	}

	return nil
}

// RunChecker running checker :)
func (c *Checker) RunChecker() error {
	err := filepath.Walk(c.RootPath,
		func(wPath string, info os.FileInfo, err error) error {
			if wPath == c.RootPath {
				return nil
			}

			if info.IsDir() {
				return nil
			}

			if wPath != c.RootPath {
				if err := c.CheckFile(wPath); err != nil {
					return err
				}
			}
			return nil
		},
	)
	if err != nil {
		return err
	}
	return nil
}
