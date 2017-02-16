package migrates

import (
	"database/sql"
	"errors"

	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Migrator struct {
	databaseVersionTable string
	db                   *sqlx.DB
	migrations           []Migration
}
type Version struct {
	Number     int `db:"num"`
	Name       string
	UpScript   string    `db:"up_script"`
	DownScript string    `db:"down_script"`
	AppliedAt  time.Time `db:"applied_at"`
}

func (v *Version) SameAsMigration(m *Migration) bool {
	if v.Number != m.Number{
		return false
	}
	if v.Name != m.Name {
		return false
	}
	if v.UpScript != m.UpScript {
		return false
	}
	if v.DownScript != m.DownScript {
		return false
	}
	return true
}

func (d *Migrator) prepareDBVersionTable() error {
	exists, err := d.dbVersionTableExists()
	if err != nil {
		return err
	}
	if !exists {
		err = d.createDBVersionTable()
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Migrator) createDBVersionTable() error {
	_, err := d.db.Exec("CREATE TABLE `" + d.databaseVersionTable + "` (" +
		" `num` INT NOT NULL," +
		" `name` VARCHAR(100) NOT NULL," +
		" `up_script` TEXT(0) NOT NULL," +
		" `down_script` TEXT(0) NOT NULL," +
		" `applied_at` DATETIME NOT NULL," +
		" PRIMARY KEY (`num`));")
	if err != nil {
		return errors.New("can't create db version table: " + err.Error())
	}
	return nil
}
func (d *Migrator) getAllVersions() ([]Version, error) {
	result := []Version{}
	err := d.db.Select(&result, "SELECT `num`, `name`, `up_script`, `down_script`, `applied_at` FROM `"+d.databaseVersionTable+"` ORDER BY `applied_at` ASC")
	if err == sql.ErrNoRows {
		return result, nil
	}
	return result, err
}
func (d *Migrator) getCurrentVersion() (*Version, error) {
	result := &Version{}
	err := d.db.Get(&result, "SELECT * FROM `"+d.databaseVersionTable+"` ORDER BY `applied_at` DESC")
	if err != nil && err == sql.ErrNoRows {
	    return nil, nil
	}
	return result, err
}
//func (d *Migrator) downgradeVersion(v *Version) (*Version, error) {
//
//	err := d.db.Get(&result, "SELECT * FROM `"+d.databaseVersionTable+"` ORDER BY `applied_at` DESC")
//	return result, err
//}
func (d *Migrator) dbVersionTableExists() (bool, error) {
	var tableName string
	err := d.db.QueryRow("SHOW TABLES LIKE '" + d.databaseVersionTable + "'").Scan(&tableName)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, errors.New("can't check if db version table exists: " + err.Error())
	}
	return true, nil
}

func NewMigrator(dbVersionTable, dsn string, migrations []Migration) (*Migrator, error) {
	d, err := mysql.ParseDSN(dsn)
	if err != nil {
		return nil, errors.New("bad dsn: " + err.Error())
	}
	d.MultiStatements = true
	d.ParseTime = true
	dsn = d.FormatDSN()
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, errors.New("can't connect to database: " + err.Error())
	}
	err = db.Ping()
	if err != nil {
		return nil, errors.New("can't connect to database: " + err.Error())
	}
	result := &Migrator{
		db:                   db,
		databaseVersionTable: dbVersionTable,
		migrations:           migrations,
	}
	return result, nil
}

