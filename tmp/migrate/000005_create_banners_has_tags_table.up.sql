create table banners_has_tags (
                                  banner_id int not null,
                                  tag_id int not null,
                                  foreign key (banner_id) references banners(banner_id)
                                      on delete cascade
                                      on update cascade
);

create index banners_has_tags_banner_id on banners_has_tags(banner_id);
