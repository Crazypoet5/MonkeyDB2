create table bat1(
    bat_id int,
    bat_name string
);
insert into bat1 values(1,'bat1'),(2,'bat2'),(3,'bat3');
insert into bat1 values(4,'bat1'),(5,'bat2'),(6,'bat3');
update bat1 set bat_name = 'bat33' where bat_name='bat3';
select * from bat1;
quit;
