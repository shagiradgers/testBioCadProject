package database

import (
	"database/sql"
	_ "modernc.org/sqlite"
)

const driverName = "sqlite"

type Config struct {
	FilePath string `yaml:"filePath"`
}

type Database struct {
	FilePath string
}

func NewDatabase(config Config) *Database {
	return &Database{
		config.FilePath,
	}
}

func (d *Database) Init() error {
	db, err := sql.Open(driverName, d.FilePath)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	if err := db.Close(); err != nil {
		return err
	}

	return nil
}

// IsFileChecked return true if file checked, else false
func (d *Database) IsFileChecked(filePath string) (bool, error) {
	db, err := sql.Open(driverName, d.FilePath)
	if err != nil {
		return false, err
	}
	defer db.Close()

	rows, err := db.Query("select * from checked_files where file_path = $1 ", filePath)
	if err != nil {
		return false, err
	}

	if rows.Next() {
		return true, nil
	}
	return false, nil
}

// AddCheckedFile add checked file to db
func (d *Database) AddCheckedFile(filepath string, status string) error {
	db, err := sql.Open(driverName, d.FilePath)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("insert into checked_files(file_path, status) values($1, $2)", filepath, status)
	if err != nil {
		return err
	}

	return nil
}

// AddModel add InputModels to db
func (d *Database) AddModel(model InputModel) error {
	db, err := sql.Open(driverName, d.FilePath)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("insert into models VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)",
		model.Mqtt, model.Invid, model.UnitGuid, model.MsgId, model.Text, model.Context, model.Class, model.Level,
		model.Area, model.Addr, model.Block, model.DataType, model.Bit, model.InvertBit)
	if err != nil {
		return err
	}
	return nil
}

// GetModels return models. If limit = -1, return all models
func (d *Database) GetModels(unitGuid string, limit int) ([]*InputModel, error) {
	db, err := sql.Open(driverName, d.FilePath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var rows *sql.Rows
	if limit < 0 {
		rows, err = db.Query("select * from models")
	} else {
		rows, err = db.Query("select * from models limit $1", limit)
	}

	if err != nil {
		return nil, err
	}

	models := make([]*InputModel, 0)

	for rows.Next() {
		model := new(InputModel)
		err = rows.Scan(&model.Mqtt, &model.Invid, &model.UnitGuid, &model.MsgId, &model.Text, &model.Context,
			&model.Class, &model.Level, &model.Area, &model.Addr, &model.Block, &model.DataType, &model.Bit,
			&model.InvertBit)
		if err != nil {
			return nil, err
		}

		models = append(models, model)
	}

	return models, nil
}
