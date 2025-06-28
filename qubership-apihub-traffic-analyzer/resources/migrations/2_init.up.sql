-- Copyright 2024-2025 NetCracker Technology Corporation
--
-- Licensed under the Apache License, Version 2.0 (the "License");
-- you may not use this file except in compliance with the License.
-- You may obtain a copy of the License at
--
--     http://www.apache.org/licenses/LICENSE-2.0
--
-- Unless required by applicable law or agreed to in writing, software
-- distributed under the License is distributed on an "AS IS" BASIS,
-- WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
-- See the License for the specific language governing permissions and
-- limitations under the License.

-- to store capture metadata
CREATE TABLE IF NOT EXISTS capture_metadata (
    capture_id varchar(36) NOT NULL,
    capture_metadata json NOT NULL,
    CONSTRAINT capture_metadata_pk PRIMARY KEY (capture_id)
);

-- sequence for ip_addresses table
CREATE SEQUENCE IF NOT EXISTS ip_address_seq
    INCREMENT BY 1
    MINVALUE 1
    START 1
    CACHE 1
    NO CYCLE;

-- table to store ip_address with service name
CREATE TABLE IF NOT EXISTS ip_addresses (
    address_id bigserial NOT NULL,
    ip_address varchar(16) NOT NULL,
    service_name varchar(256),
    CONSTRAINT ip_addresses_pk PRIMARY KEY (address_id));
-- indexes for ip_addresses
CREATE INDEX IF NOT EXISTS ip_address_ip_idx ON ip_addresses (ip_address);
CREATE INDEX IF NOT EXISTS ip_address_name_idx ON ip_addresses (service_name);

-- connection peers within capture
CREATE TABLE IF NOT EXISTS ip_peers2 (
    peer_id bigserial NOT NULL,
    source_id int8 NOT NULL,
    source_port INT NOT NULL,
    dest_id int8 NOT NULL,
    dest_port INT NOT NULL,
    capture_id varchar(36),
    CONSTRAINT ip_peers2_pk PRIMARY KEY (peer_id));
-- indexes for ip_peers2
CREATE INDEX IF NOT EXISTS ip_peers2_ip ON ip_peers2 (source_id, dest_id);