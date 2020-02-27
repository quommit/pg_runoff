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
package etl

import (
  "fmt"
  "path/filepath"
  "runoff/spawn"
  "strings"
)

const RAW_BASIN string = "basin_input"
const MODEL_BASIN string = "basin"
const MODEL_BASIN_COL string = ""
const RAW_SLOPE string = "slope_input"
const MODEL_SLOPE string = "slope"
const MODEL_SLOPE_COL string = "slope_mod"
const RAW_SOIL string = "soil_input"
const MODEL_SOIL string = "soil"
const MODEL_SOIL_COL string = "soil_mod"

func getBaseNoExt(filename string) string {
  return strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
}

type ShpReader interface {
  Summary() error
}

type Shp string

func (filename Shp) Summary() error {
  ogrinfo(string(filename))
  return nil
}

// ogrinfo_cmd var stores the ogrinfo command line text with placeholders
// needed to list fields and related metadata from the input data source.
// Please note that the double semicolon ";;" is used as a separator for
// eventually splitting the command line text into its constituent options
// and arguments.
var ogrinfo_cmd = "ogrinfo;;-so;;%v;;%v"

func ogrinfo(filename string) error {
  layername := getBaseNoExt(filename)
  proc := fmt.Sprintf(ogrinfo_cmd, filename, layername)
  args := strings.Split(proc, ";;")
  if err := spawn.ProcExec(args); err != nil {
    return err
  }
  return nil
}

type SIOSE struct {
  name, schema string
}

func NewSIOSE(name string, schema string) *SIOSE {
  s := SIOSE{name, schema}
  return &s
}

type ShpConverter interface {
  Fetch(db *SIOSE) error
  Merge(db *SIOSE, sourcefield string) error
}

type BasinShp string

func (filename BasinShp) Fetch(db *SIOSE) error {
  if err := ogr2ogr(string(filename), db, RAW_BASIN); err != nil {
    return err
  }
  return nil
}

func (filename BasinShp) Merge(db *SIOSE, sourcefield string) error {
  if err := feedModel(db, RAW_BASIN, sourcefield, MODEL_BASIN, MODEL_BASIN_COL); err != nil {
    return err
  }
  return nil
}

type SlopeShp string

func (filename SlopeShp) Fetch(db *SIOSE) error {
  if err := ogr2ogr(string(filename), db, RAW_SLOPE); err != nil {
    return err
  }
  return nil
}

func (filename SlopeShp) Merge(db *SIOSE, sourcefield string) error {
  if err := feedModel(db, RAW_SLOPE, sourcefield, MODEL_SLOPE, MODEL_SLOPE_COL); err != nil {
    return err
  }
  return nil
}

type SoilShp string

func (filename SoilShp) Fetch(db *SIOSE) error {
  if err := ogr2ogr(string(filename), db, RAW_SOIL); err != nil {
    return err
  }
  return nil
}

func (filename SoilShp) Merge(db *SIOSE, sourcefield string) error {
  if err := feedModel(db, RAW_SOIL, sourcefield, MODEL_SOIL, MODEL_SOIL_COL); err != nil {
    return err
  }
  return nil
}

// ogr2ogr_cmd var stores the ogr2ogr command line text with placeholders
// needed to import data into the target database. Please note that the
// double semicolon ";;" is used as a separator for eventually splitting
// the command line text into its constituent options and arguments.
var ogr2ogr_cmd = "ogr2ogr;;-f;;PostgreSQL;;PG:%v;;-overwrite;;-nln;;%v.%v;;-nlt;;PROMOTE_TO_MULTI;;%v;;%v"

// ogrcnstr var holds the OGR connection string to the target database
var ogrcnstr = "dbname='%v' user='postgres'"

func ogr2ogr(filename string, db *SIOSE, targettable string) error {
  layername := getBaseNoExt(filename)
  cnstr := fmt.Sprintf(ogrcnstr, db.name)
  proc := fmt.Sprintf(ogr2ogr_cmd, cnstr, db.schema, targettable, filename, layername)
  args := strings.Split(proc, ";;")
  if err := spawn.ProcExec(args); err != nil {
    return err
  }
  return nil
}

// psql_cmd var stores the psql command line text with placeholders
// needed to feed data into the model, which means inserting rows into
// the basin, slope or soil tables. Please note that the double semicolon
// ";;" is used as a separator for eventually splitting the command line
// text into its constituent options and arguments.
var psql_cmd = "psql;;-d;;%v;;-U;;postgres;;-c;;%v;;-c;;%v;;-1"

// truncate_cmd var holds the SQL statement that removes all rows from
// the basin, slope or soil model tables.
var truncate_cmd = "TRUNCATE %v.%v RESTART IDENTITY"

// insert_cmd var holds the SQL statement that selects rows from
// an input data table and inserts rows into the corresponding
// basin, slope or soil model tables.
var insert_cmd = "INSERT INTO %v.%v (%v, geom) SELECT %q::smallint, st_transform(wkb_geometry, 4258) FROM %v.%v"

// insert_geom_only_cmd var holds the SQL statement that selects rows from
// an input data table and just inserts geometry into a model table.
var insert_geom_only_cmd = "INSERT INTO %v.%v (geom) SELECT st_transform(wkb_geometry, 4258) FROM %v.%v"

func feedModel(db *SIOSE, sourcetable string, sourcefield string, targettable string, targetcol string) error {
  sqltruncate := fmt.Sprintf(truncate_cmd, db.schema, targettable)
  var sqlinsert string
  if strings.TrimSpace(sourcefield) != "" && strings.TrimSpace(targetcol) != "" {
    sqlinsert = fmt.Sprintf(insert_cmd, db.schema, targettable, targetcol, strings.ToLower(sourcefield), db.schema, sourcetable)
  } else {
    sqlinsert = fmt.Sprintf(insert_geom_only_cmd, db.schema, targettable, db.schema, sourcetable)
  }
  proc := fmt.Sprintf(psql_cmd, db.name, sqltruncate, sqlinsert)
  args := strings.Split(proc, ";;")
  if err := spawn.ProcExec(args); err != nil {
    return err
  }
  return nil
}
