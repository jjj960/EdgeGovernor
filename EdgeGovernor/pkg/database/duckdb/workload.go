package duckdb

import (
	"EdgeGovernor/pkg/models"
	"EdgeGovernor/pkg/utils"
	"fmt"
	"log"
)

// createHostloadTable 函数用于创建Hostload表
func CreateHostloadTable() {
	// 定义创建表的SQL语句
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS Hostload (
        Timestamp BIGINT,
        Hostname VARCHAR(255),
        CPUUsagePercent FLOAT,
        CPUCapacity BIGINT,
        CPUResidue BIGINT,
        MemoryUsedPercent FLOAT,
        MemoryCapacity BIGINT,
        MemoryResidue BIGINT,
        DiskUsedPercent FLOAT,
        DiskCapacity BIGINT,
        DiskResidue BIGINT,
        BytesRecv FLOAT,
        BytesSent FLOAT,
        BandWidth FLOAT
    );
    `

	// 执行SQL语句
	_, err := utils.DuckDBCli.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
	fmt.Println("Hostload table created successfully.")
}

// InsertHostload 函数用于插入数据到Hostload表
func InsertHostload(hostload models.Hostload) error {
	stmt, err := utils.DuckDBCli.Prepare("INSERT INTO Hostload (Timestamp,Hostname,CPUUsagePercent,CPUCapacity,CPUResidue,MemoryUsedPercent,MemoryCapacity,MemoryResidue,DiskUsedPercent,DiskCapacity,DiskResidue,BytesRecv,BytesSent,BandWidth) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Println(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(hostload.Timestamp,
		hostload.Hostname,
		hostload.CPUUsagePercent,
		hostload.CPUCapacity,
		hostload.CPUResidue,
		hostload.MemoryUsedPercent,
		hostload.MemoryCapacity,
		hostload.MemoryResidue,
		hostload.DiskUsedPercent,
		hostload.DiskCapacity,
		hostload.DiskResidue,
		hostload.BytesRecv,
		hostload.BytesSent,
		hostload.BandWidth)
	if err != nil {
		log.Println(err)
	}

	return nil
}

// SelectAllHostloads 函数用于查询并输出Hostload表中的所有数据
func SelectAllHostloads() error {
	// 准备SQL查询语句
	query := "SELECT * FROM Hostload"

	// 执行查询
	rows, err := utils.DuckDBCli.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	// 遍历查询结果
	for rows.Next() {
		var hostload models.Hostload
		// 使用Scan填充Hostload结构体
		err := rows.Scan(
			&hostload.Timestamp,
			&hostload.Hostname,
			&hostload.CPUUsagePercent,
			&hostload.CPUCapacity,
			&hostload.CPUResidue,
			&hostload.MemoryUsedPercent,
			&hostload.MemoryCapacity,
			&hostload.MemoryResidue,
			&hostload.DiskUsedPercent,
			&hostload.DiskCapacity,
			&hostload.DiskResidue,
			&hostload.BytesRecv,
			&hostload.BytesSent,
			&hostload.BandWidth,
		)
		if err != nil {
			return err
		}

		// 打印Hostload实例的详细信息
		fmt.Printf("%+v\n", hostload)
	}

	// 检查是否有可能的错误发生在迭代结束后
	if err = rows.Err(); err != nil {
		return err
	}

	return nil
}

// createContainerloadTable 函数用于创建ContainerLoad表
func CreateContainerloadTable(containerName string) {
	// 定义创建表的SQL语句
	createTableSQL := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s_load (
		Timestamp BIGINT,
		Name VARCHAR(255),
		ID VARCHAR(255),
		CPUPercentage DOUBLE PRECISION,
		Memory BIGINT,
		MemoryLimit BIGINT,
		MemoryPercentage DOUBLE PRECISION,
		NetworkRx DOUBLE PRECISION,
		NetworkTx DOUBLE PRECISION,
		BlockRead DOUBLE PRECISION,
		BlockWrite DOUBLE PRECISION,
		PidsCurrent BIGINT
	);
	`, containerName)

	// 执行SQL语句
	_, err := utils.DuckDBCli.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
	fmt.Println("Hostload table created successfully.")
}

// InsertContainerLoad 函数用于插入数据到指定的动态表名中
func InsertContainerload(ContainerNamePrefix string, load models.ContainerLoad) error {
	tableName := fmt.Sprintf("%s_load", ContainerNamePrefix)
	stmt, err := utils.DuckDBCli.Prepare(fmt.Sprintf(`
	INSERT INTO %s (Timestamp, Name, ID, CPUPercentage, Memory, MemoryLimit, MemoryPercentage, NetworkRx, NetworkTx, BlockRead, BlockWrite, PidsCurrent)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, tableName))
	if err != nil {
		log.Println(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(load.Timestamp,
		load.Name,
		load.ID,
		load.CPUPercentage,
		load.Memory,
		load.MemoryLimit,
		load.MemoryPercentage,
		load.NetworkRx,
		load.NetworkTx,
		load.BlockRead,
		load.BlockWrite,
		load.PidsCurrent)
	if err != nil {
		log.Println(err)
	}

	return nil
}
