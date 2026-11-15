-- +goose Up
ALTER TABLE public_services
    ALTER COLUMN id SET DEFAULT gen_random_uuid();

ALTER TABLE service_requirements
    ALTER COLUMN id SET DEFAULT gen_random_uuid();

ALTER TABLE official_sources
    ALTER COLUMN id SET DEFAULT gen_random_uuid();

ALTER TABLE service_requirements
    ADD COLUMN display_order INTEGER NOT NULL DEFAULT 0;

ALTER TABLE service_requirements
    ADD CONSTRAINT service_requirements_conditional_description_check
    CHECK (
        type <> 'CONDITIONAL'
        OR (condition_description IS NOT NULL AND btrim(condition_description) <> '')
    );

-- +goose Down
ALTER TABLE service_requirements
    DROP CONSTRAINT service_requirements_conditional_description_check;

ALTER TABLE service_requirements
    DROP COLUMN display_order;

ALTER TABLE public_services
    ALTER COLUMN id DROP DEFAULT;

ALTER TABLE service_requirements
    ALTER COLUMN id DROP DEFAULT;

ALTER TABLE official_sources
    ALTER COLUMN id DROP DEFAULT;
