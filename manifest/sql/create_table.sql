# create database `simple_tiktok`;

# use `simple_tiktok`;

-- -------------------------
-- Table structure for user
-- -------------------------
drop table if exists `user`;
create table `user` (
    `id` bigint not null primary key comment '用户ID, 雪花算法生成', 
    `name` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci not null default '' comment '用户名称, 可自定义',
    `password` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci not null default '' comment '密码, 登录使用',
    `avatar` varchar(128) not null default '' comment '用户头像',
    `background_image` varchar(128) not null default '' comment '用户个人页面顶部大图',
    `signature` varchar(1024) not null default '' comment '用户个人简介',
    `total_favorited` bigint not null default 0 comment '用户获赞总数',
    `favorite_count` bigint not null default 0 comment '用户点赞总数',
    `work_count` bigint not null default 0 comment '用户作品总数',
    `follow_count` bigint not null default 0 comment '关注总数',
    `follower_count` bigint not null default 0 comment '粉丝总数',
    `is_follow` tinyint(1) not null default 0 comment '是否关注, 0-未关注 1-已关注',
    `created_at` timestamp not null default CURRENT_TIMESTAMP comment '创建时间',
    `updated_at` timestamp not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP comment '更新时间',
    `deleted_at` timestamp default null comment '删除时间',
    UNIQUE `idx_name` (`name`) comment '用户名称唯一索引'
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- -------------------------
-- Table structure for video
-- -------------------------
drop table if exists `video`;
create table `video` (
    `id` bigint not null primary key comment '视频ID, 雪花算法生成',
    `title` varchar(64) not null default '' comment '视频标题',
    `play_url` varchar(128) not null default '' comment '视频播放链接',
    `cover_url` varchar(128) not null default '' comment '视频封面链接',
    `author_id` bigint not null comment '视频作者ID, 哪个用户发布的视频',
    `favorite_count` bigint not null default 0 comment '视频点赞数量',
    `comment_count` bigint not null default 0 comment '视频评论数量',
    `is_favorite` tinyint not null default 0 comment '是否点赞该视频, 0-未点赞 1-已点赞',
    `created_at` timestamp not null default CURRENT_TIMESTAMP comment '创建时间',
    `updated_at` timestamp not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP comment '修改时间',
    `deleted_at` timestamp default null comment '删除时间',
    index `idx_title` (`title`) comment '视频标题索引'
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- -------------------------
-- Table structure for favorite
-- -------------------------
drop table if exists `favorite`;
create table `favorite` (
    `id` bigint not null primary key comment '点赞ID, 雪花算法生成',
    `user_id` bigint not null comment '用户ID, 是谁点了赞',
    `video_id` bigint not null comment '点了哪个视频的赞',
    `created_at` timestamp not null default CURRENT_TIMESTAMP comment '创建时间',
    `updated_at` timestamp not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP comment '修改时间',
    `deleted_at` timestamp default null comment '删除时间'
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for comment
-- ----------------------------
drop table if exists `comment`;
create table `comment` (
    `id` bigint not null primary key comment '评论ID, 雪花算法生成',
    `content` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci not null default '' comment '评论内容',
    `author_id` bigint not null comment '评论者ID, 是谁评论的',
    `video_id` bigint not null comment '评论所属的视频ID',
    `created_at` timestamp not null default CURRENT_TIMESTAMP comment '创建时间',
    `updated_at` timestamp not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP comment '修改时间',
    `deleted_at` timestamp default null comment '删除时间'
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- --------------------------
-- Table structure for follow
-- --------------------------
drop table if exists `follow`;
create table `follow` (
    `id` bigint not null primary key comment '关注ID, 雪花算法生成',
    `user_id` bigint not null comment '被关注人的ID',
    `follower_id` bigint not null comment '关注者的ID',
    `created_at` timestamp not null default CURRENT_TIMESTAMP comment '创建时间',
    `updated_at` timestamp not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP comment '修改时间',
    `deleted_at` timestamp default null comment '删除时间',
    index `idx_follower_id` (`follower_id`)
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- --------------------------
-- Table structure for chat
-- --------------------------
drop table if exists `chat`;
create table `chat` (
    `id` bigint not null primary key comment '消息ID, 雪花算法生成',
    `user_id` bigint not null comment '发消息的用户ID',
    `to_user_id` bigint not null comment '接收消息的用户ID',
    `content` varchar(256) not null default '' comment '消息内容',
    `created_at` timestamp not null default CURRENT_TIMESTAMP comment '创建时间',
    `updated_at` timestamp not null default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP comment '修改时间',
    `deleted_at` timestamp default null comment '删除时间'
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;