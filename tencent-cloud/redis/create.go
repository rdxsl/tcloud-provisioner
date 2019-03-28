package redis

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"

	"github.com/pkg/errors"
	common "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	tencentErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

type TCloudRedis struct {
	InstanceName string `json:"instancename"`
	Region       string `json:"region"`
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

func (tr TCloudRedis) Validate() error {
	if tr.InstanceName == "" {
		return errors.New("instancename cannot be empty")
	}
	if len(tr.InstanceName) > 36 {
		return errors.New("instancename is too long, expected 1-36 chars https://cloud.tencent.com/document/api/239/8431")
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

func (tr TCloudRedis) Create() (exists bool, err error) {

	// Init client
	credential, cpf := NewCredential()
	client, err := redis.NewClient(credential, tr.Region, cpf)
	if err != nil {
		return false, errors.Wrap(err, "new client error")
	}

	// NOTE(jeff): This is currently not working becuase the Tencent Cloud SDK
	// does not pipe InstanceName to the API. Re-enable once that is hooked up
	//
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

	// Create instance
	resp, err := client.CreateInstances(request)
	if err != nil {
		if _, ok := err.(*tencentErrors.TencentCloudSDKError); ok {
			return false, errors.Wrap(err, "API error")
		}
		return false, err
	}

	// Parse response
	b, err := json.Marshal(resp.Response)
	if err != nil {
		// Want to return nil right here, because the instance has already
		// been successfully created
		fmt.Println(err)
		return false, nil
	}
	fmt.Printf("[INFO] Created API Response: %v\n", string(b))

	// // Find Instances Created
	// ids, err := findRedisIDByDealID(client, *resp.Response.DealId)
	// if err != nil {
	// 	fmt.Println("[ERROR] Failed to find instance by deal id:", *resp.Response.DealId)
	// 	return false, nil
	// }

	// NOTE(jehwang): Heuristically find instances created
	instance, err := findRecentlyCreatedRedisInstance(client)
	if err != nil {
		fmt.Println("[ERROR] Failed to update Redis instance name. Could not find recently created instance:", err)
		return false, nil
	}
	id := *instance.InstanceId

	// Update name
	if err := updateInstanceName(client, id, tr.InstanceName); err != nil {
		fmt.Println("[ERROR] Failed to update Redis instance name:", id, tr.InstanceName, err)
	}
	return false, nil
}

func checkInstanceExists(client *redis.Client, name string) (bool, error) {

	// https://intl.cloud.tencent.com/document/api/239/1384
	req := redis.NewDescribeInstancesRequest()
	req.Limit = common.Uint64Ptr(100)
	req.InstanceName = &name
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

// func findRedisIDByDealID(client *redis.Client, dealID string) (ids []string, err error) {
// 	req := redis.NewDescribeInstanceDealDetailRequest()
// 	req.DealIds = append(req.DealIds, &dealID)
// 	resp, err := client.DescribeInstanceDealDetail(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	for _, details := range resp.Response.DealDetails {
// 		for _, instanceID := range details.InstanceIds {
// 			ids = append(ids, *instanceID)
// 		}
// 	}
// 	return ids, nil
// }

// findRecentlyCreatedRedisInstance looks for an instance that was recently
// created using the criteria:
//	- InstanceName == InstanceId
//	- Createtime +/- 5min, or Createtime is zero
//
// Also note Createtime is given back in Asia/Shanghai timezone.
func findRecentlyCreatedRedisInstance(client *redis.Client) (*redis.InstanceSet, error) {

	// https://intl.cloud.tencent.com/document/api/239/1384
	req := redis.NewDescribeInstancesRequest()
	req.Limit = common.Uint64Ptr(100)
	resp, err := client.DescribeInstances(req)
	if err != nil {
		return nil, err
	}

	for _, item := range resp.Response.InstanceSet {
		if item.InstanceId == nil || item.InstanceName == nil {
			continue
		}

		// Look for instances that have id set to name (not tagged with name yet)
		if *item.InstanceName != *item.InstanceId {
			continue
		}

		// And recently created within the last 5 minutes or created time set to zero
		var isRecentlyCreated bool
		createdRaw := *item.Createtime
		if createdRaw == "0000-00-00 00:00:00" {
			isRecentlyCreated = true
			fmt.Printf("[INFO] Found Recently Created Instance with zero time: %v %v\n", *item.InstanceId, createdRaw)
		} else {
			shanghaiLoc, _ := time.LoadLocation("Asia/Shanghai")
			created, err := time.ParseInLocation("2006-01-02 15:04:05", createdRaw, shanghaiLoc)
			if err != nil {
				return nil, err
			}
			dur := time.Since(created)
			if dur > -5*time.Minute && dur < 5*time.Minute {
				fmt.Printf("[INFO] Found Recently Created Instance < 5min: %v %v\n", *item.InstanceId, dur)
				isRecentlyCreated = true
			}
		}
		if isRecentlyCreated {
			return item, nil
		}
	}
	return nil, errors.New("no matching criteria")
}

func updateInstanceName(client *redis.Client, id, name string) error {
	// https://cloud.tencent.com/document/api/239/8431
	req := redis.NewModifyInstanceRequest()
	operation := "rename"
	req.Operation = &operation
	req.InstanceId = &id
	req.InstanceName = &name
	_, err := client.ModifyInstance(req)
	if err != nil {
		return err
	}
	fmt.Printf("[INFO] Updated Redis instance: %v with name %v\n", id, name)
	return nil
}
