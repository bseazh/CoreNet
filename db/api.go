package db

import (
	mydb "FileStoreNetDisk_v3/db/mysql"
	"fmt"
	"log"
	"strconv"
)

type TableDirectory struct {
	ID           int
	DocumentName string
	Pid          int
	IsFile       int
}

// GetDirectory : 通过 username 获取文件目录
func GetDirectory(username string) ([]TableDirectory, error) {

	stmt, err := mydb.DBConn().Prepare(
		"select id , document_name , parent_id  , is_file " +
			" from tbl_document where user_name=? Order by is_file ASC ;")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(username)
	if err != nil {
		return nil, err
	}

	var userFiles []TableDirectory
	for rows.Next() {
		file := TableDirectory{}
		err = rows.Scan(&file.ID, &file.DocumentName, &file.Pid, &file.IsFile)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		userFiles = append(userFiles, file)
	}
	return userFiles, nil
}

// RenameFileAPI : 根据 ( username , id , newName ) 修改名称
func RenameFileAPI(username, id, newName string) error {
	stmt, err := mydb.DBConn().Prepare(
		"update tbl_document set document_name = ? " +
			" where user_name = ? and id = ? ;")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(newName, username, id)

	if err != nil {
		return err
	}
	return nil
}

// DeleteFileAPI : 根据 ( username , id  ) 删除文件
func DeleteFileAPI(username, id string) error {
	stmt, err := mydb.DBConn().Prepare(
		"delete from tbl_document " +
			" where user_name = ? and id = ? ;")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, id)

	if err != nil {
		return err
	}
	return nil
}

// RemoveDocumentAPI : 通过 ( documentName username ) 删除之间的索引 ;
func RemoveDocumentAPI(username string, ID int) bool {

	isFile, err := GetIsFile(username, ID)
	if err != nil {
		log.Println("获取 IsFile 出错!!! " + err.Error())
		return false
	}

	if isFile == 1 {
		// 文件
		id := strconv.Itoa(ID)
		err = DeleteFileAPI(username, id)
		if err != nil {
			return false
		} else {
			return true
		}
	} else {
		// 文件夹
		stmt, err := mydb.DBConn().Prepare(
			"DELETE FROM tbl_document WHERE user_name = ? and id = ?  ;")
		if err != nil {
			fmt.Println("Failed to prepare statement, err:" + err.Error())
			return false
		}
		defer stmt.Close()

		_, err = stmt.Exec(username, ID)
		if err != nil {
			fmt.Println(err.Error())
			return false
		}

		// 然后, 查找 当前 文件夹 里面 的东西
		documents, _ := GetDocumentList(username, ID)

		if documents == nil {
			return true
		}
		for _, document := range documents {
			RemoveDocumentAPI(username, document.ID)
		}

		return true
	}

}

// GetIsFile : 获取相应文件的属性（ 文件 | 文件夹 ）
func GetIsFile(username string, ID int) (int, error) {
	stmt, err := mydb.DBConn().Prepare(
		"select is_file from tbl_document where user_name = ? and id = ?;")
	if err != nil {
		fmt.Println("Failed to prepare statement, err:" + err.Error())
		return -1, err
	}
	defer stmt.Close()

	var isfile int
	err = stmt.QueryRow(username, ID).Scan(&isfile)
	if err != nil {
		fmt.Println(err.Error())
		return -1, err
	}
	return isfile, nil
}

// GetFileSha1 : 获取相应文件的属性（ 文件 | 文件夹 ）
func GetFileSha1(username string, ID int) (string, error) {
	stmt, err := mydb.DBConn().Prepare(
		"select file_sha1 from tbl_document where user_name = ? and id = ?;")
	if err != nil {
		fmt.Println("Failed to prepare statement, err:" + err.Error())
		return "", err
	}
	defer stmt.Close()

	var fileSha1 string
	err = stmt.QueryRow(username, ID).Scan(&fileSha1)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	return fileSha1, nil
}
