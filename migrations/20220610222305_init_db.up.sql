create table if not exists records
(
    id         varchar(36) default (lower(hex(randomblob(4)) || '-' || hex(randomblob(2))
        || '-' || '4' || substr(hex(randomblob(2)), 2) || '-'
        || substr('AB89', 1 + (abs(random()) % 4), 1) ||
                                          substr(hex(randomblob(2)), 2) || '-' || hex(randomblob(6)))) primary key,
    domain     text,
    ipv4       varchar(15),
    ipv6       varchar(39),
    created_at datetime    default CURRENT_TIMESTAMP,
    ttl        integer     default 604800,
    type       varchar(10)
);
