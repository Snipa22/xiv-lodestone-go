create table items
(
    id           text primary key,
    world_id     smallint    not null,
    item_id      smallint    not null,
    price        int         not null,
    total        int         not null,
    hq           bool        not null,
    quantity     int         not null,
    date_updated timestamptz not null default now()
);
create index idx_items_id_world on items (id, world_id);
create index idx_items_item_id_world on items (item_id, world_id);