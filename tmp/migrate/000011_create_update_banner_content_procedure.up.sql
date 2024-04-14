create procedure update_banner_content(b_id int, c jsonb)
    as $$
declare c_id int;
begin
        select content_id into c_id from banners where banner_id = b_id;

update banners_contents set content = c where banner_content_id = c_id;
end;
$$ language plpgsql;