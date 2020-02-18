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

//slope_field var stores the target field name given in the "field" flag
var slope_field string

// slopeCmd represents the slope command
var slopeCmd = &cobra.Command {
  Use:   "slope <filename>",
  Short: "Import slope data from shapefile",
  Long: `Import shapefile polygons with a
slope modifier attribute in the domain {-3, 3}.
Given a shapefile slope.shp with slope modifier
field named mod_value, typical command usage is:

runoff import slope -f mod_value slope.shp

This command is a wrapper around GDAL/OGR ogr2ogr and psql utilities.`,
  Args: cobra.ExactValidArgs(1),
  Run: func(cmd *cobra.Command, args []string) {
              input  := etl.SlopeShp(args[0])
              output := etl.NewSIOSE(os.Getenv("SIOSE_DB"), os.Getenv("RUNOFF_SCHEMA"))
	      fmt.Println("Fetching...")
              if fetchErr := input.Fetch(output); fetchErr != nil {
                fmt.Println("Command exited with error:", fetchErr)
              }
              fmt.Println("Merging...")
              if mergeErr := input.Merge(output, slope_field); mergeErr != nil {
                fmt.Println("Command exited with error:", mergeErr)
              }
  },
}

func init() {
  importCmd.AddCommand(slopeCmd)
  // Slope modifier field name flag definition
  slopeCmd.Flags().StringVarP(&slope_field, "field", "f", "slope_mod", "Name of attribute that holds slope modifier values.")
}
