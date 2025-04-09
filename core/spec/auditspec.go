package spec

/*
 定义规范来源：
 https://web.git.kernel.org/pub/scm/linux/kernel/git/stable/linux.git/tree/include/uapi/linux/audit.h?h=linux-3.14.y
*/

// 命令控制区块 (1000-1099)
const (
	AUDIT_GET            = 1000 + iota // 获取状态
	AUDIT_SET                          // 设置状态
	AUDIT_LIST                         // 列出系统调用规则（已弃用）
	AUDIT_ADD                          // 添加系统调用规则（已弃用）
	AUDIT_DEL                          // 删除系统调用规则（已弃用）
	AUDIT_USER                         // 用户空间消息（已弃用）
	AUDIT_LOGIN                        // 定义登录ID和信息
	AUDIT_WATCH_INS                    // 插入文件/目录监视项
	AUDIT_WATCH_REM                    // 移除文件/目录监视项
	AUDIT_WATCH_LIST                   // 列出所有文件/目录监视
	AUDIT_SIGNAL_INFO                  // 获取审计信号发送者信息
	AUDIT_ADD_RULE                     // 添加系统调用过滤规则
	AUDIT_DEL_RULE                     // 删除系统调用过滤规则
	AUDIT_LIST_RULES                   // 列出过滤规则
	AUDIT_TRIM                         // 修剪监视树中的垃圾
	AUDIT_MAKE_EQUIV                   // 追加到监视树
	AUDIT_TTY_GET                      // 获取TTY审计状态
	AUDIT_TTY_SET                      // 设置TTY审计状态
	AUDIT_SET_FEATURE                  // 启用/禁用审计功能
	AUDIT_GET_FEATURE                  // 获取已启用功能
	AUDIT_FEATURE_CHANGE               // 功能变更日志
)

// 用户空间消息 (1100-1199)
const (
	AUDIT_FIRST_USER_MSG = 1100
	AUDIT_USER_AVC       = 1107 // 差异化过滤的AVC消息
	AUDIT_USER_TTY       = 1124 // 非ICANON TTY输入
	AUDIT_LAST_USER_MSG  = 1199
)

// 审计守护进程消息 (1200-1299)
const (
	AUDIT_DAEMON_START  = 1200 // 守护进程启动记录
	AUDIT_DAEMON_END    = 1201 // 正常停止记录
	AUDIT_DAEMON_ABORT  = 1202 // 错误停止记录
	AUDIT_DAEMON_CONFIG = 1203 // 配置变更
)

// 内核事件消息 (1300-1399)
const (
	AUDIT_SYSCALL       = 1300 + iota // 系统调用事件
	_                                 // 1301（已弃用）
	AUDIT_PATH                        // 文件路径信息
	AUDIT_IPC                         // IPC记录
	AUDIT_SOCKETCALL                  // socketcall参数
	AUDIT_CONFIG_CHANGE               // 系统配置变更
	AUDIT_SOCKADDR                    // 作为参数的sockaddr
	AUDIT_CWD                         // 当前工作目录
	_
	AUDIT_EXECVE // 1309: execve参数
	_
	AUDIT_IPC_SET_PERM  // 1311: IPC新权限记录
	AUDIT_MQ_OPEN       // POSIX消息队列打开
	AUDIT_MQ_SENDRECV   // 消息发送/接收
	AUDIT_MQ_NOTIFY     // 消息通知
	AUDIT_MQ_GETSETATTR // 属性操作
	AUDIT_KERNEL_OTHER  // 第三方模块使用
	AUDIT_FD_PAIR       // 管道/socketpair记录
	AUDIT_OBJ_PID       // ptrace目标
	AUDIT_TTY           // 管理终端输入
	AUDIT_EOE           // 多记录事件结束
	AUDIT_BPRM_FCAPS    // 权能提升信息
	AUDIT_CAPSET        // capset参数
	AUDIT_MMAP          // mmap描述符和标志
	AUDIT_NETFILTER_PKT // 网络过滤器数据包
	AUDIT_NETFILTER_CFG // 网络配置变更
	AUDIT_SECCOMP       // 安全计算事件
)

// SELinux相关 (1400-1499)
const (
	AUDIT_AVC               = 1400 + iota // SELinux访问控制
	AUDIT_SELINUX_ERR                     // 内部错误
	AUDIT_AVC_PATH                        // AVC路径信息
	AUDIT_MAC_POLICY_LOAD                 // 策略文件加载
	AUDIT_MAC_STATUS                      // 强制状态变更
	AUDIT_MAC_CONFIG_CHANGE               // 布尔值变更
	AUDIT_MAC_UNLBL_ALLOW                 // 允许未标记流量
	AUDIT_MAC_CIPSOV4_ADD                 // 添加CIPSO条目
	AUDIT_MAC_CIPSOV4_DEL                 // 删除CIPSO条目
	AUDIT_MAC_MAP_ADD                     // 添加域映射
	AUDIT_MAC_MAP_DEL                     // 删除域映射
	// 1411-1414为保留值
	AUDIT_MAC_IPSEC_EVENT  = 1415 // IPSec事件
	AUDIT_MAC_UNLBL_STCADD = 1416 // 添加静态标签
	AUDIT_MAC_UNLBL_STCDEL = 1417 // 删除静态标签
)

// 内核异常记录 (1700-1799)
const (
	AUDIT_FIRST_KERN_ANOM_MSG = 1700 + iota
	AUDIT_ANOM_PROMISCUOUS    // 混杂模式变更
	AUDIT_ANOM_ABEND          // 进程异常终止
	AUDIT_ANOM_LINK           // 可疑文件链接
	AUDIT_LAST_KERN_ANOM_MSG  = 1799
)

// 完整性验证 (1800-1899)
const (
	AUDIT_INTEGRITY_DATA     = 1800 + iota // 数据完整性
	AUDIT_INTEGRITY_METADATA               // 元数据完整性
	AUDIT_INTEGRITY_STATUS                 // 启用状态
	AUDIT_INTEGRITY_HASH                   // 哈希类型
	AUDIT_INTEGRITY_PCR                    // PCR失效消息
	AUDIT_INTEGRITY_RULE                   // 策略规则
)
