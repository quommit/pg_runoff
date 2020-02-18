/*
Copyright © 2020 José Tomás Navarro Carrión

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
  "fmt"
  "github.com/spf13/cobra"
  "runoff/etl"
)

// inspectCmd represents the inspect command
var inspectCmd = &cobra.Command{
  Use:   "inspect <filename>",
  Short: "Inspect metadata found in input shapefile",
  Long: `Inspect the input shapefile and show metadata
such as field names, field data types, geometry type,
spatial reference system and spatial extension.
Given a shapefile filename.shp, command usage is:

runoff inspect filename.shp

This command is a wrapper around GDAL/OGR ogrinfo utility.`,
  Args: cobra.ExactValidArgs(1),
  Run: func(cmd *cobra.Command, args []string) {
              input  := etl.Shp(args[0])
              if err := input.Summary(); err != nil {
                fmt.Println("Command exited with error:", err)
              }
	},
}

func init() {
  rootCmd.AddCommand(inspectCmd)
}
