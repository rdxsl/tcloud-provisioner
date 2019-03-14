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

var mysqlConfigPath string

func mysqlViperParseConfig() (*tcloudmysql.TcloudMySQL, error) {

	bytes, err := ioutil.ReadFile(mysqlConfigPath)
	if err != nil {
		return nil, errors.Wrapf(err, "error reading %s file", mysqlConfigPath)
	}

	var tm tcloudmysql.TcloudMySQL
	err = json.Unmarshal(bytes, &tm)
	if err != nil {
		return nil, errors.Wrapf(err, "error marshaling %s into TcloudMySQL struct", mysqlConfigPath)
	}
	return &tm, nil
}

// createCmd represents the create command
var createMysqlCmd = &cobra.Command{
	Use:   "create",
	Short: "create a MySQL DB in Tencent Cloud",
	Long: `create a MySQL DB in Tencent Cloud. The configuration is set in the following dir
	~/.tcloud-provisioner
	`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("reading tcloud mysql config from %#v\n", mysqlConfigPath)

		tm, err := mysqlViperParseConfig()
		if err != nil {
			fmt.Println("Fatal error config file:", err)
			os.Exit(1)
		}

		exists, err := tm.Create()
		if err != nil {
			fmt.Println("[ERROR] Failed to create MySQL:", err)
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

	createMysqlCmd.Flags().StringVarP(&mysqlConfigPath, "env", "e", "", "path to mysql config file")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
