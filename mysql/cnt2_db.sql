/*
Navicat MySQL Data Transfer

Source Server         : 配置中心库
Source Server Version : 50619
Source Database       : cnt2_db

Target Server Type    : MYSQL
Target Server Version : 50619
File Encoding         : 65001

Date: 2017-12-08 11:51:28
*/
CREATE DATABASE IF NOT EXISTS cnt2_db DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_general_ci;

USE cnt2_db;

SET FOREIGN_KEY_CHECKS=0;

CREATE TABLE `user` (
  `uid` bigint(20) NOT NULL AUTO_INCREMENT,
  `username` varchar(40) DEFAULT NULL,
  `pwd` varchar(64) DEFAULT NULL,
  `create_time` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`uid`)
) ENGINE=InnoDB AUTO_INCREMENT=10005 DEFAULT CHARSET=utf8mb4;

INSERT INTO `user` (`uid`, `username`, `pwd`, `create_time`) VALUES ('10000', 'admin', 'e10adc3949ba59abbe56e057f20f883e', '2018-01-04 17:02:03');

-- ----------------------------
-- Table structure for app
-- ----------------------------
DROP TABLE IF EXISTS `app`;
CREATE TABLE `app` (
  `app` varchar(50) NOT NULL COMMENT '业务标示',
  `app_type` tinyint(1) DEFAULT '0' COMMENT '业务类型：0-服务端；9-app客户端',
  `name` varchar(100) DEFAULT NULL COMMENT '业务名称',
  `charger` varchar(50) DEFAULT NULL COMMENT '负责人',
  `charger_uid` bigint(20) DEFAULT NULL COMMENT '负责人UID',
  `create_time` datetime DEFAULT NULL,
  PRIMARY KEY (`app`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for config
-- ----------------------------
DROP TABLE IF EXISTS `config`;
CREATE TABLE `config` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `app` varchar(50) NOT NULL COMMENT '业务标示',
  `profile` varchar(50) NOT NULL COMMENT '环境名',
  `key` varchar(200) NOT NULL COMMENT '配置KEY',
  `value` longtext CHARACTER SET utf8mb4 COMMENT '配置信息',
  `version` bigint(20) NOT NULL COMMENT '版本号',
  `published_value` longtext CHARACTER SET utf8mb4 COMMENT '最近已发布配置，进程启动时全量拉取该配置',
  `published_version` bigint(20) DEFAULT NULL COMMENT '最近已发布版本号，进程启动时全量拉取该配置',
  `validator` varchar(1000) DEFAULT NULL COMMENT '验证配置的js脚本，前端执行，可空',
  `modifier` varchar(50) DEFAULT NULL COMMENT '修改人',
  `modify_time` datetime DEFAULT NULL COMMENT '修改时间',
  `description` varchar(300) DEFAULT NULL COMMENT '修改描述',
  `approve_type` tinyint(1) DEFAULT '1' COMMENT '审批类型：1-未审核；2-已审核',
  `approver` varchar(50) DEFAULT NULL COMMENT '审核人',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_config` (`app`,`profile`,`key`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for config_history
-- ----------------------------
DROP TABLE IF EXISTS `config_history`;
CREATE TABLE `config_history` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '序号（或作为版本号）',
  `app` varchar(50) DEFAULT NULL COMMENT '业务标示',
  `profile` varchar(50) DEFAULT NULL COMMENT '环境名',
  `key` varchar(200) DEFAULT NULL COMMENT '配置KEY',
  `pre_value` longtext CHARACTER SET utf8mb4 COMMENT '上一配置信息',
  `cur_value` longtext CHARACTER SET utf8mb4 COMMENT '当前配置信息',
  `pre_version` bigint(20) DEFAULT NULL COMMENT '上一版本号',
  `cur_version` bigint(20) DEFAULT NULL COMMENT '当前版本号',
  `validator` varchar(1000) DEFAULT NULL COMMENT '验证配置的js脚本，前端执行，可空',
  `modifier` varchar(50) DEFAULT NULL COMMENT '修改人',
  `modify_time` datetime DEFAULT NULL COMMENT '修改时间',
  `description` varchar(300) DEFAULT NULL COMMENT '修改描述',
  `operate_type` int(11) DEFAULT NULL COMMENT '发布类型，新增，删除，修改，回滚',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_config_his` (`app`,`profile`,`key`,`cur_version`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for node
-- ----------------------------
DROP TABLE IF EXISTS `node`;
CREATE TABLE `node` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'node id',
  `app` varchar(50) NOT NULL COMMENT '业务标示',
  `profile` varchar(50) DEFAULT NULL COMMENT '环境名',
  `sip` varchar(50) DEFAULT NULL COMMENT '环境名',
  `pid` int(11) DEFAULT '0' COMMENT 'node pid',
  `register_time` datetime DEFAULT NULL COMMENT '注册时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2148 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for node_publish
-- ----------------------------
DROP TABLE IF EXISTS `node_publish`;
CREATE TABLE `node_publish` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `node_id` bigint(20) NOT NULL COMMENT '灰度发布时记录nodeid',
  `app` varchar(50) NOT NULL COMMENT '业务标示',
  `profile` varchar(50) DEFAULT NULL COMMENT '环境名',
  `key` varchar(200) DEFAULT NULL COMMENT '配置KEY',
  `version` bigint(20) DEFAULT NULL COMMENT '版本号',
  `publish_time` datetime DEFAULT NULL COMMENT '注册时间',
  `publish_result` tinyint(1) DEFAULT NULL COMMENT '发布结果：0-未发布；1-成功；2-失败',
  `publish_type` tinyint(1) DEFAULT NULL COMMENT '发布类型：0-灰度；1-全量',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_unique` (`node_id`,`app`,`profile`,`key`,`version`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for profile
-- ----------------------------
DROP TABLE IF EXISTS `profile`;
CREATE TABLE `profile` (
  `app` varchar(50) CHARACTER SET utf8 NOT NULL,
  `profile` varchar(50) NOT NULL,
  `name` varchar(100) DEFAULT NULL,
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`app`,`profile`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for user_auth
-- ----------------------------
DROP TABLE IF EXISTS `user_auth`;
CREATE TABLE `user_auth` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `uid` bigint(20) DEFAULT NULL COMMENT 'yyuid',
  `uname` varchar(50) DEFAULT NULL,
  `app` varchar(50) DEFAULT NULL COMMENT '业务标示',
  `permission` tinyint(1) DEFAULT '1' COMMENT '权限：1-开发；9-管理',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_uid_app` (`uid`,`app`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
