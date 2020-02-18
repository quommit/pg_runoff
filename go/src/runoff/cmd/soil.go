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

//soil_field var stores the target field name given in the "field" flag
var soil_field string

// soilCmd represents the soil command
var soilCmd = &cobra.Command{
  Use:   "soil <filename>",
  Short: "Import soil data from shapefile",
  Long: `Import shapefile polygons with a
soil modifier attribute in the domain {-2, -1, 0, 1}.
Given a shapefile soil.shp with soil modifier
field named mod_value, typical command usage is:

runoff import soil -f mod_value soil.shp

This command is a wrapper around GDAL/OGR ogr2ogr and psql utilities.`,
  Args: cobra.ExactValidArgs(1),
  Run: func(cmd *cobra.Command, args []string) {
              input  := etl.SoilShp(args[0])
              output := etl.NewSIOSE(os.Getenv("SIOSE_DB"), os.Getenv("RUNOFF_SCHEMA"))
	      fmt.Println("Fetching...")
              if fetchErr := input.Fetch(output); fetchErr != nil {
                fmt.Println("Command exited with error:", fetchErr)
              }
              fmt.Println("Merging...")
              if mergeErr := input.Merge(output, soil_field); mergeErr != nil {
                fmt.Println("Command exited with error:", mergeErr)
              }
	},
}

func init() {
  importCmd.AddCommand(soilCmd)
  // Soil modifier field name flag definition
  soilCmd.Flags().StringVarP(&soil_field, "field", "f", "soil_mod", "Name of attribute that holds soil modifier values.")
}
