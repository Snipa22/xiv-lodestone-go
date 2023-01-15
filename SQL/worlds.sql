create table sq_physical_datacenters
(
    id            int primary key,
    internal_name text NOT NULL,
    display_name  text not null
);
create index idx_spd_dn_lower on sq_physical_datacenters (lower(display_name));
create index idx_spd_in_lower on sq_physical_datacenters (lower(internal_name));
create table sq_logical_datacenters
(
    id             SERIAL primary key,
    public_id      int  not null,
    internal_name  text NOT NULL,
    physical_dc_id int  not null,
    display_name   text not null,
    CONSTRAINT fk_physical_datacenter FOREIGN KEY (physical_dc_id) REFERENCES sq_physical_datacenters (id)
);
create unique index idx_logical_dc_public_phys on sq_logical_datacenters (public_id, physical_dc_id);
create index idx_logical_dc_phys_id on sq_logical_datacenters (physical_dc_id);
create index idx_sld_dn_lower on sq_logical_datacenters (lower(display_name));
create index idx_sld_in_lower on sq_logical_datacenters (lower(internal_name));
create table sq_worlds
(
    id                       SERIAL primary key,
    internal_name            text NOT NULL,
    internal_id              int  not null,
    sq_logical_datacenter_id int  not null,
    display_name             text not null,
    CONSTRAINT fk_sq_logical_datacenter FOREIGN KEY (sq_logical_datacenter_id) REFERENCES sq_logical_datacenters (id)
);
create index idx_worlds_logical_datacenter_id on sq_worlds (sq_logical_datacenter_id);
create unique index idx_worlds_internal_id on sq_worlds (internal_id);
create index idx_sw_dn_lower on sq_worlds (lower(display_name));
create index idx_sw_in_lower on sq_worlds (lower(internal_name));

insert into sq_physical_datacenters (id, internal_name, display_name)
values (1, 'North American Data Center', 'NA'),
       (2, 'European Data Center', 'EU'),
       (3, 'Oceanian Data Center', 'OCE'),
       (4, 'Japanese Data Center', 'JP'),
       (5, 'Development', 'DEV'),
       (6, 'China Data Center', 'CN'),
       (7, 'Korea Data Center', 'KR');

-- For this, we're only importing the NA, OCE, EU, JP, aka "public" on https://raw.githubusercontent.com/xivapi/ffxiv-datamining/master/csv/World.csv
insert into sq_logical_datacenters (public_id, internal_name, physical_dc_id, display_name)
values (1, 'Elemental', 4, 'Elemental'),
       (2, 'Gaia', 4, 'Gaia'),
       (3, 'Mana', 4, 'Mana'),
       (4, 'Aether', 1, 'Aether'),
       (5, 'Primal', 1, 'Primal'),
       (6, 'Chaos', 2, 'Chaos'),
       (7, 'Light', 2, 'Light'),
       (8, 'Crystal', 1, 'Crystal'),
       (9, 'Materia', 3, 'Materia'),
       (10, 'Meteor', 4, 'Meteor'),
       (11, 'Dynamis', 1, 'Dynamis');

