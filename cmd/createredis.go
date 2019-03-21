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
	tcloudredis "github.com/rdxsl/tcloud-provisioner/tencent-cloud/redis"
	"github.com/spf13/cobra"
)

func createTCloudRedisFromConfig(redisConfName string) (*tcloudredis.TCloudRedis, error) {
	byteValue, err := ioutil.ReadFile(redisConfName)
	if err != nil {
		return nil, errors.Wrapf(err, "error reading %s file", redisConfName)
	}
	var tr tcloudredis.TCloudRedis
	if err = json.Unmarshal(byteValue, &tr); err != nil {
		return nil, errors.Wrapf(err, "error unmarshaling %s into TCloudRedis struct", redisConfName)
	}
	return &tr, err
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
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tr, err := createTCloudRedisFromConfig(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

		exists, err := tr.Create()
		if err != nil {
			fmt.Println("[ERROR] Failed to create Redis instance:", err)
			os.Exit(1)
		}
		if exists {
			fmt.Printf("[INFO] Redis instance %#v already exists\n", tr.InstanceName)
			return
		}
		fmt.Printf("[INFO] Redis instance %#v successfully created\n", tr.InstanceName)
	},
}

func init() {
	redisCmd.AddCommand(createRedisCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
