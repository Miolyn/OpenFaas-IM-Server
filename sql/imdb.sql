/*
 Navicat MySQL Data Transfer

 Source Server         : Finders
 Source Server Type    : MySQL
 Source Server Version : 80021
 Source Host           : 123.56.104.212:3306
 Source Schema         : finders_imdb

 Target Server Type    : MySQL
 Target Server Version : 80021
 File Encoding         : 65001

 Date: 31/08/2020 19:22:13
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for group
-- ----------------------------
DROP TABLE IF EXISTS `group`;
CREATE TABLE `group` (
  `group_id` bigint NOT NULL AUTO_INCREMENT COMMENT '群组ID',
  `name` varchar(255) COLLATE utf8mb4_general_ci NOT NULL COMMENT '名称',
  `introduction` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '群组介绍',
  `type` int NOT NULL COMMENT '群组类型: 1.系统消息群 2.用户自建群',
  `user_num` int(10) unsigned zerofill NOT NULL COMMENT '用户数量',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`group_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Table structure for message
-- ----------------------------
DROP TABLE IF EXISTS `message`;
CREATE TABLE `message` (
  `message_id` bigint NOT NULL AUTO_INCREMENT COMMENT '消息ID',
  `device_id` bigint DEFAULT NULL COMMENT '设备唯一标示',
  `from_uid` varchar(64) COLLATE utf8mb4_general_ci NOT NULL COMMENT '消息发送人user_id',
  `to_id` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '消息接受者id',
--   `to_uid` varchar(64) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '消息接受者user_id',
--   `to_gid` bigint DEFAULT NULL COMMENT '群组ID',
  `receiver_type` int NOT NULL COMMENT '接受者类型 1.个人 2.群组',
  `type` int NOT NULL COMMENT '消息类型：1.文本 2.语音 3.视频 4.表情 6.图片',
  `status` int NOT NULL COMMENT '消息状态：1.未发送   2.已送达 (1的话会检查还有哪些用户没有收到)',
  `content` text COLLATE utf8mb4_general_ci NOT NULL COMMENT '消息内容，根据消息类型有所不同，json串',
  `created_at` datetime NOT NULL COMMENT '创建时间',
  `updated_at` datetime NOT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`message_id`),
  KEY `from_uid_idx` (`from_uid`) USING BTREE,
  KEY `to_uid_idx` (`to_uid`) USING BTREE,
  KEY `to_gid_idx` (`to_gid`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


-- ----------------------------
-- Table structure for timeline
-- ----------------------------
DROP TABLE IF EXISTS `timeline`;
CREATE TABLE `timeline` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `timeline_id` bigint DEFAULT NULL COMMENT 'timeline ID',
  `seq` bigint DEFAULT NULL COMMENT '消息序列，一个timeline下的seq从0递增。',
  `message_id` bigint DEFAULT NULL COMMENT '消息ID',
  `type` int DEFAULT NULL COMMENT 'timeline类型，1.个人 2.单聊 3.群聊',
  `status` int DEFAULT NULL COMMENT 'timeline状态，1.未读 2.已读',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `timeline_idx` (`timeline_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Table structure for timeline_map
-- ----------------------------
DROP TABLE IF EXISTS `timeline_map`;
CREATE TABLE `timeline_map` (
  `timeline_id` bigint NOT NULL AUTO_INCREMENT COMMENT 'timeline 主键',
  `type` int NOT NULL COMMENT 'timeline类型：1.个人 2.单聊 3.群聊',
  `object_ids` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'timeline对象id: type=1 时为user_id，type=2时为user_id,user_id，type=3时为group_id',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`timeline_id`),
  KEY `object_ids_idx` (`object_ids`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
                         `user_id` varchar(50) NOT NULL,
                         `password` varchar(100) DEFAULT NULL,
                         `status` int DEFAULT NULL,
                         `username` varchar(50) DEFAULT NULL,
                        `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                          `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                          `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
                         PRIMARY KEY (`user_id`),
                         UNIQUE KEY `unique_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

SET FOREIGN_KEY_CHECKS = 1;


