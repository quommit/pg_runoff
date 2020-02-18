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
  "fmt"
  "github.com/spf13/cobra"
  "os"
  "runoff/etl"
)

// basinCmd represents the basin command
var basinCmd = &cobra.Command{
  Use:   "basin <filename>",
  Short: "Import basin polygon from shapefile",
  Long: `Import the shapefile polygon covering
the drainage area which the runoff model will
be computed for. Given a shapefile basin.shp,
typical command usage is:

runoff import basin basin.shp

This command is a wrapper around GDAL/OGR ogr2ogr and psql utilities.`,
  Args: cobra.ExactValidArgs(1),
  Run: func(cmd *cobra.Command, args []string) {
              input  := etl.BasinShp(args[0])
              output := etl.NewSIOSE(os.Getenv("SIOSE_DB"), os.Getenv("RUNOFF_SCHEMA"))
	      fmt.Println("Fetching...")
              if fetchErr := input.Fetch(output); fetchErr != nil {
                fmt.Println("Command exited with error:", fetchErr)
              }
              fmt.Println("Merging...")
              if mergeErr := input.Merge(output, ""); mergeErr != nil {
                fmt.Println("Command exited with error:", mergeErr)
              }
	},
}

func init() {
	importCmd.AddCommand(basinCmd)
}
