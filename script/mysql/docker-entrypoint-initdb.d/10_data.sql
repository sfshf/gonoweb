-- docker entrypoint init db 时默认的连接编码为 latin1，导致初始化插入数据编码不是utf8mb4
SET character_set_client = utf8mb4;
SET character_set_connection = utf8mb4;

-- `t_resource` 菜单、控件、API表 起始数据；起始菜单数据只有超级管理员设置密码后能看到
INSERT INTO `t_resource` (`type`, `identifier`, `name`, `intro`, `icon`) 
VALUES (1, "/user", "user menu", "用户页", ""), -- 用户相关菜单

(1, "/domain", "domain menu", "域租户页", ""), -- 域租户相关菜单

(1, "/role", "role menu", "角色页", ""), -- 角色相关菜单

(1, "/resource", "menu/widget/api menu", "资源（菜单/控件/API）页", ""), -- 资源（菜单/控件/API）相关菜单

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

(2, "btn_add_resource", "add resouce button", "新增资源（菜单/控件/API）按钮", ""),    -- 资源（菜单/控件/API）相关控件
(2, "btn_update_resource", "update resouce button", "更新资源（菜单/控件/API）按钮", ""),
(2, "btn_delete_resource", "delete resouce button", "删除资源（菜单/控件/API）按钮", ""),
(2, "btn_search_resource", "search resouce button", "搜索资源（菜单/控件/API）按钮", ""),

(3, "GET /api/v1/ping", "ping", "PING接口", ""),  -- ping 接口

(3, "POST /api/v1/user/visit", "user visit api", "用户首次访问时传递用户代理等信息", ""), -- 用户相关API
(3, "POST /api/v1/user/sign-in", "user sign-in api", "用户登录", ""),
(3, "POST /api/v1/user/sign-out", "user sign-out api", "用户登出", ""),
(3, "PUT /api/v1/user", "user switch role api", "用户切换角色", ""),
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

(3, "POST /api/v1/resource", "add resource api", "新增资源（菜单/控件/API）", ""),  -- 资源（菜单/控件/API）相关API
(3, "GET /api/v1/resource", "get resource list api", "获取资源（菜单/控件/API）列表", ""),
(3, "GET /api/v1/resource/:id", "get resource api", "获取资源（菜单/控件/API）信息", ""),
(3, "PUT /api/v1/resource/:id", "update resource api", "更新资源（菜单/控件/API）信息", ""),
(3, "DELETE /api/v1/resource/:id", "delete resource api", "删除资源（菜单/控件/API）", ""),

(3, "GET /api/v1/casbin/domain/:dxid/role", "list role xids in the domain", "获取域租户下的角色配置", ""),  -- casbin相关API
(3, "GET /api/v1/casbin/domain/:dxid/role/:rxid/resource", "list resource xids of the role in the domain", "获取域租户下角色的资源配置", ""),
(3, "POST /api/v1/casbin/domain/:dxid/role/:rxid/resource", "alloc resources to the role in the domain", "给域租户下的角色分配资源", ""),
(3, "POST /api/v1/casbin/user/:xid", "alloc roles in the domain to the user", "给用户分配域租户下的角色", ""),
(3, "GET /api/v1/casbin/user/:xid/domain", "get domain xids of the user", "获取用户被分配到的域租户列表", ""),
(3, "GET /api/v1/casbin/user/:xid/domain/:dxid/role", "get role xids in the domain of the user", "获取用户在域租户下被分配到的角色列表", "");