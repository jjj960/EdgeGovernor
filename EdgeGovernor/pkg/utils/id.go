package utils

import (
	"github.com/bwmarrin/snowflake"
)

var SnowFlake *snowflake.Node

func GetID() (int64, int64) {
	a := SnowFlake.Generate()
	return a.Int64(), a.Time()
}
