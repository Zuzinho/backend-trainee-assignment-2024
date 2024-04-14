create function insert_banner(f_id int, t_ids int[], c jsonb, i_a bool) returns int8
as $$
DECLARE c_id int; t_id int; b_id int;
begin
insert into banners_contents (content) values (c) returning banner_content_id into c_id;
insert into banners (feature_id, content_id, is_active)
values (f_id, c_id, i_a) returning banner_id into b_id;
call insert_banner_tags(b_id, t_ids);

return b_id;
end
$$ language plpgsql;