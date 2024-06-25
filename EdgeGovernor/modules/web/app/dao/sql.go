package dao

var (
	getNodesSQL        = `SELECT Hostname FROM nodes;`
	searchNodeMsg      = `SELECT IP, Role, Status FROM nodes WHERE Hostname = ?;`
	getNodesMsg        = `SELECT Hostname, IP, Role, Status FROM nodes LIMIT ?, ?;`
	getNodesTaskSQL    = `SELECT Name, PublishTime, Image, RequestCPU, RequestMem, RequestDisk, Status FROM task_info where DeployNode = ?;`
	getFollowerNodeSQL = `SELECT Hostname, IP, Status FROM nodes WHERE Role = 'follower';`
)

var (
	deleteTaskSQL      = `DELETE FROM task_info WHERE Name = ?;`
	checkTaskStatusSQL = `SELECT PersistData, Status FROM task_info where Name = ?;`
	updateTaskMsgSQL   = `UPDATE task_info SET DeployNode = ? WHERE Name = ?;`
	getTaskNumSQL      = `SELECT count(*) FROM task_info where DeployNode = ? and Status = 'Deployed';`
)

var (
	getMirrorsMsgSQL      = `SELECT name, ip_address, port, username, password, status, detail FROM warehouses LIMIT ?, ?;`
	getMirrorMsgSQL       = `SELECT name, ip_address, port, username, password, status, detail FROM warehouses where name = ?;`
	getMirrorsNameSQL     = `SELECT name FROM warehouses;`
	checkMirrorsNameSQL   = `SELECT name FROM warehouses where name = ?;`
	insertMirrorSQL       = `INSERT INTO warehouses (name, ip_address, port, username, password, status, detail) VALUES (?, ?, ?, ?, ?, ?, ?);`
	updateMirrorStatusSQL = `UPDATE warehouses SET status = ? WHERE name = ?`
	checkMirrorStatusSQL  = `SELECT status FROM warehouses where name = ?;`
	deleteMirrorSQL       = `DELETE FROM warehouses WHERE name = ?`
)
