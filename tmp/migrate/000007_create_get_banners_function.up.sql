create type get_banners_returned_type as (
    banner_id int,
    feature_id int,
    tag_ids int[],
    content jsonb,
    is_active bool,
    created_at time,
    updated_at time
);

create function get_banners(l int, o int, t_id int, f_id int) returns setof get_banners_returned_type
as $$
begin
return Query (SELECT
            b.banner_id,
            b.feature_id,
            array_agg(bht.tag_id) as tag_ids,
            bc.content,
            b.is_active,
            b.created_at,
            b.updated_at
        FROM
            banners b
                JOIN
            banners_contents bc ON b.content_id = bc.banner_content_id
                LEFT JOIN
            banners_has_tags bht ON b.banner_id = bht.banner_id
      where (t_id is null or bht.tag_id = t_id) and (f_id is null or b.feature_id = f_id)
      GROUP BY
            b.banner_id, bc.content limit l offset o rows);
end;
$$ language plpgsql;