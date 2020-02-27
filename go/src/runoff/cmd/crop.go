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
  "os"
  "runoff/model"
)

// cropCmd represents the crop command
var cropCmd = &cobra.Command{
  Use:   "crop",
  Short: "Crop runnoff model geometries",
  Long: `Crop slope and soil polygons using the basin
boundary as the cutting line. Crop is the first step
in the sequence of model computations.`,
  Run: func(cmd *cobra.Command, args []string) {
              m := model.NewModel(os.Getenv("SIOSE_DB"), os.Getenv("RUNOFF_SCHEMA"))
	      fmt.Println("Cropping...")
              if cropErr := m.Crop(); cropErr != nil {
                fmt.Println("Command exited with error:", cropErr)
              }
  },
}

func init() {
  rootCmd.AddCommand(cropCmd)
}
