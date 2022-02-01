# 1.简介

PDFS是一个基于GFS思想的存储服务，通过将文件拆分成若干大小的块，并存放到不同存储服务器中，达到分布式存储的效果。

PDFS由管理服务器PDFS-Handler和存储服务器PDFS-Server组成。

# 2.使用

## 2.1.环境

- 操作系统：由于开发与测试都是在Linux上进行的，因此推荐Linux。不排除Windows会出现奇怪的问题。
- 数据库：Mysql（建议5.7版本），Redis
- Golang

## 2.2.下载项目

推荐下载到/usr/local/目录下，与文件树和文件块的默认存储位置一致。

> git clone https://github.com/Whalesong0622/PDFS /usr/local/

## 2.3.配置文件

请保证配置文件中redis和mysql的相关信息填写正确，若无法连接数据库，则软件无法启动。

第一次运行软件，分别进入PDFS-Handler和PDFS-Server目录下，使用go build命令，编译生成可执行文件。

初次运行会在根目录下生成config.json文件，其中PDFS-Handler的内容和含义如下：

```bash
{
	// 以下的地址表示的是ip和端口，即ip:port
    // handler的监听地址，用于云服务器等监听端口和外网端口不一致的情况
	"listen_addr":	"127.0.0.1:9999",
	// handler的访问地址
	"handler_addr": "127.0.0.1:9999",
	// handler的redis访问地址
	"handler_redis": "127.0.0.1:6379",
	// 文件目录namespace在handler服务器中的绝对路径
	"namespace_path": "/usr/local/PDFS/namespace",
	// 存放人员信息的Mysql用户名，密码，ip，端口，数据库名称和表名称
	"mysql_username": "root",
	"mysql_passwd": "123456",
	"mysql_ip": "127.0.0.1",
	"mysql_port": "3306",
	"mysql_dbname": "PDFS",
	"mysql_tbname": "PDFS_PEOPLE_TABLE"
}
```

PDFS-Server的内容和含义如下

```bash
{
	// 以下的地址表示的是ip和端口，即ip:port
	// server的监听地址，用于云服务器等监听端口和外网端口不一致的情况
	"listen_addr":"127.0.0.1:11111",
	// server的访问地址
    "server_addr": "127.0.0.1:11111",
    // handler的访问地址
	"handler_addr": "127.0.0.1:9999",
	// handler的redis访问地址
	"handler_redis": "127.0.0.1:6379",
	// 存放文件块的绝对路径
	"blocks_path": "/usr/local/PDFS/blocks"
}
```

配置完成后，运行编译生成的可执行文件。若无误，将会看到以下提示

```bash
2022/02/01 13:35:54 Loading config.
2022/02/01 13:35:54 Handler addr: 127.0.0.1:9999
2022/02/01 13:35:54 Listen addr: 127.0.0.1:9999
2022/02/01 13:35:54 Server addr: 127.0.0.1:11111
2022/02/01 13:35:54 Redis addr: 127.0.0.1:6379
2022/02/01 13:35:54 Namespace path: /usr/local/PDFS/namespace
2022/02/01 13:35:54 Found namespace dir.
2022/02/01 13:35:54 Server init success.
2022/02/01 13:35:54 Server start serving,listening to: 127.0.0.1:9999
```

# 3.试用&测试

PDFS-Debug-Client下有两个测试用客户端，分别对应管理服务器Handler和存储服务器Server。

其中管理服务器的使用需要在一个或多个存储服务器可用的情况下，才可使用。

请将想要上传的文件放到upload文件夹下。

下载的文件将会被存到download文件夹下。

## 3.1.存储服务器Server

```bash
Server-Debug-Client git:(2349226) ✗ go run Server-Debug-Client.go 
Enter your operation
1:Upload
2:Download
3:Delete
```

使用go run命令运行程序后，输入1、2、3可以执行相应操作。

- 1：上传文件
- 2：下载文件
- 3：删除文件

文件的存储位置为PDFS-Server下config.json的blocks_path。

## 3.2.管理服务器Handler

```bash
Handler-Debug-Client git:(2349226) ✗ go run Handler-Debug-Client.go 
Enter your operation
1:Create user
2:Login
3:Delete user
4:Change password
5:Upload file
6:Download file
7:Delete file
8:Add new path
9:Delete path
10:Ask files in path
Enter '0' to get operation list.
```

使用go run命令运行程序后，输入1、2、3可以执行相应操作。

- 1：新建用户
- 2：用户登陆
- 3：删除用户
- 4：更改密码
- 5：上传文件
- 6：下载文件
- 7：删除文件
- 8：新增路径
- 9：删除路径
- 10：请求路径下的文件

其中，除了1、2操作以外，其他操作都需要在登陆情况下进行。

