# gobigdrop

gobigdrop适用于MySQL大表删除，利用MySQL Drop操作机制结合Linux硬链接对大表进行删除，旨在减少线上系统的性能抖动。

### 安装说明


#### 二进制免安装




### 使用说明


```
# ./gobigdrop --help
gobigdrop version: 1.0.0
Usage: gobigdrop [-h host] [-P port] [-u user] [-p password] [-d database] [-t table] [--log logfile] [--level loglevel] [--askpass]

Options:
  -askpass
        ask MySQL password.
  -db string
        MySQL database.
  -host string
        MySQL host.
  -level string
        gobigdrop log level. (default "INFO")
  -log string
        gobigdrop log file. (default "os.Stdout")
  -password string
        MySQL password.
  -port int
        MySQL port. (default 3306)
  -table string
        MySQL table.
  -user string
        MySQL user.
```


#### 简单示例


```
# ./gobigdrop --host=148.70.127.235 --port=3307 --user=joyee --askpass --db=test --table=usertb
Enter password: 
INFO[2021-01-08 13:32:40] MySQL Connection: joyee@148.70.127.235:3307/test 
INFO[2021-01-08 13:32:40] MySQL Version: 57                            
INFO[2021-01-08 13:32:40] Checked                                       Table=usertb
INFO[2021-01-08 13:32:40] Get ibd & frm files successfully              Table=usertb
INFO[2021-01-08 13:32:40] Created                                       Hardlink=/data/mysql/3307/data/test/usertb.ibd.hlink
INFO[2021-01-08 13:32:40] Created                                       Hardlink=/data/mysql/3307/data/test/usertb.frm.hlink
INFO[2021-01-08 13:32:40] Dropped                                       Table=usertb
INFO[2021-01-08 13:32:40] Deleted                                       Hardlink=/data/mysql/3307/data/test/usertb.ibd.hlink
INFO[2021-01-08 13:32:40] Deleted                                       Hardlink=/data/mysql/3307/data/test/usertb.frm.hlink
INFO[2021-01-08 13:32:40] Done 
```