-- Perform final mapping
insert into sq_worlds (internal_id, internal_name, display_name, sq_logical_datacenter_id)
values (45, 'Carbuncle', 'Carbuncle', (select id from sq_logical_datacenters where internal_name = 'Elemental')),
       (49, 'Kujata', 'Kujata', (select id from sq_logical_datacenters where internal_name = 'Elemental')),
       (50, 'Typhon', 'Typhon', (select id from sq_logical_datacenters where internal_name = 'Elemental')),
       (58, 'Garuda', 'Garuda', (select id from sq_logical_datacenters where internal_name = 'Elemental')),
       (68, 'Atomos', 'Atomos', (select id from sq_logical_datacenters where internal_name = 'Elemental')),
       (72, 'Tonberry', 'Tonberry', (select id from sq_logical_datacenters where internal_name = 'Elemental')),
       (90, 'Aegis', 'Aegis', (select id from sq_logical_datacenters where internal_name = 'Elemental')),
       (94, 'Gungnir', 'Gungnir', (select id from sq_logical_datacenters where internal_name = 'Elemental')),
       (43, 'Alexander', 'Alexander', (select id from sq_logical_datacenters where internal_name = 'Gaia')),
       (46, 'Fenrir', 'Fenrir', (select id from sq_logical_datacenters where internal_name = 'Gaia')),
       (51, 'Ultima', 'Ultima', (select id from sq_logical_datacenters where internal_name = 'Gaia')),
       (59, 'Ifrit', 'Ifrit', (select id from sq_logical_datacenters where internal_name = 'Gaia')),
       (69, 'Bahamut', 'Bahamut', (select id from sq_logical_datacenters where internal_name = 'Gaia')),
       (76, 'Tiamat', 'Tiamat', (select id from sq_logical_datacenters where internal_name = 'Gaia')),
       (92, 'Durandal', 'Durandal', (select id from sq_logical_datacenters where internal_name = 'Gaia')),
       (98, 'Ridill', 'Ridill', (select id from sq_logical_datacenters where internal_name = 'Gaia')),
       (23, 'Asura', 'Asura', (select id from sq_logical_datacenters where internal_name = 'Mana')),
       (28, 'Pandaemonium', 'Pandaemonium', (select id from sq_logical_datacenters where internal_name = 'Mana')),
       (44, 'Anima', 'Anima', (select id from sq_logical_datacenters where internal_name = 'Mana')),
       (47, 'Hades', 'Hades', (select id from sq_logical_datacenters where internal_name = 'Mana')),
       (48, 'Ixion', 'Ixion', (select id from sq_logical_datacenters where internal_name = 'Mana')),
       (61, 'Titan', 'Titan', (select id from sq_logical_datacenters where internal_name = 'Mana')),
       (70, 'Chocobo', 'Chocobo', (select id from sq_logical_datacenters where internal_name = 'Mana')),
       (96, 'Masamune', 'Masamune', (select id from sq_logical_datacenters where internal_name = 'Mana')),
       (40, 'Jenova', 'Jenova', (select id from sq_logical_datacenters where internal_name = 'Aether')),
       (54, 'Faerie', 'Faerie', (select id from sq_logical_datacenters where internal_name = 'Aether')),
       (57, 'Siren', 'Siren', (select id from sq_logical_datacenters where internal_name = 'Aether')),
       (63, 'Gilgamesh', 'Gilgamesh', (select id from sq_logical_datacenters where internal_name = 'Aether')),
       (65, 'Midgardsormr', 'Midgardsormr', (select id from sq_logical_datacenters where internal_name = 'Aether')),
       (73, 'Adamantoise', 'Adamantoise', (select id from sq_logical_datacenters where internal_name = 'Aether')),
       (79, 'Cactuar', 'Cactuar', (select id from sq_logical_datacenters where internal_name = 'Aether')),
       (99, 'Sargatanas', 'Sargatanas', (select id from sq_logical_datacenters where internal_name = 'Aether')),
       (35, 'Famfrit', 'Famfrit', (select id from sq_logical_datacenters where internal_name = 'Primal')),
       (53, 'Exodus', 'Exodus', (select id from sq_logical_datacenters where internal_name = 'Primal')),
       (55, 'Lamia', 'Lamia', (select id from sq_logical_datacenters where internal_name = 'Primal')),
       (64, 'Leviathan', 'Leviathan', (select id from sq_logical_datacenters where internal_name = 'Primal')),
       (77, 'Ultros', 'Ultros', (select id from sq_logical_datacenters where internal_name = 'Primal')),
       (78, 'Behemoth', 'Behemoth', (select id from sq_logical_datacenters where internal_name = 'Primal')),
       (93, 'Excalibur', 'Excalibur', (select id from sq_logical_datacenters where internal_name = 'Primal')),
       (95, 'Hyperion', 'Hyperion', (select id from sq_logical_datacenters where internal_name = 'Primal')),
       (39, 'Omega', 'Omega', (select id from sq_logical_datacenters where internal_name = 'Chaos')),
       (71, 'Moogle', 'Moogle', (select id from sq_logical_datacenters where internal_name = 'Chaos')),
       (80, 'Cerberus', 'Cerberus', (select id from sq_logical_datacenters where internal_name = 'Chaos')),
       (83, 'Louisoix', 'Louisoix', (select id from sq_logical_datacenters where internal_name = 'Chaos')),
       (85, 'Spriggan', 'Spriggan', (select id from sq_logical_datacenters where internal_name = 'Chaos')),
       (97, 'Ragnarok', 'Ragnarok', (select id from sq_logical_datacenters where internal_name = 'Chaos')),
       (400, 'Sagittarius', 'Sagittarius', (select id from sq_logical_datacenters where internal_name = 'Chaos')),
       (401, 'Phantom', 'Phantom', (select id from sq_logical_datacenters where internal_name = 'Chaos')),
       (33, 'Twintania', 'Twintania', (select id from sq_logical_datacenters where internal_name = 'Light')),
       (36, 'Lich', 'Lich', (select id from sq_logical_datacenters where internal_name = 'Light')),
       (42, 'Zodiark', 'Zodiark', (select id from sq_logical_datacenters where internal_name = 'Light')),
       (56, 'Phoenix', 'Phoenix', (select id from sq_logical_datacenters where internal_name = 'Light')),
       (66, 'Odin', 'Odin', (select id from sq_logical_datacenters where internal_name = 'Light')),
       (67, 'Shiva', 'Shiva', (select id from sq_logical_datacenters where internal_name = 'Light')),
       (402, 'Alpha', 'Alpha', (select id from sq_logical_datacenters where internal_name = 'Light')),
       (403, 'Raiden', 'Raiden', (select id from sq_logical_datacenters where internal_name = 'Light')),
       (34, 'Brynhildr', 'Brynhildr', (select id from sq_logical_datacenters where internal_name = 'Crystal')),
       (37, 'Mateus', 'Mateus', (select id from sq_logical_datacenters where internal_name = 'Crystal')),
       (41, 'Zalera', 'Zalera', (select id from sq_logical_datacenters where internal_name = 'Crystal')),
       (62, 'Diabolos', 'Diabolos', (select id from sq_logical_datacenters where internal_name = 'Crystal')),
       (74, 'Coeurl', 'Coeurl', (select id from sq_logical_datacenters where internal_name = 'Crystal')),
       (75, 'Malboro', 'Malboro', (select id from sq_logical_datacenters where internal_name = 'Crystal')),
       (81, 'Goblin', 'Goblin', (select id from sq_logical_datacenters where internal_name = 'Crystal')),
       (91, 'Balmung', 'Balmung', (select id from sq_logical_datacenters where internal_name = 'Crystal')),
       (21, 'Ravana', 'Ravana', (select id from sq_logical_datacenters where internal_name = 'Materia')),
       (22, 'Bismarck', 'Bismarck', (select id from sq_logical_datacenters where internal_name = 'Materia')),
       (86, 'Sephirot', 'Sephirot', (select id from sq_logical_datacenters where internal_name = 'Materia')),
       (87, 'Sophia', 'Sophia', (select id from sq_logical_datacenters where internal_name = 'Materia')),
       (88, 'Zurvan', 'Zurvan', (select id from sq_logical_datacenters where internal_name = 'Materia')),
       (24, 'Belias', 'Belias', (select id from sq_logical_datacenters where internal_name = 'Meteor')),
       (29, 'Shinryu', 'Shinryu', (select id from sq_logical_datacenters where internal_name = 'Meteor')),
       (30, 'Unicorn', 'Unicorn', (select id from sq_logical_datacenters where internal_name = 'Meteor')),
       (31, 'Yojimbo', 'Yojimbo', (select id from sq_logical_datacenters where internal_name = 'Meteor')),
       (32, 'Zeromus', 'Zeromus', (select id from sq_logical_datacenters where internal_name = 'Meteor')),
       (52, 'Valefor', 'Valefor', (select id from sq_logical_datacenters where internal_name = 'Meteor')),
       (60, 'Ramuh', 'Ramuh', (select id from sq_logical_datacenters where internal_name = 'Meteor')),
       (82, 'Mandragora', 'Mandragora', (select id from sq_logical_datacenters where internal_name = 'Meteor')),
       (404, 'Marilith', 'Marilith', (select id from sq_logical_datacenters where internal_name = 'Dynamis')),
       (405, 'Seraph', 'Seraph', (select id from sq_logical_datacenters where internal_name = 'Dynamis')),
       (406, 'Halicarnassus', 'Halicarnassus', (select id from sq_logical_datacenters where internal_name = 'Dynamis')),
       (407, 'Maduin', 'Maduin', (select id from sq_logical_datacenters where internal_name = 'Dynamis'));