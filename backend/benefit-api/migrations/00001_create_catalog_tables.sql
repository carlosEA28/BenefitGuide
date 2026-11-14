-- +goose Up
CREATE TABLE public_services (
    id UUID PRIMARY KEY,
    external_id VARCHAR(150),
    source_system VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    category VARCHAR(100) NOT NULL,
    agency VARCHAR(255) NOT NULL,
    status VARCHAR(30) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_synced_at TIMESTAMPTZ,
    CONSTRAINT public_services_source_system_check
        CHECK (source_system IN ('GOV_BR', 'MANUAL')),
    CONSTRAINT public_services_status_check
        CHECK (status IN ('ACTIVE', 'INACTIVE'))
);

CREATE UNIQUE INDEX public_services_source_external_id_key
    ON public_services (source_system, external_id)
    WHERE external_id IS NOT NULL;

CREATE TABLE service_requirements (
    id UUID PRIMARY KEY,
    public_service_id UUID NOT NULL,
    description TEXT NOT NULL,
    type VARCHAR(30) NOT NULL,
    condition_description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT service_requirements_public_service_id_fkey
        FOREIGN KEY (public_service_id) REFERENCES public_services (id) ON DELETE CASCADE,
    CONSTRAINT service_requirements_type_check
        CHECK (type IN ('MANDATORY', 'CONDITIONAL'))
);

CREATE INDEX service_requirements_public_service_id_idx
    ON service_requirements (public_service_id);

CREATE TABLE official_sources (
    id UUID PRIMARY KEY,
    public_service_id UUID NOT NULL,
    title VARCHAR(255) NOT NULL,
    url TEXT NOT NULL,
    agency VARCHAR(255) NOT NULL,
    last_verified_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT official_sources_public_service_id_fkey
        FOREIGN KEY (public_service_id) REFERENCES public_services (id) ON DELETE RESTRICT
);

CREATE INDEX official_sources_public_service_id_idx
    ON official_sources (public_service_id);

-- +goose Down
DROP TABLE official_sources;
DROP TABLE service_requirements;
DROP TABLE public_services;
