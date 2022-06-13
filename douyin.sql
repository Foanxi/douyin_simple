/*
 Navicat Premium Data Transfer

 Source Server         : Demo
 Source Server Type    : MySQL
 Source Server Version : 80026
 Source Host           : localhost:3306
 Source Schema         : douyin

 Target Server Type    : MySQL
 Target Server Version : 80026
 File Encoding         : 65001

 Date: 13/06/2022 15:26:42
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for author_fans
-- ----------------------------
DROP TABLE IF EXISTS `author_fans`;
CREATE TABLE `author_fans`  (
  `author_id` bigint NOT NULL,
  `favourite_id` bigint NOT NULL,
  PRIMARY KEY (`author_id`, `favourite_id`) USING BTREE,
  INDEX `author_fans_ibfk_2`(`favourite_id`) USING BTREE,
  CONSTRAINT `author_fans_ibfk_1` FOREIGN KEY (`author_id`) REFERENCES `user` (`Id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `author_fans_ibfk_2` FOREIGN KEY (`favourite_id`) REFERENCES `user` (`Id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for comment
-- ----------------------------
DROP TABLE IF EXISTS `comment`;
CREATE TABLE `comment`  (
  `comment_id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` bigint NOT NULL,
  `video_id` bigint NOT NULL,
  `comment_text` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `create_time` char(6) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  PRIMARY KEY (`comment_id`, `user_id`, `video_id`) USING BTREE,
  INDEX `user_id`(`user_id`) USING BTREE,
  INDEX `video_id`(`video_id`) USING BTREE,
  CONSTRAINT `comment_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `user` (`Id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `comment_ibfk_2` FOREIGN KEY (`video_id`) REFERENCES `video` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB AUTO_INCREMENT = 14 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for favouriter_video
-- ----------------------------
DROP TABLE IF EXISTS `favouriter_video`;
CREATE TABLE `favouriter_video`  (
  `user_id` bigint NOT NULL,
  `video_id` bigint NOT NULL,
  `favourite` tinyint(1) NULL DEFAULT NULL,
  PRIMARY KEY (`user_id`, `video_id`) USING BTREE,
  INDEX `user_video_ibfk_2`(`video_id`) USING BTREE,
  CONSTRAINT `favouriter_video_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `user` (`Id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `favouriter_video_ibfk_2` FOREIGN KEY (`video_id`) REFERENCES `video` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`  (
  `Id` bigint NOT NULL AUTO_INCREMENT,
  `Name` varchar(32) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `Password` varchar(32) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `FollowCount` bigint NULL DEFAULT NULL,
  `FollowerCount` bigint NULL DEFAULT NULL,
  PRIMARY KEY (`Id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 13 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for video
-- ----------------------------
DROP TABLE IF EXISTS `video`;
CREATE TABLE `video`  (
  `id` bigint NOT NULL,
  `author_id` bigint NOT NULL,
  `play_url` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `cover_url` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `favourite_count` bigint NOT NULL,
  `comment_count` bigint NOT NULL,
  `create_time` datetime NOT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `video_ibfk_1`(`author_id`) USING BTREE,
  CONSTRAINT `video_ibfk_1` FOREIGN KEY (`author_id`) REFERENCES `user` (`Id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
