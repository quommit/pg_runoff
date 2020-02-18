/*
Copyright © 2020 José Tomás Navarro Carrión <jt.navarro@ua.es>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
  "github.com/spf13/cobra"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
  Use:   "import",
  Short: "Import runoff model data",
  Long: `Import geometry and attribute needed
for runoff model computation. To get
help on import subcommands run:

runoff import --help`,
}

func init() {
	rootCmd.AddCommand(importCmd)
}
