alter table entries drop column description;
update migration set current = 1;