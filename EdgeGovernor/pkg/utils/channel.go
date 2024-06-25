package utils

import "EdgeGovernor/pkg/models"

var ModuleControlChannel chan bool //模块控制channel

var TaskMonitorChannel chan bool
var PopTaskChannel chan bool

var AlarmMsgChannel chan models.Msg

var JobOperationChannel chan models.JobOperation
