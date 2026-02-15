-- `t_user_agent` 用户代理表

CREATE TABLE `t_user_agent` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` bigint NULL,
  `user_xid` varchar(64) NULL,
  `ip` varchar(64) NOT NULL,
  `ua` text NOT NULL,
  `trace_id` varchar(64) NOT NULL,
  `token` text NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- `t_user` 用户表

CREATE TABLE `t_user` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` bigint NULL,
  `xid` varchar(64) NOT NULL UNIQUE,
  `email` varchar(64) NOT NULL UNIQUE,
  `nick_name` varchar(32) NULL,
  `real_name` varchar(32) NULL,
  `password` varchar(64) DEFAULT NULL,
  `avatar` varchar(256) NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT uk_nick_name UNIQUE (nick_name, deleted_at)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- `t_role` 角色表

CREATE TABLE `t_role` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` bigint NULL,
  `xid` varchar(64) NOT NULL UNIQUE,
  `name` varchar(64) NOT NULL,
  `intro` varchar(256) DEFAULT NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT uk_name UNIQUE (name, deleted_at)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- `t_menu_widget_api` 菜单、控件、API表

CREATE TABLE `t_menu_widget_api` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` bigint NULL,
  `type` tinyint NOT NULL, -- 1: 菜单；2: 控件；3: api
  `identifier` varchar(256) NOT NULL UNIQUE,
  `name` varchar(64) NOT NULL,
  `intro` varchar(256) DEFAULT NULL,
  `icon` varchar(256) DEFAULT NULL, -- 菜单图标url
  PRIMARY KEY (`id`),
  CONSTRAINT uk_name UNIQUE (name, deleted_at)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- `t_domain` 域租户表

CREATE TABLE `t_domain` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` bigint NULL,
  `xid` varchar(64) NOT NULL UNIQUE,
  `name` varchar(64) NOT NULL,
  `intro` varchar(256) DEFAULT NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT uk_name UNIQUE (name, deleted_at)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- `t_casbin_rule` casbin规则表
-- 主要是用casbin的 `RBAC with Domains` 模型
-- ptype -> p, g
-- 当ptype为g时：v0->user xid, v1->role xid, v2->domain xid
-- 当ptype为p时：v0->role xid, v1->domain xid, v2->resource identifier, v3->method
--      当resource为非API（菜单、控件）时，v3存储空字符串
--      当resource为API时，v3存储相应的GET、POST、DELETE方法

CREATE TABLE `t_casbin_rule` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` bigint NULL,
  `ptype` varchar(512) DEFAULT NULL,
  `v0` varchar(512) DEFAULT NULL,
  `v1` varchar(512) DEFAULT NULL,
  `v2` varchar(512) DEFAULT NULL,
  `v3` varchar(512) DEFAULT NULL,
  `v4` varchar(512) DEFAULT NULL,
  `v5` varchar(512) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
