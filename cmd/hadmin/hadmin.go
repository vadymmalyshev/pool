package main

import (
	"flag"
	internalAdmin "git.tor.ph/hiveon/pool/internal/admin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

var cmdRoot = &cobra.Command{
	Use:   "hadmin",
	Short: "Hiveon Admin server",
	Run:   internalAdmin.RunServer,
}

var cmdAdmin = &cobra.Command{
	Use:   "admin",
	Short: "Add/remove admin rights to user",
	TraverseChildren: true,
	Run:   internalAdmin.AddAdmin,
}

var cmdAdd = &cobra.Command{
	Use:   "add",
	Short: "Add admin rights for user",
	Run:   internalAdmin.AddAction,
}

var cmdRemove = &cobra.Command{
	Use:   "remove",
	Short: "Remove admin rights for user",
	Run:   internalAdmin.RemoveAction,
}

var cmdMigrate = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate db",
	Run:   internalAdmin.DoMigrate,
}

func init() {
	flag.Parse()
	cmdAdmin.AddCommand(cmdAdd)
	cmdAdmin.AddCommand(cmdRemove)
	cmdRoot.AddCommand(cmdMigrate)
	cmdRoot.AddCommand(cmdAdmin)
}

func main() {
	if err := cmdRoot.Execute(); err != nil {
		logrus.Infof("can't run admin server: %s", err)
		os.Exit(1)
	}

	os.Exit(0)
}
