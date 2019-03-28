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
	"fmt"
	"os"

	tcloudredis "github.com/rdxsl/tcloud-provisioner/tencent-cloud/redis"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listRedisCmd = &cobra.Command{
	Use:   "list",
	Short: "List Redis instances in Tencent Cloud",
	Long:  "List Redis instances in Tencent Cloud",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		region := "na-siliconvalley"
		if len(args) > 1 {
			region = args[1]
		}
		if err := tcloudredis.List(region); err != nil {
			fmt.Println("[ERROR]", err)
			os.Exit(1)
		}
	},
}

func init() {
	redisCmd.AddCommand(listRedisCmd)
}
