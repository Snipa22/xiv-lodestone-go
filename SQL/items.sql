create table items
(
    id           text primary key,
    world_id     smallint    not null,
    item_id      integer    not null,
    price        int         not null,
    total        int         not null,
    hq           bool        not null,
    quantity     int         not null,
    date_updated timestamptz not null default now()
);
create index idx_items_id_world on items (id, world_id);
create index idx_items_item_id_world on items (item_id, world_id);

create table sales
(
    id          serial primary key,
    hq          bool        not null,
    date_loaded timestamptz not null default now(),
    price       int         not null,
    quantity    int         not null,
    total       int         not null,
    world_id    smallint    not null,
    item_id     integer    not null
);
create index idx_sales_item_id_world_time on sales (item_id, world_id, date_loaded);