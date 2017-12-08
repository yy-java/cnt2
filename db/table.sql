CREATE TABLE `app` (
  `app` varchar(50) NOT NULL COMMENT '业务标示',
  `app_type` TINYINT(1) DEFAULT 0 COMMENT '业务类型：0-服务端；9-app客户端',
  `name` varchar(100) DEFAULT NULL COMMENT '业务名称',
  `charger` varchar(50) DEFAULT NULL COMMENT '负责人',
  `create_time` datetime DEFAULT NULL,
  PRIMARY KEY (`app`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
CREATE TABLE `user_auth` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `uid` bigint(20) DEFAULT NULL COMMENT 'yyuid',
  `uname` varchar(50) DEFAULT NULL COMMENT '',
  `app` varchar(50) DEFAULT NULL COMMENT '业务标示',
  `permission` TINYINT(1) DEFAULT 1 COMMENT '权限：1-开发；9-管理',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
CREATE TABLE `config` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `app` varchar(50) NOT NULL COMMENT '业务标示',
  `profile` varchar(50) NOT NULL COMMENT '环境名',
  `key` varchar(200) NOT NULL COMMENT '配置KEY',
  `value` varchar(2000) CHARACTER SET utf8mb4 COMMENT '配置信息',
  `version` bigint(20) NOT NULL COMMENT '版本号',
  `published_value` varchar(2000) CHARACTER SET utf8mb4 COMMENT '最近已发布配置，进程启动时全量拉取该配置',
  `published_version` bigint(20) DEFAULT NULL COMMENT '最近已发布版本号，进程启动时全量拉取该配置',
  `validator` varchar(1000) DEFAULT NULL COMMENT '验证配置的js脚本，前端执行，可空',
  `modifier` varchar(50) DEFAULT NULL COMMENT '修改人',
  `modify_time` datetime DEFAULT NULL COMMENT '修改时间',
  `description` varchar(300) DEFAULT NULL COMMENT '修改描述',
  `approve_type` TINYINT(1) DEFAULT 0 COMMENT '审批类型：0-未审核；1-已审核',
  `approver` varchar(50) DEFAULT NULL COMMENT '审核人',	
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_config` (`app`, `profile`, `key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
CREATE TABLE `config_history` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '序号（或作为版本号）',
  `app` varchar(50) DEFAULT NULL COMMENT '业务标示',
  `profile` varchar(50) DEFAULT NULL COMMENT '环境名',
  `key` varchar(200) DEFAULT NULL COMMENT '配置KEY',
  `pre_value` varchar(2000) CHARACTER SET utf8mb4 COMMENT '上一配置信息',
  `cur_value` varchar(2000) CHARACTER SET utf8mb4 COMMENT '当前配置信息',
  `pre_version` bigint(20) DEFAULT NULL COMMENT '上一版本号',
  `cur_version` bigint(20) DEFAULT NULL COMMENT '当前版本号',
  `validator` varchar(1000) DEFAULT NULL COMMENT '验证配置的js脚本，前端执行，可空',
  `modifier` varchar(50) DEFAULT NULL COMMENT '修改人',
  `modify_time` datetime DEFAULT NULL COMMENT '修改时间',
  `description` varchar(300) DEFAULT NULL COMMENT '修改描述',
  `operate_type` int(11) DEFAULT NULL COMMENT '发布类型，新增，删除，修改，回滚',
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_config_his` (`app`, `profile`, `key`, `cur_version`)  
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
CREATE TABLE `node` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'node id',
  `app` varchar(50) NOT NULL COMMENT '业务标示',
  `profile` varchar(50) DEFAULT NULL COMMENT '环境名',
  `sip` varchar(50) DEFAULT NULL COMMENT '环境名',
  `pid` int(11) DEFAULT 0 COMMENT 'node pid',
  `register_time` datetime DEFAULT NULL COMMENT '注册时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
CREATE TABLE `node_publish` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '',
  `node_id` bigint(20) NOT NULL COMMENT '灰度发布时记录nodeid',
  `app` varchar(50) NOT NULL COMMENT '业务标示',
  `profile` varchar(50) DEFAULT NULL COMMENT '环境名',
  `key` varchar(200) DEFAULT NULL COMMENT '配置KEY',
  `version` bigint(20) DEFAULT NULL COMMENT '版本号',
  `publish_time` datetime DEFAULT NULL COMMENT '注册时间',
  `publish_result` TINYINT(1) DEFAULT NULL COMMENT '发布结果：0-已发布；1-成功；2-失败',
  `publish_type` TINYINT(1) DEFAULT NULL COMMENT '发布类型：0-灰度；1-全量',
  PRIMARY KEY (`id`)  
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
CREATE TABLE `profile` (
  `app` varchar(50) NOT NULL,
  `profile` varchar(50) NOT NULL,
  `name` varchar(100) DEFAULT NULL,
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`app`,`profile`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
