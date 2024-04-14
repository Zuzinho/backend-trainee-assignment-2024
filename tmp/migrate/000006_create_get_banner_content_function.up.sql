create type user_role_enum as enum('User', 'Admin');

create function get_banner_content(t_id int, f_id int, role user_role_enum) returns jsonb
as $$
DECLARE b_id int; c_id int; is_act boolean;
begin
    select banner_id into b_id from banners join banners_has_tags bht on banners.banner_id = bht.banner_id where feature_id = f_id and bht.tag_id = t_id;
    select content_id into c_id from banners where banner_id = b_id;

    select is_active into is_act from banners where banner_id = b_id;
    if not is_act and role like 'User' then
        return null;
    end if;

    return (select content from banners_contents where banner_content_id = c_id);
end;
$$ language plpgsql;