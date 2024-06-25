package algorithm

import (
	"EdgeGovernor/pkg/sec"
	"EdgeGovernor/pkg/utils"
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	utils.GetETCDCli()
	sec.Safer = sec.NewMSecLayer("2TXQfxBE#TrULn6FZx")
	AddAlgorithmStatus("TOPSIS", "http://192.168.47.128:50052/scheduler", "Schedule", "active")
	fmt.Println(GetAlgorithmURL("TOPSIS", "Schedule"))
}
