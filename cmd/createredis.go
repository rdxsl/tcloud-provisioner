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

	tcloudredis "github.com/rdxsl/tcloud-provisioner/tencent-cloud/redis"
	"github.com/spf13/cobra"
)

var redisConfName string

func redisParseConfig(tm *tcloudredis.TcloudRedis) error {
	// Open our jsonFile
	jsonFile, err := os.Open(redisConfName)
	if err != nil {
		return err
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}
	return json.Unmarshal(byteValue, tm)
}

// createCmd represents the create command
var createRedisCmd = &cobra.Command{
	Use:   "create",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var tr tcloudredis.TcloudRedis
		err := redisParseConfig(&tr)

		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

		err = tr.Create()

		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			// we might want to define exit code later
			os.Exit(2)
		}
	},
}

func init() {
	redisCmd.AddCommand(createRedisCmd)
	// createRedisCmd.Flags().StringVarP(&redisConfName, "conf", "c", "", "location of the redis conf json file")

	createRedisCmd.Flags().StringVarP(&redisConfName, "env", "e", "", "sub directory inside of ./conf what has a region's config files")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
