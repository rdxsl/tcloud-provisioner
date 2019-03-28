package mysql

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pkg/errors"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	common "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	tencentErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

type TCloudMySQL struct {
	Region        string `json:"region"`
	Zone          string `json:"zone"`
	Instance      int64  `json:"instance"`
	Memory        int64  `json:"memory"`
	Volume        int64  `json:"volume"`
	VpcId         string `json:"vpcid"`
	SubnetId      string `json:"subnetid"`
	Password      string `json:"password"`
	InstanceName  string `json:"instancename"`
	EngineVersion string `json:"engineversion"`
	ProjectId     int64  `json:"projectid"`
}

func (tm TCloudMySQL) Validate() error {
	if tm.InstanceName == "" {
		return errors.New("instancename cannot be empty")
	}
	return nil
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

func (tm TCloudMySQL) Create() (bool, error) {

	// Init client
	credential, cpf := NewCredential()
	client, err := cdb.NewClient(credential, tm.Region, cpf)
	if err != nil {
		return false, err
	}

	// Check if instance with the same name already exists
	exists, err := checkInstanceExists(client, tm.InstanceName)
	if err != nil || exists {
		return exists, errors.Wrap(err, "failed to check instances")
	}

	// Create instance
	request := cdb.NewCreateDBInstanceHourRequest()
	request.GoodsNum = common.Int64Ptr(tm.Instance)
	request.Memory = common.Int64Ptr(tm.Memory)
	request.Volume = common.Int64Ptr(tm.Volume)
	request.Zone = common.StringPtr(tm.Zone)
	request.UniqVpcId = common.StringPtr(tm.VpcId)
	request.UniqSubnetId = common.StringPtr(tm.SubnetId)
	request.Password = common.StringPtr(tm.Password)
	request.InstanceName = common.StringPtr(tm.InstanceName)
	request.EngineVersion = common.StringPtr(tm.EngineVersion)
	request.ProjectId = common.Int64Ptr(tm.ProjectId)
	response, err := client.CreateDBInstanceHour(request)
	if _, ok := err.(*tencentErrors.TencentCloudSDKError); ok {
		return false, errors.Wrap(err, "API error")
	}
	if err != nil {
		return false, err
	}

	// Parse response
	b, err := json.Marshal(response.Response)
	if err != nil {
		// Want to return nil right here, because the instance has already
		// been successfully created
		fmt.Println(err)
		return false, nil
	}
	fmt.Printf("[INFO] API Response: %v\n", string(b))
	return false, nil
}

func checkInstanceExists(client *cdb.Client, name string) (bool, error) {
	// https://intl.cloud.tencent.com/document/api/236/1266
	req := cdb.NewDescribeDBInstancesRequest()
	req.Limit = common.Uint64Ptr(100)
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
