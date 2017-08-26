package migratego

import "gopkg.in/urfave/cli.v1"

var cliCommands = []cli.Command{}

func RunToolCli(m *migrateApplication, args []string) error {
	tool := cli.NewApp()
	tool.HelpName = "migratego"
	client, err := m.getDriverClient()
	if err != nil {
		return err
	}
	err = client.PrepareTransactionsTable()
	if err != nil {
		return err
	}
	tool.Version = "1.0.0"
	tool.Usage = "Tool to manipulate with database versions"
	tool.Metadata = map[string]interface{}{
		"client": client,
		"app":    m,
	}
	tool.Commands = cliCommands
	tool.Run(args)
	return nil
}
