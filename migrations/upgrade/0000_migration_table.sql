create table migration (
  id      int auto_increment not null,
  current int not null,
  primary key (`id`)
);
insert into migration(current) values(0);