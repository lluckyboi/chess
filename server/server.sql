create table action
(
    action_id  int auto_increment
        primary key,
    history_id int          not null,
    data       varchar(600) null comment '棋盘数据'
);

create table history
(
    id      int auto_increment
        primary key,
    user_id int not null
);

create table user
(
    username  varchar(20)   null,
    id        int auto_increment
        primary key,
    user_mail varchar(20)   not null,
    win_count int default 0 null
)
    comment '用户';