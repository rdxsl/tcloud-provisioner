// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	tcloudmysql "github.com/rdxsl/tcloud-provisioner/tencent-cloud/mysql"
	"github.com/spf13/cobra"
)

func createTCloudMySQLFromConfig(mysqlConfigPath string) (*tcloudmysql.TCloudMySQL, error) {
	byteValue, err := ioutil.ReadFile(mysqlConfigPath)
	if err != nil {
		return nil, errors.Wrapf(err, "error reading %s file", mysqlConfigPath)
	}
	var tm tcloudmysql.TCloudMySQL
	if err = json.Unmarshal(byteValue, &tm); err != nil {
		return nil, errors.Wrapf(err, "error unmarshaling %s into TCloudMySQL struct", mysqlConfigPath)
	}
	if err = tm.Validate(); err != nil {
		return nil, errors.Wrap(err, "config invalid")
	}
	return &tm, nil
}

// createCmd represents the create command
var createMysqlCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a MySQL DB instance in Tencent Cloud",
	Long:  "Create a MySQL DB instance in Tencent Cloud",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tm, err := createTCloudMySQLFromConfig(args[0])
		if err != nil {
			fmt.Println("Fatal error config file:", err)
			os.Exit(1)
		}

		exists, err := tm.Create()
		if err != nil {
			fmt.Println("[ERROR] Failed to create MySQL instance:", err)
			os.Exit(1)
		}
		if exists {
			fmt.Printf("[INFO] MySQL instance %#v already exists\n", tm.InstanceName)
			return
		}
		fmt.Printf("[INFO] MySQL instance %#v successfully created\n", tm.InstanceName)
	},
}

func init() {
	mysqlCmd.AddCommand(createMysqlCmd)
}
