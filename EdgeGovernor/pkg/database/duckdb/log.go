package duckdb

import (
	"EdgeGovernor/pkg/models"
	"EdgeGovernor/pkg/utils"
	"fmt"
	"log"
)

//logEntry := models.OperationLog{
//ID:            a,
//NodeName:      "Node1",
//NodeIP:        "192.168.1.1",
//OperationType: "Update",
//Description:   "Update system configuration",
//Result:        true,
//CreatedAt:     time.Now(),
//}

// createTableOperationLog 函数用于创建OperationLog表
func CreateTableOperationLog() {
	_, err := utils.DuckDBCli.Exec(`
    CREATE TABLE IF NOT EXISTS OperationLog (
        ID BIGINT,
        NodeName VARCHAR(255),
        NodeIP VARCHAR(255),
        OperationType VARCHAR(255),
        Description TEXT,
        Result BOOLEAN,
        CreatedAt TIMESTAMP
    );
    `)
	if err != nil {
		log.Fatal(err)
	}
}

// insertOperationLog 函数用于向OperationLog表插入数据
func InsertOperationLog(logEntry models.OperationLog) {
	stmt, err := utils.DuckDBCli.Prepare("INSERT INTO OperationLog (ID, NodeName, NodeIP, OperationType, Description, Result, CreatedAt) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(logEntry.ID, logEntry.NodeName, logEntry.NodeIP, logEntry.OperationType, logEntry.Description, logEntry.Result, logEntry.CreatedAt)
	if err != nil {
		log.Fatal(err)
	}
}

// dropTableOperationLog 函数用于删除OperationLog表
func DropTableOperationLog() {
	_, err := utils.DuckDBCli.Exec("DROP TABLE OperationLog")
	if err != nil {
		log.Fatal(err)
	}
}

// selectAllOperationLogs 函数用于选择并打印OperationLog表中的所有记录
func SelectAllOperationLogs() {
	rows, err := utils.DuckDBCli.Query("SELECT ID, NodeName, NodeIP, OperationType, Description, Result, CreatedAt FROM OperationLog")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var logEntry models.OperationLog
		err := rows.Scan(&logEntry.ID, &logEntry.NodeName, &logEntry.NodeIP, &logEntry.OperationType, &logEntry.Description, &logEntry.Result, &logEntry.CreatedAt)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID: %d, NodeName: %s, NodeIP: %s, OperationType: %s, Description: %s, Result: %t, CreatedAt: %s\n",
			logEntry.ID, logEntry.NodeName, logEntry.NodeIP, logEntry.OperationType, logEntry.Description, logEntry.Result, logEntry.CreatedAt.Format("2006-01-02 15:04:05"))
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
}
