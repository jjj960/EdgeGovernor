package resource

import (
	"EdgeGovernor/pkg/utils"
	"context"
	"fmt"
	"github.com/buger/jsonparser"
	clientv3 "go.etcd.io/etcd/client/v3"
	"testing"
)

func Test(t *testing.T) {
	utils.GetETCDCli()
	prefix := "/menet/node"

	resp, _ := utils.ETCDCli.Get(context.Background(), prefix, clientv3.WithPrefix())
	for _, kv := range resp.Kvs {
		re, aa, bb, _ := jsonparser.Get(kv.Value, "Hostname")
		fmt.Println(string(re))
		fmt.Println(aa)
		fmt.Println(bb)
	}
}