```bash
1
Connection established.Enter new username and passwd.
username:
debug
passwd:
debug
Create new user success.
Enter '0' to get operation list.
2
Connection established.Enter your username and passwd.
username:
debug
passwd:
debug
Login success.
Cookie: [149 175 90 37 54 121 81 186 162 255 108 212 113 196 131 241 95 185 11 173]
```

登陆之后，服务器会返回长度为20的Cookie，其他操作需要利用这个Cookie来标识自己的身份。

# **有关路径的操作**

5-9操作都需要输入路径，该路径是**相对于该登录用户的路径，且包含文件/文件夹 名。**

用户的文件空间存在于2.3.中Handler的配置文件namespace_path下的文件夹，其中文件夹的名字为用户名。

如用户debug，则在默认配置下有/usr/local/bin/PDFS/namespace/debug/，他的根目录也是/usr/local/bin/PDFS/namespace/debug/。

若希望将文件example.txt存放到该用户下的dir文件夹，则需要提供的path为/dir/example.txt，绝对路径为/usr/local/bin/PDFS/namespace/user/dir/example.txt

如debug想要将GFS.pdf上传到其根目录下，则进行如下操作

```bash
5
Please make sure login before upload file.
Please put your file into "upload" diretory.
Enter filename and relative path.
filename:
GFS.pdf
path:
/GFS.pdf
Error when reading file. EOF
Send file finish
```

GFS.pdf会出现在namespace文件夹下的/debug文件夹中。

**注意，此处的文件只是用于构建文件树用，没有任何信息。真正的文件存储于存储服务器的blocks_path下。**

```
namespace git:(2349226) ✗ tree  
.
└── debug
    └── GFS.pdf
```

若debug想要在根目录下创建dir目录，并将GFS.pdf上传到该dir目录下，则进行如下操作

```
8
Please make sure login before adding new path.
Enter path.
path:
/dir
Add new path success.
Enter '0' to get operation list.
5
Please make sure login before upload file.
Please put your file into "upload" diretory.
Enter filename and relative path.
filename:
GFS.pdf 
path:
/dir/GFS.pdf
Error when reading file. EOF
Send file finish
```

GFS.pdf会出现在namespace文件夹下的/debug/dir文件夹中。

```bash
namespace git:(2349226) ✗ tree
.
└── debug
    ├── dir
    │   └── GFS.pdf
    └── GFS.pdf
```



# 4.PDFS-Handler协议

## 4.1.操作

PDFS-Handler为管理服务器，拥有以下操作：

- 1.新建用户
- 2.删除用户
- 3.修改用户密码
- 4.用户登陆
- 5.上传文件
- 6.读取文件
- 7.删除文件
- 8.新建路径
- 9.删除路径
- 10.请求路径下存在的文件
- 11.存储服务器请求注册（开发中）

## 4.2.协议

### 4.2.1.请求协议

协议为字节流，通过Socket进行信息的交互与传输。

协议头通用的第一位含义如下：

| 变量 | 位置下标 | 含义                                     |
| ---- | -------- | ---------------------------------------- |
| op   | [0,0]    | 表示请求服务器进行3.1.中的其中一项操作。 |

服务器被请求后会返回字节流，表示成功与否或其他信息。其中，通用的错误码有

| 变量              | 含义                               | 值   |
| ----------------- | ---------------------------------- | ---- |
| UNKNOWN_ERR       | 未知错误                           | 0    |
| COOKIES_NOT_FOUND | 没找到COOKIE，可能是过期或没有登陆 | 20   |

#### 4.2.1.1.新建用户

新建用户时，用户需要提供长度不超过20的字符串作为用户名和密码。

若提供的用户名或密码小于20，则需要用字节0进行填充。

| 变量     | 位置下标 | 含义                             |
| -------- | -------- | -------------------------------- |
| op（1）  | [0,0]    | 表示请求服务器进行新建用户操作。 |
| username | [1,20]   | 用户名字符串                     |
| passwd   | [21,40]  | 密码字符串                       |

#### 4.2.1.2.登陆

登陆时，用户需要输入密码进行身份确认。

| 变量     | 位置下标 | 含义                         |
| -------- | -------- | ---------------------------- |
| op（4）  | [0,0]    | 表示请求服务器进行登陆操作。 |
| username | [1,20]   | 用户名字符串                 |
| passwd   | [21,40]  | 密码字符串                   |

用户在登陆之后，服务器会对于该用户生成并返回一个长度为20的Cookie，有效期为一小时。

除了新建和登陆操作之外，该用户的请求中都需要利用该Cookie进行身份标识。

返回的信息有

| 变量           | 含义       | 值   |
| -------------- | ---------- | ---- |
| LOGIN_SUCCESS  | 登陆成功   | 7    |
| USER_NOT_EXIST | 用户不存在 | 8    |
| PASSWD_ERROR   | 密码错误   | 9    |

#### 4.2.1.3.删除用户

