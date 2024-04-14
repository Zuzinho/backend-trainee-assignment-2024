create function set_created_updated_at()
    returns trigger as $$
begin
    if new.created_at is not null
           or new.updated_at is not null then
        raise exception 'You can not assign value to created_at or updated_at';
end if;

    new.created_at := now();
    new.updated_at = new.created_at;
return new;
end
$$ language plpgsql;

create trigger banners_insert
    before insert
    on banners
    for each row
    execute function set_created_updated_at();