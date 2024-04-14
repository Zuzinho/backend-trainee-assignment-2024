create function delete_banner()
    returns trigger as $$
begin
delete from banners_contents where banner_content_id = old.content_id;

return old;
end
$$ language plpgsql;

create trigger banners_delete
    before update
    on banners
    for each row
    execute function delete_banner();