删除时，用户需要输入密码进行身份确认。

| 变量    | 位置下标 | 含义                         |
| ------- | -------- | ---------------------------- |
| op（2） | [0,0]    | 表示请求服务器进行登陆操作。 |
| cookie  | [1,20]   | 用于标识用户                 |
| passwd  | [21,40]  | 密码字符串                   |

返回的信息有

| 变量                  | 含义               | 值   |
| --------------------- | ------------------ | ---- |
| DEL_USER_FAILED       | 删除用户成功       | 3    |
| DEL_USER_PASSWD_ERROR | 密码错误，删除失败 | 4    |

#### 4.2.1.4.修改密码

修改密码时，用户需要输入密码进行身份确认。

| 变量    | 位置下标 | 含义                             |
| ------- | -------- | -------------------------------- |
| op（3） | [0,0]    | 表示请求服务器进行修改密码操作。 |
| cookie  | [1,20]   | 用于标识用户                     |
| passwd  | [21,40]  | 新的密码字符串                   |

返回的信息有

| 变量                  | 含义               | 值   |
| --------------------- | ------------------ | ---- |
| CHANGE_PASSWD_SUCCESS | 删除用户成功       | 5    |
| CHANGE_PASSWD_FAILED  | 密码错误，删除失败 | 6    |

#### 4.2.1.5.上传文件

上传文件需要提供cookie和path。

其中，path表示的希望将文件存放到的、相对于用户根目录的路径信息，包括文件名。

如配置文件中提到的namespace路径，PDFS-Handler会为每一个用户在该路径下创建一个文件夹表示该用户的文件空间，其中文件夹的名字为用户的名字。

如用户user，则在默认配置下有/usr/local/bin/PDFS/namespace/user/，他的根目录也是/usr/local/bin/PDFS/namespace/user/。

若希望将文件example.txt存放到该用户下的dir文件夹，则需要提供的path为/dir/example.txt，绝对路径为/usr/local/bin/PDFS/namespace/user/dir/example.txt

| 变量    | 位置下标 | 含义                         |
| ------- | -------- | ---------------------------- |
| op（5） | [0,0]    | 表示请求服务器进行上传操作。 |
| cookie  | [1,20]   | 用于标识用户。               |
| path    | [21,~]   | 路径字符串。                 |

上传文件的流程如下

- 客户端向Handler发送上表中的比特流
- 若解析无误，Handler会回送表示OK的字节255
- 客户端收到255之后，开始向建立的socket写入文件的字节流。**注意，若文件写入完毕后，客户端应该主动关闭socket，否则Handler会被阻塞。**

返回的信息有

| 变量             | 含义         | 值   |
| ---------------- | ------------ | ---- |
| WRITE_OP_SUCCESS | 上传文件成功 | 10   |

#### 4.2.1.6.下载文件

下载文件需要提供cookie和path。

path的含义参考3.2.1.5.上传文件

| 变量    | 位置下标 | 含义                         |
| ------- | -------- | ---------------------------- |
| op（6） | [0,0]    | 表示请求服务器进行下载操作。 |
| cookie  | [1,20]   | 用于标识用户。               |
| path    | [21,~]   | 路径字符串。                 |

下载文件的流程如下

- 客户端向Handler发送上表中的比特流
- 若解析无误，Handler会回送字节流，含义如下

| 变量                | 位置下标            | 含义                                                     |
| ------------------- | ------------------- | -------------------------------------------------------- |
| op（11）            | [0,0]               | 服务器返回码，表示这是下载文件的返回信息                 |
| blockname           | [1,64]              | 文件通过规则进行sha256得到的映射                         |
| blocknums           | [65,65]             | 表示该文件被分成了多少个块                               |
| ip_length_1         | [66,66]             | 文件第1部分所在服务器的地址的长度                        |
| ip_1                | [67,66+ip_length_1] | 文件第1部分所在服务器的地址，形如127.0.0.1:11111         |
| ...                 | ...                 | ...                                                      |
| ip_length_blocknums | ...                 | 文件第blocknums部分所在服务器的地址的长度                |
| ip_blocknums        | ...                 | 文件第blocknums部分所在服务器的地址，形如127.0.0.1:11111 |

- 客户端解析完字节流之后，对于每一个块，依次向存放了该块的存储服务器进行请求。其中向存储服务器请求的下载协议在4.2.中。
- 块的名字为

| 变量      | 位置下标 | 含义                                                         |
| --------- | -------- | ------------------------------------------------------------ |
| blockname | [0,63]   | 文件通过规则进行sha256得到的映射                             |
| '-'       | [64,64]  | 固定字符                                                     |
| blocknum  | [65,~]   | 想要请求的块的编号减一。<br />假设一个文件的blockname为example，它被分成了n个块<br />那么这n个块的名字分别为example-0,example-1...example-(n-1) |

