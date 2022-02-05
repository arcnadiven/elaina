/*
Copyright © 2022 arcnadiven@github.com

*/
package cmd

import (
	"github.com/arcnadiven/elaina/backstore"
	"github.com/arcnadiven/elaina/server"
	"github.com/arcnadiven/elaina/tracelog"

	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start will run a grpc server for CSI client",
	Long: `Elaina is a out-of-tree persistentvolume plugin base on MySQL in kubernetes

This application is a grpc server for CSI client, 
it will bound the persistentvolume to a localpath
and insert a record to MySQL`,
	Run: Start,
}

var (
	// TODO: add log config later
	logfile string

	dbUser, dbPassword, dbHost, dbPort, dbName string

	dbArgs map[string]string
)

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// log file flags
	startCmd.Flags().StringVar(&logfile, "log-file", "", "output all log to this file if non-empty")

	// database flags
	startCmd.Flags().StringVar(&dbUser, "db-user", "root", "database username for establish connection")
	startCmd.Flags().StringVar(&dbPassword, "db-passwd", "", "database user password for establish connection")
	startCmd.Flags().StringVar(&dbHost, "db-host", "localhost", "database server ip for establish connection")
	startCmd.Flags().StringVar(&dbPort, "db-port", "3306", "database server port for establish connection")
	startCmd.Flags().StringVar(&dbName, "db-name", "", "database name for establish connection")
	startCmd.Flags().StringToStringVar(&dbArgs, "db-args", nil, "database connect arguments for establish connection")
}

func Start(cmd *cobra.Command, args []string) {
	bl := tracelog.NewBaseLogger(logfile)
	bl.Infoln(
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
		dbName,
		dbArgs,
	)
	conf := &backstore.ConnConfig{
		Username:     dbUser,
		Password:     dbPassword,
		Host:         dbHost,
		Port:         dbPort,
		DataBaseName: dbName,
		ConnArgs:     dbArgs,
	}
	sqlCli, err := backstore.NewSQLClient(bl, conf)
	if err != nil {
		bl.Errorln(err)
		return
	}
	sqlOpera := backstore.NewStoreOperator(sqlCli)
	srv := server.NewElainaServer(bl, sqlOpera)
	if err := srv.RunGRPCServer(); err != nil {
		bl.Errorln(err)
		return
	}
}
