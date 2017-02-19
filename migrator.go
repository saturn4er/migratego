package migratego

import (
	"database/sql"
	"errors"

	"strings"

	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/keighl/barkup"
	"github.com/saturn4er/migratego/types"
)

type Migrator struct {
	databaseVersionTable string
	driver               string
	dsn                  string
	db                   *sqlx.DB
	migrations           []types.Migration
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

	t := getDriverQueryBuilder(d.driver).CreateTable(d.databaseVersionTable, func(table types.CreateTableGenerator) {
		table.Column("num", "int").NotNull().Primary()
		table.Column("name", "text").NotNull()
		table.Column("up_script", "text").NotNull()
		table.Column("down_script", "text").NotNull()
		table.Column("applied_at", "datetime").NotNull()
	})
	_, err := d.db.Exec(t.Sql())
	if err != nil {
		return errors.New("can't create db version table: " + err.Error())
	}
	return nil
}
func (d *Migrator) getAppliedMigrations() ([]types.Migration, error) {
	result := []types.Migration{}
	err := d.db.Select(&result, "SELECT `num`, `name`, `up_script`, `down_script`, `applied_at` FROM `"+d.databaseVersionTable+"` ORDER BY `applied_at` ASC")
	if err == sql.ErrNoRows {
		return result, nil
	}
	return result, err
}
func (d *Migrator) getLatestMigration() (*types.Migration, error) {
	result := &types.Migration{}
	err := d.db.Get(&result, "SELECT * FROM `"+d.databaseVersionTable+"` ORDER BY `applied_at` DESC")
	if err != nil && err == sql.ErrNoRows {
		return nil, nil
	}
	return result, err
}
func (d *Migrator) applyMigration(migration *types.Migration, down bool) error {
	var query string
	if down {
		query = migration.DownScript
	} else {
		query = migration.UpScript
	}
	_, err := d.db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
func (d *Migrator) deleteMigration(m *types.Migration) error {
	_, err := d.db.Exec("DELETE FROM `"+d.databaseVersionTable+"` WHERE `num`=?", m.Number)
	return err
}
func (d *Migrator) addMigration(m *types.Migration) error {
	now := time.Now()
	m.AppliedAt = &now
	_, err := d.db.NamedExec("INSERT INTO `"+d.databaseVersionTable+"` (`num`, `name`, `up_script`, `down_script`,`applied_at`) VALUES (:num, :name, :up_script, :down_script, :applied_at);", m)
	return err
}
func (d *Migrator) backupDatabase(path string) (string, error) {
	dsn, err := mysql.ParseDSN(d.dsn)
	if err != nil {
		return "", errors.New("error parsing dsn: " + err.Error())
	}
	addr := strings.Split(dsn.Addr, ":")
	var host = "127.0.0.1"
	var port = "3306"
	if len(addr) > 0 {
		host = addr[0]
	}
	if len(addr) > 1 {
		port = addr[1]
	}
	db := &barkup.MySQL{
		Host:     host,
		Port:     port,
		DB:       dsn.DBName,
		User:     dsn.User,
		Password: dsn.Passwd,
	}

	export := db.Export()
	if export.Error != nil {
		return "", export.Error
	}
	cpErr := export.To(path, nil)
	if cpErr != nil {
		return export.Filename(), errors.New("can't copy backup to dst path:" + cpErr.Error())
	}
	return export.Filename(), nil
}
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

func newMigrator(dbVersionTable, driver, dsn string, migrations []types.Migration) (*Migrator, error) {
	d, err := mysql.ParseDSN(dsn)
	if err != nil {
		return nil, errors.New("bad dsn: " + err.Error())
	}
	d.MultiStatements = true
	d.ParseTime = true
	dsn = d.FormatDSN()
	db, err := sqlx.Open(driver, dsn)
	if err != nil {
		return nil, errors.New("can't connect to database: " + err.Error())
	}
	err = db.Ping()
	if err != nil {
		return nil, errors.New("can't connect to database: " + err.Error())
	}
	result := &Migrator{
		driver:               driver,
		db:                   db,
		dsn:                  dsn,
		databaseVersionTable: dbVersionTable,
		migrations:           migrations,
	}
	return result, nil
}