返回的信息有

| 变量                | 含义                     | 值   |
| ------------------- | ------------------------ | ---- |
| READ_OP_RETURN      | 上传文件成功             | 11   |
| READ_FILE_NOT_EXIST | 下载文件失败，文件不存在 | 12   |

#### 4.2.1.7.删除文件

下载文件需要提供cookie和path。

path的含义参考3.2.1.5.上传文件

| 变量    | 位置下标 | 含义                             |
| ------- | -------- | -------------------------------- |
| op（7） | [0,0]    | 表示请求服务器进行删除文件操作。 |
| cookie  | [1,20]   | 用于标识用户。                   |
| path    | [21,~]   | 路径字符串。                     |

返回的信息有

| 变量               | 含义                     | 值   |
| ------------------ | ------------------------ | ---- |
| DEL_FILE_SUCCESS   | 删除文件成功             | 13   |
| DEL_FILE_NOT_EXIST | 删除文件失败，文件不存在 | 14   |

#### 4.2.1.8.新增路径

新增路径需要提供cookie和path。

path的含义参考3.2.1.5.上传文件

| 变量    | 位置下标 | 含义                             |
| ------- | -------- | -------------------------------- |
| op（8） | [0,0]    | 表示请求服务器进行新增路径操作。 |
| cookie  | [1,20]   | 用于标识用户。                   |
| path    | [21,~]   | 路径字符串。                     |

返回的信息有

| 变量                | 含义                     | 值   |
| ------------------- | ------------------------ | ---- |
| CREATE_PATH_SUCCESS | 新增路径成功             | 15   |
| CREATE_PATH_EXIST   | 创建路径失败，路径已存在 | 16   |

#### 4.2.1.9.删除路径

删除路径需要提供cookie和path。

path的含义参考3.2.1.5.上传文件

注意，删除路径会将路径下的所有文件一同删除。

| 变量    | 位置下标 | 含义                             |
| ------- | -------- | -------------------------------- |
| op（9） | [0,0]    | 表示请求服务器进行删除路径操作。 |
| cookie  | [1,20]   | 用于标识用户。                   |
| path    | [21,~]   | 路径字符串。                     |

返回的信息有

| 变量               | 含义                     | 值   |
| ------------------ | ------------------------ | ---- |
| DEL_PATH_SUCCESS   | 删除路径成功             | 17   |
| DEL_PATH_NOT_EXIST | 删除路径失败，路径不存在 | 18   |

#### 4.2.1.10.请求路径下存在的文件

需要提供cookie和path。

path的含义参考3.2.1.5.上传文件

| 变量      | 位置下标 | 含义                             |
| --------- | -------- | -------------------------------- |
| op（255） | [0,0]    | 表示请求服务器进行删除路径操作。 |
| cookie    | [1,20]   | 用于标识用户。                   |
| path      | [21,~]   | 路径字符串。                     |

返回的信息格式为

| 位置下标                 | 变量                      | 含义                                                   |
| ------------------------ | ------------------------- | ------------------------------------------------------ |
| [0,0]                    | ASK_FILES_IN_PATH（21）   | 表示该请求的返回信息。                                 |
| [1,1]                    | filenums                  | 表示该路径下有多少个文件。                             |
| [2,2]                    | file_type_1               | 第一个文件的类型。如果是文件则为1，文件夹则为2。       |
| [3,3]                    | file_name_length_1        | 第一个文件的名字的长度。                               |
| [4,3+file_name_length_1] | file_name_1               | 第一个文件的名字。                                     |
| ...                      | ...                       | ..                                                     |
| ...                      | file_type_filenums        | 第filenums个文件的类型。如果是文件则为1，文件夹则为2。 |
| ...                      | file_name_length_filenums | 第filenums个文件的名字的长度。                         |
| ...                      | file_name_filenums        | 第filenums个文件的名字。                               |

# 5.PDFS-Server协议

## 5.1.操作

PDFS-Server为存储服务器，拥有以下操作：

- 1.上传文件块
- 2.读取文件块
- 3.删除文件块

PDFS-Server统一将块存储到配置文件的blocks_path下。

## 5.2.协议

| 变量      | 位置下标 | 含义                           |
| --------- | -------- | ------------------------------ |
| op（1）   | [0,0]    | 表示请求服务器进行上传块操作。 |
| op（2）   | [0,0]    | 表示请求服务器进行下载块操作。 |
| op（3）   | [0,0]    | 表示请求服务器进行删除块操作。 |
| blockname | [1,~]    | 请求的块的名字。               |

返回值如下

| 变量           | 含义         | 值   |
| -------------- | ------------ | ---- |
| OK             | 操作成功。   | 0    |
| UNKNOWN_ERR    | 未知错误。   | 1    |
| FILE_NOT_EXIST | 文件不存在。 | 2    |

