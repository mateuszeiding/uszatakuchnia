-- +goose Up

DROP TABLE IF EXISTS aromas CASCADE;
DROP TABLE IF EXISTS ingredient_taste CASCADE;
DROP TABLE IF EXISTS taste CASCADE;
DROP TABLE IF EXISTS ingredients CASCADE;
DROP TABLE IF EXISTS aroma_types CASCADE;
DROP TABLE IF EXISTS ingredient_types CASCADE;
DROP TABLE IF EXISTS taste_types CASCADE;
DROP TABLE IF EXISTS images CASCADE;

-- =============================================
-- taste_types
CREATE TABLE taste_types (
    id   BIGSERIAL PRIMARY KEY,
    code VARCHAR(64) NOT NULL UNIQUE,
    name VARCHAR(128) NOT NULL
);

-- ingredient_types
CREATE TABLE ingredient_types (
    id   BIGSERIAL PRIMARY KEY,
    code VARCHAR(64) NOT NULL UNIQUE,
    name VARCHAR(128) NOT NULL
);

-- aroma_types
CREATE TABLE aroma_types (
    id   BIGSERIAL PRIMARY KEY,
    code VARCHAR(64) NOT NULL UNIQUE,
    name VARCHAR(128) NOT NULL
);

-- ingredients
CREATE TABLE ingredients (
    id          BIGSERIAL PRIMARY KEY,
    parent_id   BIGINT REFERENCES ingredients(id) ON DELETE SET NULL,
    name        VARCHAR(255) NOT NULL,
    type_id     BIGINT NOT NULL REFERENCES ingredient_types(id) ON DELETE RESTRICT,
    is_allergen BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE UNIQUE INDEX ux_ingredients_parent_name ON ingredients(parent_id, name);
CREATE INDEX ix_ingredients_parent ON ingredients(parent_id);
CREATE INDEX ix_ingredients_type ON ingredients(type_id);

-- taste
CREATE TABLE taste (
    id      BIGSERIAL PRIMARY KEY,
    type_id BIGINT NOT NULL UNIQUE REFERENCES taste_types(id) ON DELETE RESTRICT
);

-- ingredient_taste
CREATE TABLE ingredient_taste (
    ingredient_id BIGINT NOT NULL REFERENCES ingredients(id) ON DELETE RESTRICT,
    taste_id      BIGINT NOT NULL REFERENCES taste(id) ON DELETE RESTRICT,
    intensity     INT NOT NULL DEFAULT 50,
    PRIMARY KEY (ingredient_id, taste_id)
);

CREATE INDEX ix_ingtaste_taste ON ingredient_taste(taste_id);

-- aromas
CREATE TABLE aromas (
    id            BIGSERIAL PRIMARY KEY,
    ingredient_id BIGINT NOT NULL REFERENCES ingredients(id) ON DELETE RESTRICT,
    name          VARCHAR(255) NOT NULL,
    type_id       BIGINT NOT NULL REFERENCES aroma_types(id) ON DELETE RESTRICT,
    intensity     INT NOT NULL DEFAULT 50
);

CREATE UNIQUE INDEX ux_aromas_ing_name ON aromas(ingredient_id, name);
CREATE INDEX ix_aromas_ingredient ON aromas(ingredient_id);
CREATE INDEX ix_aromas_type ON aromas(type_id);

-- images
CREATE TABLE images (
    id                 BIGSERIAL PRIMARY KEY,
    ref_type           VARCHAR(32) NOT NULL,
    ref_id             BIGINT NOT NULL,
    provider           VARCHAR(32) NOT NULL,
    photo_id           VARCHAR(255) NOT NULL,
    image_url_small    TEXT NOT NULL,
    image_url_regular  TEXT NOT NULL,
    image_url_raw      TEXT NOT NULL,
    author_name        VARCHAR(255) NOT NULL,
    author_username    VARCHAR(255) NOT NULL,
    author_profile_url TEXT NOT NULL,
    photo_page_url     TEXT NOT NULL
);

CREATE UNIQUE INDEX ux_images_ref ON images(ref_type, ref_id);

-- +goose Down

DROP TABLE IF EXISTS images CASCADE;
DROP TABLE IF EXISTS aromas CASCADE;
DROP TABLE IF EXISTS ingredient_taste CASCADE;
DROP TABLE IF EXISTS taste CASCADE;
DROP TABLE IF EXISTS ingredients CASCADE;
DROP TABLE IF EXISTS aroma_types CASCADE;
DROP TABLE IF EXISTS ingredient_types CASCADE;
DROP TABLE IF EXISTS taste_types CASCADE;