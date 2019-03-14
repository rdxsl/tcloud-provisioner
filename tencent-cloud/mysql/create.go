package mysql

import (
	"encoding/json"
	"fmt"
	"os"

	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	common "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

type TcloudMySQL struct {
	InstanceName string
	Region       string
	Zone         string
	Instance     int64
	Memory       int64
	Volume       int64
}

func NewCredential() (*common.Credential, *profile.ClientProfile) {
	credential := common.NewCredential(
		os.Getenv("TENCENTCLOUD_SECRET_ID"),
		os.Getenv("TENCENTCLOUD_SECRET_KEY"),
	)

	cpf := profile.NewClientProfile()
	cpf.HttpProfile.ReqMethod = "POST"
	cpf.HttpProfile.ReqTimeout = 60
	cpf.SignMethod = "HmacSHA1"

	return credential, cpf
}

func (tm TcloudMySQL) Create() {

	// Init client
	credential, cpf := NewCredential()
	client, _ := cdb.NewClient(credential, tm.Region, cpf)

	exists, err := checkInstanceExists(client, tm.InstanceName)
	if err != nil {
		fmt.Println("Failed to check if instances already exist:", err)
		os.Exit(1)
	}
	if exists {
		fmt.Printf("[INFO] %#v already exists\n", tm.InstanceName)
		os.Exit(0)
	}

	request := cdb.NewCreateDBInstanceHourRequest()
	request.GoodsNum = common.Int64Ptr(tm.Instance)
	request.Memory = common.Int64Ptr(tm.Memory)
	request.Volume = common.Int64Ptr(tm.Volume)
	request.Zone = common.StringPtr(tm.Zone)
	response, err := client.CreateDBInstanceHour(request)

	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return
	}
	// unexpected errors
	if err != nil {
		panic(err)
	}
	b, _ := json.Marshal(response.Response)
	fmt.Printf("%s", b)

	fmt.Printf("[INFO] %#v successfully created\n", tm.InstanceName)
	os.Exit(0)
}

func checkInstanceExists(client *cdb.Client, name string) (bool, error) {
	req := cdb.NewDescribeDBInstancesRequest()
	resp, err := client.DescribeDBInstances(req)
	if err != nil {
		return false, err
	}
	for _, item := range resp.Response.Items {
		if item.InstanceName == nil {
			continue
		}
		if *item.InstanceName == name {
			return true, nil
		}
	}
	return false, nil
}
