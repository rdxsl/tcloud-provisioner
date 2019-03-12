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

	tcloudmysql "github.com/rdxsl/tcloud-provisioner/tencent-cloud/mysql"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var tcloudEnvName string

// createCmd represents the create command
var createMysqlCmd = &cobra.Command{
	Use:   "create",
	Short: "create a MySQL DB in Tencent Cloud",
	Long: `create a MySQL DB in Tencent Cloud. The configuration is set in the following dir
	~/.tcloud-provisioner
	`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("reading tcloud mysql config from ./conf/" + tcloudEnvName)
		viper.SetConfigName("mysql")
		viper.AddConfigPath("./conf/" + tcloudEnvName)
		err := viper.ReadInConfig()
		if err != nil { // Handle errors reading the config file
			fmt.Println(fmt.Errorf("Fatal error config file: %s \n", err))
		}
		var tm tcloudmysql.TcloudMySQL
		tm.Instance = int64(viper.Get("Instance").(float64))
		tm.Memory = int64(viper.Get("Memory").(float64))
		tm.Volume = int64(viper.Get("Volume").(float64))
		tm.Region = string(viper.Get("Region").(string))
		tm.Zone = string(viper.Get("Zone").(string))
		tm.Create()
	},
}

func init() {
	mysqlCmd.AddCommand(createMysqlCmd)
	createMysqlCmd.Flags().StringVarP(&tcloudEnvName, "env", "e", "", "sub directory inside of ./conf what has a region's config files")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
