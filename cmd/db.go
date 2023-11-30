package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/icon-project/centralized-relay/relayer/lvldb"
	"github.com/spf13/cobra"
)

func dbCmd(a *appState) *cobra.Command {
	var db *lvldb.LVLDB
	dbCMD := &cobra.Command{
		Use:     "db",
		Short:   "Manage the database",
		Aliases: []string{"d"},
		Example: strings.TrimSpace(fmt.Sprintf(`$ %s db [command]`, appName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			db, err = lvldb.NewLvlDB(filepath.Join(defaultHome, defaultDBName))
			if err != nil {
				return err
			}
			return nil
		},
	}

	pruneCmd := &cobra.Command{
		Use:   "prune",
		Short: "Prune the database",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Pruning the database...")
			db.ClearStore()
		},
	}

	messagesCmd := &cobra.Command{
		Use:   "messages",
		Short: "Get messages stored in the database",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO:
			fmt.Println("Getting messages stored in the database...")
		},
	}

	dbCMD.AddCommand(pruneCmd)
	dbCMD.AddCommand(messagesCmd)

	return dbCMD
}
