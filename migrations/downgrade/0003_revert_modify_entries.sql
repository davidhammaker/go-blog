alter table entries modify ref text not null;
alter table entries drop column refHost;
alter table entries drop column refPath;
update migration set current = 2;