create procedure insert_banner_tags(b_id int, t_ids int[])
    as $$
declare t_id int;
begin
        foreach t_id in array t_ids
            loop
                insert into banners_has_tags (banner_id, tag_id) values (b_id, t_id);
end loop;
end;
$$ language plpgsql;
