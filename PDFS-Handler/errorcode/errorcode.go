package errorcode

// 返回值

// 未知错误
const UNKNOWN_ERR byte = 0

// 新增用户成功
const NEW_USER_SUCCESS byte = 1

// 用户已存在（用于新建用户时冲突）
const USER_EXIST byte = 2

// 删除用户成功
const DEL_USER_FAILED byte = 3

// 删除用户密码核对失败
const DEL_USER_PASSWD_ERROR byte = 4

// 修改密码成功
const CHANGE_PASSWD_SUCCESS byte = 5

// 修改密码失败
const CHANGE_PASSWD_FAILED byte = 6

// 登录成功
const LOGIN_SUCCESS byte = 7

// 用户不存在
const USER_NOT_EXIST byte = 8

// 密码错误
const PASSWD_ERROR byte = 9

// 上传文件成功
const WRITE_OP_SUCCESS byte = 10

// 下载文件返回数据
const READ_OP_RETURN byte = 11

// 下载文件失败，文件不存在
const READ_FILE_NOT_EXIST byte = 12

// 删除文件成功
const DEL_FILE_SUCCESS byte = 13

// 删除文件失败，文件不存在
const DEL_FILE_NOT_EXIST byte = 14

// 创建路径成功
const CREATE_PATH_SUCCESS byte = 15

// 创建路径失败，路径已存在
const CREATE_PATH_EXIST byte = 16

// 删除路径成功
const DEL_PATH_SUCCESS byte = 17

// 创建路径失败，路径已存在
const DEL_PATH_NOT_EXIST byte = 18

// 权限不足
const NO_ACCESS_CONTROL byte = 19

// 没找到COOKIE
const COOKIES_NOT_FOUND byte = 20

// 请求文件目录下的文件或文件夹
const ASK_FILES_IN_PATH byte = 21

// 请求文件目录下的文件或文件夹失败，路径不存在
const ASK_FILES_IN_PATH_FAILED byte = 22

//
const OK byte = 255
