package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"gobigdrop/mysql"
	"gobigdrop/utils"

	"github.com/howeyc/gopass"
)

var (
	host     string
	port     int
	user     string
	password string
	database string
	table    string
	log      string
	loglevel string
	askpass  bool
)

// usage 自定义usage信息
func usage() {
	fmt.Fprintf(
		os.Stderr,
		"gobigdrop version: 1.0.0\nUsage: gobigdrop [--host host] [--port port] [--user user] [--password password] [--db database] [--table table] [--log logfile] [--level loglevel] [--askpass]\n\nOptions:\n")
	flag.PrintDefaults()
}

func init() {
	flag.StringVar(&host, "host", "", "MySQL host.")
	flag.IntVar(&port, "port", 3306, "MySQL port.")
	flag.StringVar(&user, "user", "", "MySQL user.")
	flag.StringVar(&password, "password", "", "MySQL password.")
	flag.StringVar(&database, "db", "", "MySQL database.")
	flag.StringVar(&table, "table", "", "MySQL table.")
	flag.StringVar(&log, "log", "os.Stdout", "gobigdrop log file.")
	flag.StringVar(&loglevel, "level", "INFO", "gobigdrop log level.")
	flag.BoolVar(&askpass, "askpass", false, "ask MySQL password.")

	// 改变默认的 Usage
	flag.Usage = usage
}

func main() {
	// 自定义Usage
	flag.Usage = usage
	flag.Parse()

	// 是否交互获取password
	if askpass == true {
		passwd, err := gopass.GetPasswdPrompt("Enter password: ", false, os.Stdin, os.Stdout)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		password = string(passwd)
	}

	// 开启logger
	logger, f, err := utils.Log(log, loglevel, false)
	defer f.Close()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// 获取DB连接
	conn, err := mysql.GetMySQLConn(host, port, user, password, database)
	defer conn.Close()

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	} else {
		logger.Info("MySQL Connection: " + user + "@" + host + ":" + strconv.Itoa(port) + "/" + database)
	}

	// 查询MySQL版本
	verint, err := mysql.GetMySQLVersion(conn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	} else {
		logger.Info("MySQL Version: " + strconv.Itoa(verint))
	}

	// 验证待MySQL表元信息
	tabbool, err := mysql.CheckMySQLTable(conn, database, table)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	} else {
		logger.Info("Table: " + table + ", Checked")
	}

	if tabbool == true {
		switch verint {
		case 56:
			mysql.MySQLSafeDrop(conn, database, table, logger)
		case 57:
			mysql.MySQLSafeDrop(conn, database, table, logger)
		case 80:
			mysql.MySQL80SafeDrop(conn, database, table, logger)
		}
	} else {
		// 待删除表不存在则打印日志退出
		logger.Info("Table: " + table + ", Does not exist")
		os.Exit(1)
	}

	// end
	logger.Info("Done")
}
