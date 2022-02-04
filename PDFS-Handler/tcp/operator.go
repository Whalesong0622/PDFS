package tcp

// 操作

const NEW_USER_OP byte = 1      // 新建用户
const DEL_USER_OP byte = 2      // 删除用户
const CHANGE_PASSWD_OP byte = 3 // 修改密码
const LOGIN_OP byte = 4         // 用户登陆
const WRITE_OP byte = 5         // 上传文件
const READ_OP byte = 6          // 读取文件
const DEL_OP byte = 7           // 删除文件
const ADD_DIR_OP byte = 8       // 新建路径
const DEL_DIR_OP byte = 9       // 删除路径
const ASK_FILES_OP byte = 255   // 请求该目录下文件
