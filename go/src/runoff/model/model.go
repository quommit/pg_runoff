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
package model

import (
  "fmt"
  "runoff/spawn"
  "strings"
)

const MODEL_SLOPE_CROP string = "slope_crop"
const MODEL_SOIL_CROP string = "soil_crop"
const MODEL_LINK = "link"
const MODEL_C_LEAF_NODES = "cover_leaf_nodes"
const MODEL_C_ALLOC = "cover_allocation"
const MODEL_WMEAN = "weighted_mean"
const MODEL_BASIN_P0 = "basin_p0"

type Model struct {
  name, schema string
}

func NewModel(name string, schema string) *Model {
  m := Model{name, schema}
  return &m
}

type Builder interface {
  Crop() error
  Link() error
  Alloc() error
  ComputeWMean() error
  GetP0() error
}

func (db Model) Crop() error {
  if err := refreshCropViews(&db); err != nil {
    return err
  }
  return nil
}

func (db Model) Link() error {
  if err := refreshLinkView(&db); err != nil {
    return err
  }
  return nil
}

func (db Model) Allocate() error {
  if err := refreshCoverViews(&db); err != nil {
    return err
  }
  return nil
}

func (db Model) ComputeWMean() error {
  if err := refreshWMeanView(&db); err != nil {
    return err
  }
  return nil
}

func (db Model) GetP0() (string, error) {
  r, err := selectBasinP0(&db)
  if err != nil {
    return "", err
  }
  return r, nil
}

var refresh_cmd = "REFRESH MATERIALIZED VIEW %v.%v"
var select_scalar_cmd = "SELECT DISTINCT %v FROM %v.%v LIMIT 1"

// psql_crop_cmd var stores the psql command line text with placeholders
// needed to refresh the model's crop views, which means fetching rows
// into the MODEL_SLOPE_CROP and MODEL_SOIL_CROP materialized
// views. Please note that the double semicolon ";;" is used
// as a separator for eventually splitting the command line
// text into its constituent options and arguments.
var psql_crop_cmd = "psql;;-d;;%v;;-U;;postgres;;-c;;%v;;-c;;%v;;-1"

func refreshCropViews(db *Model) error {
  sqlrefresh_slope_crop := fmt.Sprintf(refresh_cmd, db.schema, MODEL_SLOPE_CROP)
  sqlrefresh_soil_crop := fmt.Sprintf(refresh_cmd, db.schema, MODEL_SOIL_CROP)
  proc := fmt.Sprintf(psql_crop_cmd, db.name, sqlrefresh_slope_crop, sqlrefresh_soil_crop)
  args := strings.Split(proc, ";;")
  if err := spawn.ProcExec(args); err != nil {
    return err
  }
  return nil
}

// psql_link_cmd var stores the psql command line text with placeholders
// needed to refresh the model's link view. Please note that the double
// semicolon ";;" is used as a separator for eventually splitting the
// command line text into its constituent options and arguments.
var psql_link_cmd = "psql;;-d;;%v;;-U;;postgres;;-c;;%v"

func refreshLinkView(db *Model) error {
  sqlrefresh_link := fmt.Sprintf(refresh_cmd, db.schema, MODEL_LINK)
  proc := fmt.Sprintf(psql_link_cmd, db.name, sqlrefresh_link)
  args := strings.Split(proc, ";;")
  if err := spawn.ProcExec(args); err != nil {
    return err
  }
  return nil
}

// psql_allocate_cmd var stores the psql command line text with placeholders
// needed to refresh the model's cover views, which means fetching
// rows into the MODEL_C_LEAF_NODES and MODEL_C_ALLOC materialized
// views. Please note that the double semicolon ";;" is used as a
// separator for eventually splitting the command line text into
// its constituent options and arguments.
var psql_allocate_cmd = "psql;;-d;;%v;;-U;;postgres;;-c;;%v;;-c;;%v;;-1"

func refreshCoverViews(db *Model) error {
  sqlrefresh_cover_leaf_nodes := fmt.Sprintf(refresh_cmd, db.schema, MODEL_C_LEAF_NODES)
  sqlrefresh_cover_allocation := fmt.Sprintf(refresh_cmd, db.schema, MODEL_C_ALLOC)
  proc := fmt.Sprintf(psql_allocate_cmd, db.name, sqlrefresh_cover_leaf_nodes, sqlrefresh_cover_allocation)
  args := strings.Split(proc, ";;")
  if err := spawn.ProcExec(args); err != nil {
    return err
  }
  return nil
}

// psql_wmean_cmd var stores the psql command line text with placeholders
// needed to refresh the n an p0 weighted means view. Please note that the double
// semicolon ";;" is used as a separator for eventually splitting the
// command line text into its constituent options and arguments.
var psql_wmean_cmd = "psql;;-d;;%v;;-U;;postgres;;-c;;%v"

func refreshWMeanView(db *Model) error {
  sqlrefresh_wmean := fmt.Sprintf(refresh_cmd, db.schema, MODEL_WMEAN)
  proc := fmt.Sprintf(psql_wmean_cmd, db.name, sqlrefresh_wmean)
  args := strings.Split(proc, ";;")
  if err := spawn.ProcExec(args); err != nil {
    return err
  }
  return nil
}

// psql_getp0_cmd var stores the psql command line text with placeholders
// needed to get the computed basin's p0 value. Please note that the double
// semicolon ";;" is used as a separator for eventually splitting the
// command line text into its constituent options and arguments.
var psql_getp0_cmd = "psql;;-d;;%v;;-U;;postgres;;-c;;%v;;-t"

func selectBasinP0(db *Model) (string, error) {
  sqlselect_p0 := fmt.Sprintf(select_scalar_cmd, MODEL_BASIN_P0, db.schema, MODEL_WMEAN)
  proc := fmt.Sprintf(psql_getp0_cmd, db.name, sqlselect_p0)
  args := strings.Split(proc, ";;")
  r, err := spawn.ProcGet(args)
  if err != nil {
    return "", err
  }
  return r, nil
}