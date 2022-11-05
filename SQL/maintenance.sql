create table ls_maint
(
    id         text primary key,
    region     smallint    not null,
    title      text        not null,
    uri        text        not null unique,
    maint_body    text        not null,
    square_edit timestamptz not null default now(),
    date_found timestamptz not null default now()
);
alter table ls_maint add column maint_body text not null;
create unique index idx_id_region_maint on ls_maint (id, region);