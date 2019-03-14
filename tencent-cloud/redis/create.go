package redis

import (
	"encoding/json"
	"fmt"
	"os"

	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"

	common "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

type TcloudRedis struct {
	Region      string `json:"region"`
	Instance    uint64 `json:"instance"`
	ZoneId      uint64 `json:"zoneid"`
	TypeId      uint64 `json:"typeid"`
	MemSize     uint64 `json:"memsize"`
	Password    string `json:"password"`
	ProjectId   int64  `json:"projectid"`
	Vport       uint64 `json:"vport"`
	BillingMode int64  `json:"billingmode"`
	Period      uint64 `json:"billperiod"`
	VpcId       string `json:"vpcid"`
	SubnetId    string `json:"subnetid"`
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

func (tr TcloudRedis) Create() {
	credential, cpf := NewCredential()
	client, _ := redis.NewClient(credential, tr.Region, cpf)
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

	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s\n", err)
		return
	}
	// unexpected errors
	if err != nil {
		panic(err)
	}
	b, _ := json.Marshal(response.Response)
	fmt.Printf("%s\n", b)
}
