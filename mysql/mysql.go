package mysql

import (
	"errors"
	"strconv"
	"strings"

	"gobigdrop/utils"

	"github.com/siddontang/go-mysql/client"
	"github.com/sirupsen/logrus"
)

// GetMySQLConn 获取MySQL连接
func GetMySQLConn(host string, port int, user string, password string, db string) (*client.Conn, error) {
	// 获取MySQL连接
	conn, err := client.Connect(host+":"+strconv.Itoa(port), user, password, db)
	if err != nil {
		return nil, err
	}

	// Ping测试
	if err := conn.Ping(); err != nil {
		return nil, err
	}

	return conn, nil
}

// GetMySQLVersion 获取MySQL版本
func GetMySQLVersion(c *client.Conn) (int, error) {
	// 获取MySQL version参数结果
	verres, err := c.Execute("show global variables like 'version'")
	defer verres.Close()

	if err != nil {
		return 0, err
	}

	verstr, err := verres.GetStringByName(0, "Value")
	if err != nil {
		return 0, err
	}

	versli := strings.Split(verstr, ".")
	ver, err := strconv.Atoi(versli[0] + versli[1])
	if err != nil {
		return 0, err
	}

	return ver, nil
}

// GetMySQLTableFile 获取MySQL表对应的文件
func GetMySQLTableFile(c *client.Conn, db string, table string) (map[string]string, error) {
	// 获取MySQL数据目录
	showres, err := c.Execute("show variables like 'datadir'")
	defer showres.Close()

	if err != nil {
		return nil, err
	}

	// 判断datadir是否以"/"结尾，是则去除
	datadir, err := showres.GetStringByName(0, "Value")
	if err != nil {
		return nil, err
	}

	if strings.HasSuffix(datadir, "/") == true {
		datadir = strings.TrimSuffix(datadir, "/")
	}

	// 判断ibd和frm文件是否存在
	filemap := make(map[string]string)

	ibdfile := datadir + "/" + db + "/" + table + ".ibd"
	if bv, err := utils.PathExists(ibdfile); err != nil {
		return nil, err
	} else {
		if bv == true {
			filemap["ibd"] = ibdfile
		} else {
			filemap["ibd"] = ""
		}
	}

	frmfile := datadir + "/" + db + "/" + table + ".frm"
	if bv, err := utils.PathExists(frmfile); err != nil {
		return nil, err
	} else {
		if bv == true {
			filemap["frm"] = frmfile
		} else {
			filemap["frm"] = ""
		}
	}

	return filemap, nil
}

// CheckMySQLTable 判断MySQL是否存在某张表
func CheckMySQLTable(c *client.Conn, db string, table string) (bool, error) {
	// 查询表信息
	cntres, err := c.Execute("select count(*) as count from information_schema.tables where table_schema=? and table_name=?", db, table)
	defer cntres.Close()

	if err != nil {
		return false, err
	}

	cnt, err := cntres.GetIntByName(0, "count")
	if err != nil {
		return false, err
	}

	// 如果行数为1，则返回true
	if cnt == 1 {
		return true, nil
	}

	nuerr := errors.New(db + "." + table + " check failed")
	return false, nuerr
}

// DropMySQLTable 在MySQL中进行DROP TABLE操作c
func DropMySQLTable(c *client.Conn, db string, table string) error {
	dropres, err := c.Execute("drop table " + db + "." + table)
	defer dropres.Close()

	if err != nil {
		return err
	}

	return nil
}

// MySQLSafeDrop MySQL删除大表流程，创建硬链接 -> DROP TABLE -> 删除硬链接
func MySQLSafeDrop(conn *client.Conn, db string, table string, log *logrus.Logger) {
	//  MySQL表存在则获取表对应的ibd和frm文件
	fmap, err := GetMySQLTableFile(conn, db, table)
	if err != nil {
		log.Error(err.Error())
		return
	} else {
		log.Info("Table: " + table + ", Get ibd & frm files successfully")
	}

	// 待操作文件slice
	filesli := []string{fmap["ibd"], fmap["frm"]}

	// 创建硬链接，后缀.hlink
	for i, f := range filesli {
		if hl, err := utils.CreateHardLink(f); err != nil {
			log.Error(err.Error())

			// rollback上一个硬链接的创建
			if i-1 >= 0 {
				if hl, err := utils.DropHardLink(filesli[i-1]); err != nil {
					log.Error(err.Error())
					return
				} else {
					log.Info("Hardlink: " + hl + ", Deleted")
				}
			}

			return
		} else {
			log.Info("Hardlink: " + hl + ", Created")
		}
	}

	// 在MySQL中进行DROP TABLE操作
	if err := DropMySQLTable(conn, db, table); err != nil {
		log.Error(err)
		return
	} else {
		log.Info("Table: " + table + ", Dropped")
	}

	// 删除硬链接，后缀.hlink
	for _, f := range filesli {
		if hl, err := utils.DropHardLink(f); err != nil {
			log.Error(err.Error())
			return
		} else {
			log.Info("Hardlink: " + hl + ", Deleted")
		}
	}
}

// MySQL80SafeDrop MySQL删除大表流程，创建硬链接 -> DROP TABLE -> 删除硬链接
func MySQL80SafeDrop(conn *client.Conn, db string, table string, log *logrus.Logger) {
	// MySQL表存在则获取表对应的ibd和frm文件
	fmap, err := GetMySQLTableFile(conn, db, table)
	if err != nil {
		log.Error(err.Error())
		return
	} else {
		log.Info("Table: " + table + ", Get ibd file successfully")
	}

	// 创建硬链接，后缀.hlink
	f := fmap["ibd"]
	if hl, err := utils.CreateHardLink(f); err != nil {
		log.Error(err.Error())
		return
	} else {
		log.Info("Hardlink: " + hl + ", Created")
	}

	// 在MySQL中进行DROP TABLE操作
	if err := DropMySQLTable(conn, db, table); err != nil {
		log.Error(err)
		return
	} else {
		log.Info("Table: " + table + ", Dropped")
	}

	// 删除硬链接，后缀.hlink
	if hl, err := utils.DropHardLink(f); err != nil {
		log.Error(err.Error())
		return
	} else {
		log.Info("Hardlink: " + hl + ", Deleted")
	}
}
