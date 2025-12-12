-- `t_menu_widget_api` 菜单、控件、API表 起始数据；起始菜单数据只有超级管理员设置密码后能看到
INSERT INTO `t_menu_widget_api` (`type`, `identifier`, `name`, `intro`, `icon`) 
VALUES (1, "/user", "user menu", "用户页", ""), -- 用户相关菜单

(1, "/domain", "domain menu", "域租户页", ""), -- 域租户相关菜单

(1, "/role", "role menu", "角色页", ""), -- 角色相关菜单

(1, "/menu", "frontend menu", "前端菜单页", ""), -- 菜单相关菜单

(1, "/widget", "frontend widget", "前端控件页", ""), -- 控件相关菜单

(1, "/api", "backend api", "后端API页", ""), -- API相关菜单

(2, "btn_add_user", "add user button", "新增用户按钮", ""),  -- 用户相关控件
(2, "btn_update_user", "update user button", "更新用户按钮", ""),
(2, "btn_delete_user", "delete user button", "删除用户按钮", ""),
(2, "btn_search_user", "search user button", "搜索用户按钮", ""),

(2, "btn_add_domain", "add domain button", "新增域租户按钮", ""),    -- 域租户相关控件
(2, "btn_update_domain", "update domain button", "更新域租户按钮", ""),
(2, "btn_delete_domain", "delete domain button", "删除域租户按钮", ""),
(2, "btn_search_domain", "search domain button", "搜索域租户按钮", ""),

(2, "btn_add_role", "add role button", "新增角色按钮", ""),    -- 角色相关控件
(2, "btn_update_role", "update role button", "更新角色按钮", ""),
(2, "btn_delete_role", "delete role button", "删除角色按钮", ""),
(2, "btn_search_role", "search role button", "搜索角色按钮", ""),

(2, "btn_add_menu", "add menu button", "新增菜单按钮", ""),    -- 菜单相关控件
(2, "btn_update_menu", "update menu button", "更新菜单按钮", ""),
(2, "btn_delete_menu", "delete menu button", "删除菜单按钮", ""),
(2, "btn_search_menu", "search menu button", "搜索菜单按钮", ""),

(2, "btn_add_widget", "add widget button", "新增控件按钮", ""),    -- 控件相关控件
(2, "btn_update_widget", "update widget button", "更新控件按钮", ""),
(2, "btn_delete_widget", "delete widget button", "删除控件按钮", ""),
(2, "btn_search_widget", "search widget button", "搜索控件按钮", ""),

(2, "btn_add_api", "add api button", "新增API按钮", ""),    -- API相关控件
(2, "btn_update_api", "update api button", "更新API按钮", ""),
(2, "btn_delete_api", "delete api button", "删除API按钮", ""),
(2, "btn_search_api", "search api button", "搜索API按钮", ""),

(3, "GET /api/v1/ping", "ping", "PING接口", "域租户菜单"),  -- ping 接口

(3, "POST /api/v1/user/visit", "user visit api", "用户首次访问时传递用户代理等信息", ""), -- 用户相关API
(3, "POST /api/v1/user/sign-in", "user sign-in api", "用户登录", ""),
(3, "POST /api/v1/user/sign-out", "user sign-out api", "用户登出", ""),
(3, "POST /api/v1/user", "add user api", "新增用户", ""),
(3, "GET /api/v1/user", "get user list api", "获取用户列表", ""),
(3, "GET /api/v1/user/:xid", "get user api", "获取用户信息", ""),
(3, "PUT /api/v1/user/:xid", "update user api", "更新用户信息", ""),
(3, "DELETE /api/v1/user/:xid", "delete user api", "删除用户", ""),

(3, "POST /api/v1/domain", "add domain api", "新增域租户", ""),  -- 域租户相关API
(3, "GET /api/v1/domain", "get domain list api", "获取域租户列表", ""),
(3, "GET /api/v1/domain/:xid", "get domain api", "获取域租户信息", ""),
(3, "PUT /api/v1/domain/:xid", "update domain api", "更新域租户信息", ""),
(3, "DELETE /api/v1/domain/:xid", "delete domain api", "删除域租户", ""),

(3, "POST /api/v1/role", "add role api", "新增角色", ""),  -- 角色相关API
(3, "GET /api/v1/role", "get role list api", "获取角色列表", ""),
(3, "GET /api/v1/role/:xid", "get role api", "获取角色信息", ""),
(3, "PUT /api/v1/role/:xid", "update role api", "更新角色信息", ""),
(3, "DELETE /api/v1/role/:xid", "delete role api", "删除角色", ""),

(3, "POST /api/v1/menu", "add menu api", "新增菜单", ""),  -- 菜单相关API
(3, "GET /api/v1/menu", "get menu list api", "获取菜单列表", ""),
(3, "GET /api/v1/menu/:xid", "get menu api", "获取菜单信息", ""),
(3, "PUT /api/v1/menu/:xid", "update menu api", "更新菜单信息", ""),
(3, "DELETE /api/v1/menu/:xid", "delete menu api", "删除菜单", ""),

(3, "POST /api/v1/widget", "add widget api", "新增控件", ""),  -- 控件相关API
(3, "GET /api/v1/widget", "get widget list api", "获取控件列表", ""),
(3, "GET /api/v1/widget/:xid", "get widget api", "获取控件信息", ""),
(3, "PUT /api/v1/widget/:xid", "update widget api", "更新控件信息", ""),
(3, "DELETE /api/v1/widget/:xid", "delete widget api", "删除控件", ""),

(3, "POST /api/v1/api", "add api api", "新增API", ""),  -- API相关API
(3, "GET /api/v1/api", "get api list api", "获取API列表", ""),
(3, "GET /api/v1/api/:xid", "get api api", "获取API信息", ""),
(3, "PUT /api/v1/api/:xid", "update api api", "更新API信息", ""),
(3, "DELETE /api/v1/api/:xid", "delete api api", "删除API", "");