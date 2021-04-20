package db

import (
	mydb "FileStoreNetDisk_v3/db/mysql"
	"fmt"
	"log"
)

// CreateDocument : 创建新的文件夹
func CreateDocument(username, documentName string, parentID int, documentSize int64) bool {

	stmt, err := mydb.DBConn().Prepare(
		"insert ignore into tbl_document (`user_name`,`document_name`,`parent_id`) values (?,?,?);")
	if err != nil {
		fmt.Println("Failed to insert, err:" + err.Error())
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(username, documentName, parentID)
	if err != nil {
		fmt.Println("Failed to insert, err:" + err.Error())
		return false
	}

	if rowsAffected, err := ret.RowsAffected(); nil == err && rowsAffected > 0 {
		UpdateDocument(username, documentName, parentID, documentSize)
		return true
	}
	return false
}

// OpenFolder : 通过 ( username documentName parentID) 找到文件(夹)更新其对应的文件(夹)名;
func OpenFolder(username, documentName string, parentID int) (int, error) {
	id, err := QueryUserID(username, parentID, documentName)
	return id, err
}

// GoUpFolder : 通过( username , parentID ) 找到上一层的文件夹;
func GoUpFolder(username string, parentID int) (int, error) {
	stmt, err := mydb.DBConn().Prepare(
		"select parent_id " +
			" from tbl_document where user_name=? and id = ? ")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	var resID int
	err = stmt.QueryRow(username, parentID).Scan(&resID)
	if err != nil {
		return resID, err
	}
	return resID, err
}

// RenameDocument : 通过 ( documentName username ) 找到文件(夹)更新其对应的文件(夹)名;
func RenameDocument(username, documentName, filename string) bool {

	stmt, err := mydb.DBConn().Prepare(
		"UPDATE tbl_document SET document_name = ? where user_name = ? and document_name = ?")
	if err != nil {
		fmt.Println("Failed to prepare statement, err:" + err.Error())
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(filename, username, documentName)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

// TableDocument : 文件表结构体
type TableDocument struct {
	ID           int
	UserName     string
	FileSha1     string
	ParentID     int
	DocumentName string
	DocumentSize int64
	UploadAt     string
	IsFile       int
}

// GetDocumentList : 从mysql批量获取文件元信息
func GetDocumentList(username string, parentID int) ([]TableDocument, error) {

	stmt, err := mydb.DBConn().Prepare(
		"select id , document_name , file_sha1 , document_size , update_at , is_file " +
			" from tbl_document where user_name=? and parent_id = ? order by is_file ASC")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(username, parentID)
	if err != nil {
		return nil, err
	}

	var userFiles []TableDocument
	for rows.Next() {
		file := TableDocument{}
		file.UserName = username
		file.ParentID = parentID
		err = rows.Scan(&file.ID, &file.DocumentName, &file.FileSha1, &file.DocumentSize, &file.UploadAt, &file.IsFile)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		userFiles = append(userFiles, file)
	}
	return userFiles, nil

}

// GetDocumentID : 获取文件夹的ID;
func GetDocumentID(username string, parentID int, documentName string) (int, error) {

	stmt, err := mydb.DBConn().Prepare(
		"select id " +
			" from tbl_document where user_name=? and parent_id = ? and document_name = ? ")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var res int
	err = stmt.QueryRow(username, parentID, documentName).Scan(&res)
	if err != nil {
		return 0, err
	}
	return res, nil

}

// RemoveDocument : 通过 ( documentName username ) 删除之间的索引 ;
func RemoveDocument(username, documentName string, parentID int, isFile int, documentSize int64) bool {

	ID, err := GetDocumentID(username, parentID, documentName)
	if err != nil {
		log.Println("获取ID出错!!! " + err.Error())
		return false
	}

	if isFile == 1 {
		// 文件
		stmt, err := mydb.DBConn().Prepare(
			"DELETE FROM tbl_document WHERE user_name = ? and document_name = ? and parent_ID = ? ;")
		if err != nil {
			fmt.Println("Failed to prepare statement, err:" + err.Error())
			return false
		}
		defer stmt.Close()

		_, err = stmt.Exec(username, documentName, parentID)
		if err != nil {
			fmt.Println(err.Error())
			return false
		}
		UpdateDocument(username, documentName, parentID, -documentSize)
		return true
	} else {
		// 文件夹
		// 首先删除文件夹
		stmt, err := mydb.DBConn().Prepare(
			"DELETE FROM tbl_document WHERE user_name = ? and document_name = ? and parent_ID = ? ;")
		if err != nil {
			fmt.Println("Failed to prepare statement, err:" + err.Error())
			return false
		}
		defer stmt.Close()

		_, err = stmt.Exec(username, documentName, parentID)
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
			log.Printf("db/document/RemoveDocument/Deletedocument DocumentName: %s , ParentID : %d , IsFile : %v\n",
				document.UserName, ID, document.IsFile)
			RemoveDocument(document.UserName, document.DocumentName,
				ID, document.IsFile, -document.DocumentSize)
		}

		return true
	}

}

// UpdateDocument : 增删文件 对当前文件夹大小进行修正
func UpdateDocument(username, documentName string, parentID int, fileSize int64) bool {
	stmt, err := mydb.DBConn().Prepare(
		"UPDATE tbl_document SET document_size = document_size + ? " +
			"where user_name = ? and document_name = ? and parent_ID = ? ")
	if err != nil {
		fmt.Println("Failed to prepare statement, err:" + err.Error())
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(fileSize, username, documentName, parentID)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}
