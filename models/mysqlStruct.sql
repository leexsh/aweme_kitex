# CREATE DATABASE if not exists 'aweme_community';
# DROP TABLE if exists 'user';
CREATE TABLE `user`
(
    `id`          bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增id',
    `user_identity` varchar(256) NOT NULL COMMENT '真实id',
    `name`        varchar(32)        NOT NULL DEFAULT '' COMMENT '用户名称',
    `password`    varchar(32)        NOT NULL DEFAULT '' COMMENT '密码',
    `follow_count` bigint(20) unsigned NOT NULL DEFAULT 0 COMMENT '关注总数',
    `follower_count` bigint(20) unsigned NOT NULL DEFAULT 0 COMMENT '粉丝总数',
    `is_follow` tinyint NOT NULL default 0 COMMENT '是否主播',
    `created_time` timestamp null comment '创建时间',
    `updated_time` timestamp null comment '更新时间',
    `deleted_time` timestamp null comment '删除时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='用户表';

INSERT INTO `user`(id, name, password, follow_count, follower_count, is_follow)
VALUES (1, 'Jerry', 'jerry', 12, 23, 1),
       (2, 'Tom', 'Tom', 1, 0, 0),
       (3, 'Amy', 'Amy', 0, 1, 1);

DROP TABLE IF EXISTS `video`;
CREATE TABLE `video`
(
    `id`          bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '视频id',
    `user_id`     bigint(20) unsigned NOT NULL DEFAULT 0 COMMENT '用户id',
    `title`       varchar(128)        NOT NULL DEFAULT '' COMMENT '标题',
    `play_url`    varchar(128)        NOT NULL DEFAULT '' COMMENT '播放地址',
    `cover_url`   varchar(128)        NOT NULL DEFAULT '' COMMENT '封面地址',
    'favorite_count' bigint(20) unsigned NOT NULL DEFAULT 0 COMMENT '点赞总数',
    'comment_count'  bigint(20) unsigned NOT NULL DEFAULT 0 COMMENT '评论总数',
    'create_date'    timestamp         NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`)
    ) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='视频表';

CREATE TABLE `aweme_community`.`video`  (
    `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '视频id',
    `user_id` bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户id',
    `title` varchar(128) NOT NULL DEFAULT '' COMMENT '标题\r\n',
    `play_url` varchar(512) NOT NULL DEFAULT '' COMMENT '播放地址',
    `cover_url` varchar(512) NOT NULL DEFAULT '' COMMENT '封面地址',
    `favourite_count` bigint(20) NOT NULL default 0 COMMENT '点赞总数',
    `comment_count` bigint(20) NOT NULL default 0 COMMENT '评论总数',
    `create_time` timestamp(0) NOT NULL default 0 COMMENT '创建时间',
    PRIMARY KEY (`id`)
    ) ENGINE = InnoDB
CHARACTER SET = utf8mb4 COMMENT = '视频表' ROW_FORMAT = DEFAULT;

insert into `aweme_community`.`video`
values (1, 1, 'ckx', 'https://www.bilibili.com/video/BV1Ve4y147D2?t=4.7',
        'https://c-ssl.duitang.com/uploads/item/202006/13/20200613202923_flfxg.jpg', 3, 2, CURRENT_TIME);

DROP TABLE IF EXISTS `favorite`;
CREATE TABLE `favorite`
(
    `id`    bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增id',
    `favour_identity`       varchar(256) NOT NULL COMMENT '真实id',
    `user_identity` varchar(256) NOT NULL COMMENT '用户id',
    `video_identity` varchar(256) NOT NULL COMMENT '视频id',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='点赞表';


insert into  `aweme_community`.`favorite`
values (1, 1, 1) (2, 2, 1), (3, 3, 1);


CREATE TABLE `aweme_community`.`comment`  (
    `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '评论id',
    `user_id` bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户id',
    `video_id` bigint(20) UNSIGNED NOT NULL DEFAULT 0 COMMENT '视频id',
    `contents` text NOT NULL COMMENT '评论内容',
    `create_date` timestamp(0) NOT NULL DEFAULT  CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COMMENT = '评论表';

INSERT INTO `aweme_community`.`comment`(id, user_id, video_id, contents, create_date)
VALUES (1, 1, 1, '这视频也太模糊了！', CURRENT_TIMESTAMP),
       (2, 2, 1, '老CCTV了～', CURRENT_TIMESTAMP);

CREATE TABLE `aweme_community`.`relation` (
    `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '关注id',
    `user_id` bigint(20) UNSIGNED NOT NULL default 0 COMMENT '用户id',
    `to_user_id` bigint(20) UNSIGNED NOT NULL default 0 COMMENT '对方用户id',
    `status` tinyint NOT NULL default -1 COMMENT '关注状态',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COMMENT ='关注表';