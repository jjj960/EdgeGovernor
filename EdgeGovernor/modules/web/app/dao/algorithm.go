package dao

import "EdgeGovernor/pkg/cache/algorithm"

func GetAlgorithmMsg(alName string) algorithm.AlgorithmStatus {
	msg := algorithm.GetAlgorithmMsg(alName)
	return msg
}

func GetAlgorithmNames() []string {
	algorithmName := algorithm.GetAlgorithmNames()
	return algorithmName
}

func GetAlgorithms() []map[string]string {
	return algorithm.GetAlgorithmsMsg()
}

func AddAlgorithms(name, url, use, detail string) {
	algorithm.AddAlgorithmStatus(name, url, detail, use)
}
