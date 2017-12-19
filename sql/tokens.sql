/*
 Navicat Premium Data Transfer

 Source Server         : datica-users
 Source Server Type    : PostgreSQL
 Source Server Version : 90503
 Source Host           : localhost
 Source Database       : datica_users_dev
 Source Schema         : public

 Target Server Type    : PostgreSQL
 Target Server Version : 90503
 File Encoding         : utf-8

 Date: 12/19/2017 15:39:52 PM
*/

-- ----------------------------
--  Table structure for tokens
-- ----------------------------
CREATE TABLE "public"."tokens" (
	"token_id" uuid NOT NULL,
	"token" varchar(255) NOT NULL COLLATE "default",
	"email" varchar(75) NOT NULL COLLATE "default"
)
WITH (OIDS=FALSE);
ALTER TABLE "public"."tokens" OWNER TO "admin";

-- ----------------------------
--  Primary key structure for table tokens
-- ----------------------------
ALTER TABLE "public"."tokens" ADD PRIMARY KEY ("token_id") NOT DEFERRABLE INITIALLY IMMEDIATE;

-- ----------------------------
--  Indexes structure for table tokens
-- ----------------------------
CREATE UNIQUE INDEX  "token_idx" ON "public"."tokens" USING btree(token COLLATE "default" "pg_catalog"."text_ops" ASC NULLS LAST);

