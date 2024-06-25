package logging

import (
	"fmt"
	"testing"
)

func TestGetHostWorkload(t *testing.T) {
	result := GetHostWorkload()
	fmt.Println(result)

}
