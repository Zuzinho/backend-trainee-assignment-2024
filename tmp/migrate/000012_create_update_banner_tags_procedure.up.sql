create procedure update_banner_tags(b_id int, t_ids int[])
    as $$
begin
delete from banners_has_tags where banner_id = b_id;
call insert_banner_tags(b_id, t_ids);
end;
$$ language plpgsql;