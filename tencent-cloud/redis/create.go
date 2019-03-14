package redis

import (
	"encoding/json"
	"fmt"
	"os"

	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"

	"github.com/pkg/errors"
	common "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	tencentErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

type TCloudRedis struct {
	Region       string `json:"region"`
	InstanceName string `json:"instanceName"`
	Instance     uint64 `json:"instance"`
	ZoneId       uint64 `json:"zoneid"`
	TypeId       uint64 `json:"typeid"`
	MemSize      uint64 `json:"memsize"`
	Password     string `json:"password"`
	ProjectId    int64  `json:"projectid"`
	Vport        uint64 `json:"vport"`
	BillingMode  int64  `json:"billingmode"`
	Period       uint64 `json:"billperiod"`
	VpcId        string `json:"vpcid"`
	SubnetId     string `json:"subnetid"`
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

func (tr TCloudRedis) Create() (exists bool, err error) {

	// Init client
	credential, cpf := NewCredential()
	client, err := redis.NewClient(credential, tr.Region, cpf)
	if err != nil {
		return false, errors.Wrap(err, "new client error")
	}

	// Check if instance with the same name already exists
	exists, err = checkInstanceExists(client, tr.InstanceName)
	if err != nil || exists {
		return exists, errors.Wrap(err, "failed to check instances")
	}

	// Create redis
	request := redis.NewCreateInstancesRequest()
	//Availability Zone ID
	//Please see the comment on the bottom for different Availability Zone ID
	request.ZoneId = common.Uint64Ptr(tr.ZoneId)
	//Instance Type: 2 for Redis2.8 Master-Salve Edition
	//Other type currently not available on Intl console
	request.TypeId = common.Uint64Ptr(tr.TypeId)
	//Instance memory size(unit: MB)
	request.MemSize = common.Uint64Ptr(tr.MemSize)
	//Instance quantity you want to purchase
	request.GoodsNum = common.Uint64Ptr(tr.Instance)
	//Instance password
	request.Password = common.StringPtr(tr.Password)
	//put 0, Redis instance goes to default project
	//put certain project ID, goes to certain project
	request.ProjectId = common.Int64Ptr(tr.ProjectId)
	//set the port (default port 6379)
	request.VPort = common.Uint64Ptr(tr.Vport)
	//You don't need to change BillingMode and Period for intl site
	request.BillingMode = common.Int64Ptr(tr.BillingMode)
	request.Period = common.Uint64Ptr(tr.Period)

	// // VPC's ID
	request.VpcId = common.StringPtr(tr.VpcId)
	// // VPC's subnet ID
	request.SubnetId = common.StringPtr(tr.SubnetId)
	// security groups' ID
	// var SGroup = []string {"sg-clygup2w"}
	// request.SecurityGroupIdList = common.StringPtrs(SGroup)

	response, err := client.CreateInstances(request)
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
	fmt.Printf("%v", string(b))
	return false, nil
}

func checkInstanceExists(client *redis.Client, name string) (bool, error) {
	// https://intl.cloud.tencent.com/document/api/239/1384
	req := redis.NewDescribeInstancesRequest()
	req.Limit = common.Uint64Ptr(100)
	resp, err := client.DescribeInstances(req)
	if err != nil {
		return false, err
	}

	for _, item := range resp.Response.InstanceSet {
		if item.InstanceName == nil {
			continue
		}
		if *item.InstanceName == name {
			return true, nil
		}
	}
	return false, nil
}
