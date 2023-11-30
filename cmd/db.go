package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func dbCmd(a *appState) *cobra.Command {
	db := &cobra.Command{
		Use:     "db",
		Short:   "Manage the database",
		Aliases: []string{"d"},
		Example: strings.TrimSpace(fmt.Sprintf(`$ %s db [command]`, appName)),
	}
	return db
}
