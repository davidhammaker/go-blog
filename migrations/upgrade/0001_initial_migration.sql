create table entries (
  id      int auto_increment not null,
  title   varchar(128) not null,
  ref     text not null,
  created timestamp default current_timestamp,
  updated timestamp default current_timestamp on update current_timestamp,
  primary key (`id`)
);
update migration set current = 1;