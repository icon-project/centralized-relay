package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func db(app *appState) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "db",
		Aliases: []string{"d"},
		Short:   "Manage the database",
		Example: strings.TrimSpace(fmt.Sprintf(`
		$ %s db [command],
		$ %s db [command] --help
		`, appName, appName)),
	}

	return cmd
}
