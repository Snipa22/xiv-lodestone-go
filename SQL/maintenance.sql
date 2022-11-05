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
create unique index idx_id_region_maint on ls_maint (id, region);

create table ls_status
(
    id         text primary key,
    region     smallint    not null,
    title      text        not null,
    uri        text        not null unique,
    status_body    text        not null,
    square_edit timestamptz not null default now(),
    date_found timestamptz not null default now()
);
create unique index idx_id_region_status on ls_status (id, region);

create table ls_notices
(
    id         text primary key,
    region     smallint    not null,
    title      text        not null,
    uri        text        not null unique,
    notice_body    text        not null,
    square_edit timestamptz not null default now(),
    date_found timestamptz not null default now()
);
create unique index idx_id_region_notice on ls_notices (id, region);