package duckdb

var (
	createHostloadSQL = `CREATE TABLE IF NOT EXISTS host_workload (
				Timestamp           BIGINT(20),
				Hostname            TEXT,
				CPUUsagePercent     REAL,
				CPUCapacity         REAL,
				CPUResidue          REAL,
				MemoryUsedPercent   REAL,
				MemoryCapacity      REAL,
				MemoryResidue       REAL,
				DiskUsedPercent     REAL,
				DiskCapacity        REAL,
				DiskResidue         REAL,
				BytesReceived       REAL,
				BytesSent           REAL
			);`
)

var (
	insertOperationLogSQL = `
	INSERT INTO operation_log (
		operator, operation_time, operation_type, operation_detail,
		operation_result, ip_address, remarks
	) VALUES (?, ?, ?, ?, ?, ?, ?);
	`

	createOperationLogTableSQL = `
		CREATE TABLE IF NOT EXISTS operation_log (
			id INTEGER PRIMARY KEY AUTO_INCREMENT,
			operator VARCHAR(255),
			operation_time DATETIME DEFAULT CURRENT_TIMESTAMP,
			operation_type VARCHAR(50),
			operation_detail TEXT,
			operation_result VARCHAR(50),
			ip_address VARCHAR(50),
			remarks TEXT
		);`

	selectOperationLogSQL = `
	SELECT
		id,
		operator,
		operation_time,
		operation_type,
		operation_detail,
		operation_result,
		ip_address,
		remarks
	FROM
		operation_log;
	`
)
