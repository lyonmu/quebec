package code

var (
	// 登录相关
	CaptchaGeneratorFailed = Response{Code: 51001, Message: "验证码生成失败"}
	UserNotExists          = Response{Code: 51002, Message: "用户不存在"}
	UserPasswordError      = Response{Code: 51003, Message: "密码错误"}
	InvalidCaptcha         = Response{Code: 51004, Message: "无效的验证码"}
	CaptchaTooFrequent     = Response{Code: 51023, Message: "验证码请求过于频繁，请稍后再试"}
	// 角色相关
	RoleNotExists      = Response{Code: 51005, Message: "角色不存在"}
	RoleDisabled       = Response{Code: 51006, Message: "该角色被禁用"}
	RoleAddFailed      = Response{Code: 51007, Message: "角色添加失败"}
	RoleDelFailed      = Response{Code: 51008, Message: "角色删除失败"}
	RoleEditFailed     = Response{Code: 51009, Message: "角色编辑失败"}
	RoleQueryFailed    = Response{Code: 51010, Message: "角色查询失败"}
	RoleNameDuplicate  = Response{Code: 51011, Message: "角色名称重复"}
	RoleEnableFailed   = Response{Code: 51012, Message: "角色启停失败"}
	RoleSystemNotAllow = Response{Code: 51013, Message: "系统角色不可操作"}

	// 用户相关
	UserSystemNotAllow     = Response{Code: 51014, Message: "系统用户不可操作"}
	UserAddFailed          = Response{Code: 51015, Message: "用户添加失败"}
	UserDelFailed          = Response{Code: 51016, Message: "用户删除失败"}
	UserEditFailed         = Response{Code: 51017, Message: "用户编辑失败"}
	UserQueryFailed        = Response{Code: 51018, Message: "用户查询失败"}
	UserNameDuplicate      = Response{Code: 51019, Message: "用户名称重复"}
	UserEnableFailed       = Response{Code: 51020, Message: "用户启停失败"}
	PasswordIncorrect      = Response{Code: 51021, Message: "密码不正确"}
	UserPasswordEditFailed = Response{Code: 51022, Message: "修改密码失败"}

	// 菜单相关
	MenuNotExists      = Response{Code: 51024, Message: "菜单不存在"}
	MenuAddFailed      = Response{Code: 51025, Message: "菜单添加失败"}
	MenuDelFailed      = Response{Code: 51026, Message: "菜单删除失败"}
	MenuEditFailed     = Response{Code: 51027, Message: "菜单编辑失败"}
	MenuQueryFailed    = Response{Code: 51028, Message: "菜单查询失败"}
	MenuNameDuplicate  = Response{Code: 51029, Message: "菜单名称重复"}
	MenuEnableFailed   = Response{Code: 51030, Message: "菜单启停失败"}

	// 角色菜单绑定相关
	RoleMenuBindFailed = Response{Code: 51031, Message: "角色菜单绑定失败"}

	// 操作日志相关
	LogCreateFailed = Response{Code: 51032, Message: "创建操作日志失败"}
	LogQueryFailed  = Response{Code: 51033, Message: "查询操作日志失败"}
)
