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
	Region   string
	Zone     string
	Instance int64
	Memory   int64
	Volume   int64
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
	request.ZoneId = common.Uint64Ptr(150001)
	//Instance Type: 2 for Redis2.8 Master-Salve Edition
	//Other type currently not available on Intl console
	request.TypeId = common.Uint64Ptr(2)
	//Instance memory size(unit: MB)
	request.MemSize = common.Uint64Ptr(1024)
	//Instance quantity you want to purchase
	request.GoodsNum = common.Uint64Ptr(1)
	//Instance password
	request.Password = common.StringPtr("test12341234")
	//put 0, Redis instance goes to default project
	//put certain project ID, goes to certain project
	request.ProjectId = common.Int64Ptr(0)
	//set the port (default port 6379)
	request.VPort = common.Uint64Ptr(6379)
	//You don't need to change BillingMode and Period for intl site
	request.BillingMode = common.Int64Ptr(0)
	request.Period = common.Uint64Ptr(1)

	// // VPC's ID
	request.VpcId = common.StringPtr("vpc-gfehxte9")
	// // VPC's subnet ID
	request.SubnetId = common.StringPtr("subnet-fiyvufly")
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
