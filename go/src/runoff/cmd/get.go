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

// getCmd represents the get command
var getCmd = &cobra.Command{
  Use:   "get",
  Short: "Execute runoff model algorithm",
  Long: `Compute the following steps:
(1) Allocate cover area and cover weigth values to
each polygon output by the link command.
(2) Get the target basin's threshold-runoff value.
Get is the third step in the sequence of
model computations.`,
  Run: func(cmd *cobra.Command, args []string) {
              m := model.NewModel(os.Getenv("SIOSE_DB"), os.Getenv("RUNOFF_SCHEMA"))
	      fmt.Println("Allocating cover area and cover weight...")
              if allocateErr := m.Allocate(); allocateErr != nil {
                fmt.Println("Command exited with error:", allocateErr)
              }
	      fmt.Println("Computing p0 an n weighted mean values...")
              if wmeanErr := m.ComputeWMean(); wmeanErr != nil {
                fmt.Println("Command exited with error:", wmeanErr)
              }
	      p0, p0Err := m.GetP0()
	      if p0Err != nil {
	        fmt.Println("Command exited with error:", p0Err)
	      } else {
	        fmt.Println("Basin P0 = ", p0)
	      }
	},
}

func init() {
  rootCmd.AddCommand(getCmd)
}
