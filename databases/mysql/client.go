package mysql

import "github.com/saturn4er/migratego/types"

type Client struct {
}

func (c *Client) PrepareTransactionsTable() error {
	return nil
}

func (c *Client) ApplyMigration(migration types.Migration, down bool) error {
	return nil
}

func NewClient(dsn, transactionsTableName string) types.Client {
	result := new(Client)
	return result
}
