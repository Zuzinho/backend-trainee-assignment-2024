create table users (
                       login text primary key,
                       password text not null,
                       role user_role_enum not null default 'User'
);

create index users_login on users(login);