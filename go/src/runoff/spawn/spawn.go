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
package spawn

import (
  "bytes"
  "fmt"
  "os"
  "os/exec"
  "syscall"
)

func ProcExec(args []string) error {
  fmt.Printf("%v\n", args)
  var stdout, stderr bytes.Buffer
  proc := exec.Command(args[0], args[1:]...)
  proc.Env = os.Environ()
  // Safeguard for exec.Command clean exit
  // For more info check: https://medium.com/@ganeshmaharaj/clean-exit-of-golangs-exec-command-897832ac3fa5
  proc.SysProcAttr = &syscall.SysProcAttr{
    Pdeathsig: syscall.SIGTERM,
  }
  proc.Stdout = &stdout
  proc.Stderr = &stderr
  if err := proc.Run(); err != nil {
    fmt.Println(stderr.String())
    return err
  }
  fmt.Println(stdout.String())
  return nil
}
