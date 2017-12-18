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

 Date: 12/18/2017 07:57:30 AM
*/

-- ----------------------------
--  Table structure for users
-- ----------------------------
CREATE TABLE "public"."users" (
	"user_id" uuid NOT NULL,
	"name" varchar(75) COLLATE "default",
	"email" varchar(75) NOT NULL COLLATE "default",
	"salt" varchar(64) NOT NULL COLLATE "default",
	"password" varchar(128) NOT NULL COLLATE "default",
	"created_at" timestamp(6) NOT NULL,
	"updated_at" timestamp(6) NOT NULL
)
WITH (OIDS=FALSE);
ALTER TABLE "public"."users" OWNER TO "admin";

-- ----------------------------
--  Primary key structure for table users
-- ----------------------------
ALTER TABLE "public"."users" ADD PRIMARY KEY ("user_id") NOT DEFERRABLE INITIALLY IMMEDIATE;

-- ----------------------------
--  Uniques structure for table users
-- ----------------------------
ALTER TABLE "public"."users" ADD CONSTRAINT "email_uniq" UNIQUE ("email") NOT DEFERRABLE INITIALLY IMMEDIATE;

-- ----------------------------
--  Indexes structure for table users
-- ----------------------------
CREATE UNIQUE INDEX  "email_idx" ON "public"."users" USING btree(email COLLATE "default" "pg_catalog"."text_ops" ASC NULLS LAST);

