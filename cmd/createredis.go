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

	tcloudredis "github.com/rdxsl/tcloud-provisioner/tencent-cloud/redis"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var redisConfName string

func redisParseConfig(tm *tcloudredis.TcloudRedis) error {
	switch viper.Get("Instance").(type) {
	case float64:
		tm.Instance = int64(viper.Get("Instance").(float64))
	default:
		return fmt.Errorf("mysql config needs to have 'Instance' with type 'int' without quotes, currently we have %+v\n", viper.Get("Instance"))
	}
	switch viper.Get("Memory").(type) {
	case float64:
		tm.Memory = int64(viper.Get("Memory").(float64))
	default:
		return fmt.Errorf("mysql config needs to have 'Memory' with type 'int' without quotes, currently we have %+v\n", viper.Get("Memory"))
	}
	switch viper.Get("Volume").(type) {
	case float64:
		tm.Volume = int64(viper.Get("Volume").(float64))
	default:
		return fmt.Errorf("mysql config needs to have 'Volume' with type 'int' without quotes, currently we have %+v\n", viper.Get("Volume"))
	}
	switch viper.Get("Region").(type) {
	case string:
		tm.Region = string(viper.Get("Region").(string))
	default:
		return fmt.Errorf("mysql config needs to have 'Region' with type 'string' with quotes, currently we have %+v\n", viper.Get("Region"))
	}
	switch viper.Get("Zone").(type) {
	case string:
		tm.Zone = string(viper.Get("Zone").(string))
	default:
		return fmt.Errorf("mysql config needs to have 'Zone' with type 'string' with quotes, currently we have %+v\n", viper.Get("Zone"))
	}
	return nil
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
		viper.SetConfigFile(redisConfName)
		err := viper.ReadInConfig()
		if err != nil {
			fmt.Println(fmt.Errorf("Fatal error config file: %s \n", err))
		}
		var tr tcloudredis.TcloudRedis
		err = redisParseConfig(&tr)
		tr.Create()

		// if err == nil {
		// 	tm.Create()
		// } else {
		// 	fmt.Println(err)
		// }
		// fmt.Println(viper.Get("Instance"))
		// tcloudredis.TcloudRedis()
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
