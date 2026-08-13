alter table button
    add column sort_order int default 0;

update button
set sort_order = 0
where data in ('ru', 'yes');
update button
set sort_order = 1
where data in ('en', 'no');