# CREATE DATABASE if not exists 'aweme_community';
# DROP TABLE if exists 'user';
create table user
(
    id             bigint unsigned auto_increment comment '用户id' primary key,
    user_id        varchar(256)                not null comment 'uuid generate',
    name           varchar(32)     default ''  not null comment '用户名称',
    password       varchar(32)     default ''  not null comment '密码',
    token          varchar(512)                not null comment '鉴权',
    follow_count   bigint unsigned default '0' not null comment '关注总数',
    follower_count bigint unsigned default '0' not null comment '粉丝总数',
    created_at     timestamp                   null comment '创建时间',
    updated_at     timestamp                   null comment '更新时间',
    deleted_at     timestamp                   null comment '删除时间'
)
    comment '用户表' engine = InnoDB;
INSERT INTO `user`(id, name, password, follow_count, follower_count, is_follow)
VALUES (1, 'Jerry', 'jerry', 12, 23, 1),
       (2, 'Tom', 'Tom', 1, 0, 0),
       (3, 'Amy', 'Amy', 0, 1, 1);

DROP TABLE IF EXISTS `video`;
create table video
(
    id              bigint unsigned auto_increment comment '视频id' primary key,
    video_id        varchar(256)                               not null comment 'video ID',
    user_id         varchar(256)                               not null comment '作者',
    title           varchar(128) default ''                    not null comment '标题',
    play_url        varchar(512) default ''                    not null comment '播放地址',
    cover_url       varchar(512) default ''                    not null comment '封面地址',
    favourite_count bigint       default 0                     not null comment '点赞总数',
    comment_count   bigint       default 0                     not null comment '评论总数',
    created_at      timestamp    default '0000-00-00 00:00:00' not null comment '创建时间',
    updated_at      timestamp                                  null comment '更新时间',
    deleted_at      timestamp                                  null comment '删除时间'
)
    comment '视频表' engine = InnoDB;

insert into `aweme_community`.`video`
values (1, 1, 'ckx', 'https://www.bilibili.com/video/BV1Ve4y147D2?t=4.7',
        'https://c-ssl.duitang.com/uploads/item/202006/13/20200613202923_flfxg.jpg', 3, 2, CURRENT_TIME);

DROP TABLE IF EXISTS `favorite`;
create table favourite
(
    id        bigint unsigned auto_increment comment '自增id' primary key,
    favour_id varchar(256) default '' not null,
    user_id   varchar(265)            not null comment '用户ID',
    video_id  varchar(256)            not null comment '视频id'
)
    comment '点赞表' engine = InnoDB;

insert into  `aweme_community`.`favorite`
values (1, 1, 1) (2, 2, 1), (3, 3, 1);

-- auto-generated definition
create table comment
(
    id         bigint unsigned auto_increment comment '自增id' primary key,
    comment_id varchar(256)                           not null comment '评论ID',
    user_id    varchar(256) default '0'               not null comment '用户id',
    video_id   varchar(256) default '0'               not null comment '视频id',
    content    text                                   not null comment '评论内容',
    created_at timestamp    default CURRENT_TIMESTAMP null comment '创建时间',
    updated_at timestamp                              null comment '更新时间',
    deleted_at timestamp                              null comment '删除时间'
)
    comment '评论表' engine = InnoDB;


INSERT INTO `aweme_community`.`comment`(id, user_id, video_id, contents, create_date)
VALUES (1, 1, 1, '这视频也太模糊了！', CURRENT_TIMESTAMP),
       (2, 2, 1, '老CCTV了～', CURRENT_TIMESTAMP);
-- auto-generated definition
create table relation
(
    id          bigint unsigned auto_increment comment '关注id' primary key,
    relation_id varchar(256)             not null comment 'id',
    user_id     varchar(256) default '0' not null comment '用户id',
    to_user_id  varchar(256) default '0' not null comment '对方用户id',
    status      tinyint      default -2  not null comment '关注状态'
)
    comment '关注表' engine = InnoDB;

