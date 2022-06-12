create table user
(
    username  varchar(20) null,
    id        int auto_increment
        primary key,
    user_mail varchar(20) not null
    win_cout int default 0
)
    comment '用户';

