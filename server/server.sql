create table user
(
    username  varchar(20) null,
    id        int auto_increment
        primary key,
    user_mail varchar(20) not null
)
    comment '用户';