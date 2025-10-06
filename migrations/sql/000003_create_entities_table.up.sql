SET search_path TO scarlet;
create table entities (
    entity_id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    entity_type VARCHAR(100) NOT NULL,
    address TEXT,
    phone BIGINT,
    email VARCHAR(255),
    image_url TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    active BOOLEAN DEFAULT TRUE
);

create index idx_entities_name ON entities (name);
create index idx_entities_entity_type ON entities (entity_type);