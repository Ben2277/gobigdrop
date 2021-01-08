# gobigdrop

gobigdrop适用于MySQL大表删除，利用MySQL Drop操作机制结合Linux硬链接对大表进行删除，旨在减少线上系统的性能抖动。

&nbsp;
### 安装说明
[二进制免安装](https://github.com/Ben2277/gobigdrop/releases)

&nbsp;
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

&nbsp;
#### 简单示例
```
# ./gobigdrop --host=xxx.xxx.xxx.xxx --port=xxxx --user=xxx --askpass --db=xxx --table=xxx
Enter password: 
INFO[2021-01-08 13:32:40] MySQL Connection: xxx@xxx.xxx.xxx.xxx:xxx/xxx
INFO[2021-01-08 13:32:40] MySQL Version: xxx                            
INFO[2021-01-08 13:32:40] Checked                                       Table=xxx
INFO[2021-01-08 13:32:40] Get ibd & frm files successfully              Table=xxx
INFO[2021-01-08 13:32:40] Created                                       Hardlink=/xxx/xxx.ibd.hlink
INFO[2021-01-08 13:32:40] Created                                       Hardlink=/xxx/xxx.frm.hlink
INFO[2021-01-08 13:32:40] Dropped                                       Table=xxx
INFO[2021-01-08 13:32:40] Deleted                                       Hardlink=/xxx/xxx.ibd.hlink
INFO[2021-01-08 13:32:40] Deleted                                       Hardlink=/xxx/xxx.frm.hlink
INFO[2021-01-08 13:32:40] Done 
```
