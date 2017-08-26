package mysql

import (
	"database/sql"
	"errors"

	"strings"

	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/saturn4er/barkup"
)

type MysqlClient struct {
	tableName string
	DB        *sqlx.DB
	dsn       string
}

func (c *MysqlClient) JoinQueries(queries []string) string {
	return strings.Join(queries, ";")
}

func (c *MysqlClient) Backup(path string) (string, error) {
	dsn, err := mysql.ParseDSN(c.dsn)
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

func (c *MysqlClient) ApplyMigration(migration *DBMigration, down bool) error {
	var query string
	if down {
		query = migration.DownScript
	} else {
		query = migration.UpScript
	}
	_, err := c.DB.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (c *MysqlClient) InsertMigration(migration *DBMigration) error {
	now := time.Now()
	migration.AppliedAt = &now
	_, err := c.DB.NamedExec("INSERT INTO `"+c.tableName+"` (`num`, `name`, `up_script`, `down_script`,`applied_at`) VALUES (:num, :name, :up_script, :down_script, :applied_at);", migration)
	return err
}

func (c *MysqlClient) GetAppliedMigrations() ([]DBMigration, error) {
	result := []DBMigration{}
	err := c.DB.Select(&result, "SELECT `num`, `name`, `up_script`, `down_script`, `applied_at` FROM `"+c.tableName+"` ORDER BY `applied_at` ASC")
	if err == sql.ErrNoRows {
		return result, nil
	}
	return result, err
}

func (c *MysqlClient) RemoveMigration(migration *DBMigration) error {
	_, err := c.DB.Exec("DELETE FROM `"+c.tableName+"` WHERE `num`=?", migration.Number)
	return err
}

func (c *MysqlClient) PrepareTransactionsTable() error {
	exists, err := c.dbVersionTableExists()
	if err != nil {
		return err
	}
	if !exists {
		err = c.createDBVersionTable()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *MysqlClient) dbVersionTableExists() (bool, error) {
	var tableName string
	err := c.DB.QueryRow("SHOW TABLES LIKE '" + c.tableName + "'").Scan(&tableName)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, errors.New("can't check if db version table exists: " + err.Error())
	}
	return true, nil
}
func (c *MysqlClient) createDBVersionTable() error {
	t := (&MysqlQueryBuilder{}).CreateTable(c.tableName, func(table CreateTableGenerator) {
		table.Column("num", "int").NotNull().Primary()
		table.Column("name", "text").NotNull()
		table.Column("up_script", "text").NotNull()
		table.Column("down_script", "text").NotNull()
		table.Column("applied_at", "datetime").NotNull()
	})
	_, err := c.DB.Exec(t.Sql())
	if err != nil {
		return errors.New("can't create db version table: " + err.Error())
	}
	return nil
}
func NewClient(dsn, transactionsTableName string) (DBClient, error) {
	result := new(MysqlClient)
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
	result.DB = db
	result.dsn = dsn
	result.tableName = transactionsTableName
	return result, nil
}
