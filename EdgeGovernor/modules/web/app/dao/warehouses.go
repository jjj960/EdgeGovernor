package dao

//import (
//	"EdgeGovernor/pkg/constants"
//	"EdgeGovernor/pkg/database/mariadb"
//	"log"
//)
//
//func GetMirrorsMsg(start int, end int) []map[string]string {
//	rows, err := mariadb.Db.Query(getMirrorsMsgSQL, start, end-start)
//	if err != nil {
//		log.Println(err)
//	}
//	defer rows.Close()
//	var result []map[string]string
//	// 遍历结果集
//	for rows.Next() {
//		var name, ip_address, port, username, password, status, detail string
//		err := rows.Scan(&name, &ip_address, &port, &username, &password, &status, &detail)
//		if err != nil {
//			log.Println(err)
//		}
//		if name == "DockerHub" {
//			ip_address, port, username, password, status = "-", "-", "-", "-", "是"
//		}
//		node1 := make(map[string]string)
//		node1["mirrorRes"] = name
//		node1["mirrorIP"] = ip_address
//		node1["mirrorPort"] = port
//		node1["loginName"] = username
//		node1["password"] = password
//		node1["use"] = status
//		node1["remark"] = detail
//		result = append(result, node1)
//	}
//
//	// 检查是否有错误发生
//	if err = rows.Err(); err != nil {
//		log.Println(err)
//	}
//
//	return result
//}
//
//func GetMirrorsName() []string {
//	rows, err := mariadb.Db.Query(getMirrorsNameSQL)
//	if err != nil {
//		log.Println(err)
//	}
//	defer rows.Close()
//	var result []string
//	// 遍历结果集
//	for rows.Next() {
//		var name string
//		err := rows.Scan(&name)
//		if err != nil {
//			log.Println(err)
//		}
//		result = append(result, name)
//	}
//
//	// 检查是否有错误发生
//	if err = rows.Err(); err != nil {
//		log.Println(err)
//	}
//
//	return result
//}
//
//func CheckMirrorNameUniqueness(mirrorName string) bool {
//	rows, err := mariadb.Db.Query(checkMirrorsNameSQL, mirrorName)
//	if err != nil {
//		log.Println(err)
//	}
//	defer rows.Close()
//	var count int
//	// 遍历结果集
//	for rows.Next() {
//		count++
//	}
//	// 检查是否有错误发生
//	if err = rows.Err(); err != nil {
//		log.Println(err)
//	}
//
//	return count > 1
//}
//
//func AddMirror(mirrorName string, ip string, port string, username string, password string, status string, detail string) {
//	// 执行插入操作
//	_, err := mariadb.Db.Exec(insertMirrorSQL, mirrorName, ip, port, username, password, status, detail)
//	if err != nil {
//		log.Println(err)
//		return
//	}
//	constants.MirrorCount++
//}
//
//func ChangeMirrorStatus(mirrorName string) {
//	status := checkMirrorStatus(mirrorName)
//	var newStatus string
//	if status == "是" {
//		newStatus = "否"
//	} else {
//		newStatus = "是"
//	}
//	_, err := mariadb.Db.Exec(updateMirrorStatusSQL, newStatus, mirrorName)
//	if err != nil {
//		log.Println(err)
//	}
//}
//
//func checkMirrorStatus(mirrorName string) string {
//	rows, err := mariadb.Db.Query(checkMirrorStatusSQL, mirrorName)
//	if err != nil {
//		log.Println(err)
//	}
//	defer rows.Close()
//	var status string
//	// 遍历结果集
//	for rows.Next() {
//		err := rows.Scan(&status)
//		if err != nil {
//			log.Println(err)
//		}
//	}
//	// 检查是否有错误发生
//	if err = rows.Err(); err != nil {
//		log.Println(err)
//	}
//
//	return status
//}
//
//func CheckMirrorNameisExist(mirrorName string) bool {
//	rows, err := mariadb.Db.Query(checkMirrorsNameSQL, mirrorName)
//	if err != nil {
//		log.Println(err)
//	}
//	defer rows.Close()
//	var count int
//	// 遍历结果集
//	for rows.Next() {
//		count++
//	}
//	// 检查是否有错误发生
//	if err = rows.Err(); err != nil {
//		log.Println(err)
//	}
//
//	return count >= 1
//}
//
//func DeleteMirror(mirrorName string) {
//	_, err := mariadb.Db.Exec(deleteMirrorSQL, mirrorName)
//	if err != nil {
//		log.Println(err)
//	}
//}
//
//func SearchMirror(mirrorName string) map[string]string {
//	rows, err := mariadb.Db.Query(getMirrorMsgSQL, mirrorName)
//	if err != nil {
//		log.Println(err)
//	}
//	defer rows.Close()
//	var node1 = make(map[string]string)
//	// 遍历结果集
//	for rows.Next() {
//		var name, ip_address, port, username, password, status, detail string
//		err := rows.Scan(&name, &ip_address, &port, &username, &password, &status, &detail)
//		if err != nil {
//			log.Println(err)
//		}
//		node1["mirrorRes"] = name
//		node1["mirrorIP"] = ip_address
//		node1["mirrorPort"] = port
//		node1["loginName"] = username
//		node1["password"] = password
//		node1["use"] = status
//		node1["remark"] = detail
//	}
//
//	// 检查是否有错误发生
//	if err = rows.Err(); err != nil {
//		log.Println(err)
//	}
//
//	return node1
//}
