alter table entries modify ref text null;
alter table entries add refHost text null;
alter table entries add refPath text null;
update migration set current = 3;