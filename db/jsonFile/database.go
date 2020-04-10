package jsonFile

import (
	"dapi/db"
	"dapi/interfaces"
	"encoding/json"
	"os"
)

type config struct {
	interfaces.IBaseDatabase
	fileName string
}

func NewDB(fileName string) *config {
	c := new(config)
	c.fileName = fileName
	c.IBaseDatabase = db.NewBaseDatabase()

	return c
}

//Save state to JSON config file
func (c *config) Save() error {
	configFile, err := os.OpenFile(c.fileName, os.O_WRONLY, os.ModePerm)

	if err != nil {
		return err
	}

	defer func() {
		_ = configFile.Close()
	}()

	jsonEncoder := json.NewEncoder(configFile)
	jsonEncoder.SetIndent("", "\t")
	err = jsonEncoder.Encode(c.IBaseDatabase)

	return err
}

//Load Config from JSON file
func (c *config) Load() error {
	configFile, err := os.Open(c.fileName)

	if err != nil {
		return err
	}

	defer func() {
		_ = configFile.Close()
	}()

	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(c.IBaseDatabase)

	if err != nil {
		return err
	}

	return nil
}
