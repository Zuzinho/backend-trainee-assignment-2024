create table banners (
                         banner_id serial primary key,
                         feature_id int not null,
                         content_id int not null,
                         is_active boolean not null default false,
                         created_at time not null,
                         updated_at time not null,
                         foreign key (content_id) references banners_contents(banner_content_id)
                             on delete cascade
                             on update cascade
);

create index banners_feature_id on banners(feature_id);