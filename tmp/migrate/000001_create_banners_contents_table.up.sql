create table banners_contents (
                                  banner_content_id serial primary key,
                                  content jsonb
);

create index banners_contents_banner_content_id on banners_contents(banner_content_id);
