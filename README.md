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
# ./gobigdrop --host=xxx.xxx.xxx.xxx --port=3306 --user=xxx --askpass --db=xxx --table=xxx
Enter password: 
INFO[2021-01-06 16:27:03] MySQL Connection: user@host:port/db 
INFO[2021-01-06 16:27:03] MySQL Version: 57                            
INFO[2021-01-06 16:27:03] Table: table, Checked                       
INFO[2021-01-06 16:27:03] Table: table, Get ibd & frm files successfully 
INFO[2021-01-06 16:27:03] Hardlink: /../../db/table.ibd.hlink, Created 
INFO[2021-01-06 16:27:03] Hardlink: /../../db/table.frm.hlink, Created 
INFO[2021-01-06 16:27:03] Table: table, Dropped                       
INFO[2021-01-06 16:27:03] Hardlink: /../../db/table.ibd.hlink, Deleted 
INFO[2021-01-06 16:27:03] Hardlink: /../../db/table.frm.hlink, Deleted 
INFO[2021-01-06 16:27:03] Done
```
