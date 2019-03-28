package redis

import (
	"github.com/k0kubun/pp"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"

	common "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
)

// List lists Redis instances
func List(region string) error {

	// Init client
	credential, cpf := NewCredential()
	client, err := redis.NewClient(credential, region, cpf)
	if err != nil {
		return err
	}

	// List instances
	req := redis.NewDescribeInstancesRequest()
	req.Limit = common.Uint64Ptr(100)
	resp, err := client.DescribeInstances(req)
	if err != nil {
		return err
	}
	pp.Println(resp)
	return nil
}
