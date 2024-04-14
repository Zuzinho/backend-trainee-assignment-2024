create function set_updated_at()
    returns trigger as $$
begin
    if new.created_at is not null
        or new.updated_at is not null then
        raise exception 'You can not assign value to created_at or updated_at';
end if;

    new.updated_at := now();
return new;
end
$$ language plpgsql;

create trigger banners_update
    before update
    on banners
    for each row
    execute function set_updated_at();