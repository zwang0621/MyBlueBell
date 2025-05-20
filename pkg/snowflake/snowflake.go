package snowflake

import (
	"time"

	"github.com/bwmarrin/snowflake"
)

// 基于雪花算法生成分布式id
var node *snowflake.Node

func Init(startTime string, machineID int64) (err error) {
	var st time.Time
	st, err = time.Parse("2006-01-02", startTime)
	if err != nil {
		return
	}
	snowflake.Epoch = st.UnixNano() / 1000000
	node, err = snowflake.NewNode(machineID)
	return
}

func GenID() int64 {
	return node.Generate().Int64()
}

// func main() {
// 	if err := Init("2025-04-16", 1); err != nil {
// 		fmt.Printf("init failed, err:%v\n", err)
// 		return
// 	}
// 	id := GenID()
// 	fmt.Println(id)
// }
