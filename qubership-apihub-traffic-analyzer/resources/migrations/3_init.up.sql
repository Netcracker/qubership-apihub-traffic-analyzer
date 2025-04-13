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

-- just to resolve some constraints
CREATE TABLE if not exists http_headers (
    header_id varchar(36) NOT NULL,
    "name" varchar(256) NOT NULL,
    value text NULL,
    CONSTRAINT http_headers_pk PRIMARY KEY (header_id)
);
CREATE UNIQUE INDEX if not exists http_headers_name_idx ON http_headers USING btree (name, value);

-- report statuses
CREATE TABLE if not exists report_status (
    report_status_id serial4 NOT NULL,
    report_status varchar NOT NULL,
    created_at timestamp DEFAULT now() NOT NULL,
    retired_at timestamp NULL,
    CONSTRAINT report_status_pk PRIMARY KEY (report_status_id)
);

-- report types
CREATE TABLE if not exists report_types (
    report_type_id serial4 NOT NULL,
    report_type varchar NOT NULL,
    created_at timestamp DEFAULT now() NOT NULL,
    retired_at timestamp NULL,
    CONSTRAINT report_types_pk PRIMARY KEY (report_type_id)
);

-- stored_reports reports stored in DB
CREATE TABLE if not exists stored_reports (
    report_id bigserial NOT NULL,
    created_at timestamp DEFAULT now() NOT NULL,
    report_parameters json NOT NULL,
    report_type_id int4 NOT NULL,
    report_status_id int4 NOT NULL,
    completed_at timestamp NULL,
    report_uuid varchar NOT NULL,
    CONSTRAINT stored_reports_pk PRIMARY KEY (report_id),
    CONSTRAINT stored_reports_unique UNIQUE (report_uuid),
    CONSTRAINT stored_reports_report_status_fk FOREIGN KEY (report_status_id) REFERENCES report_status(report_status_id) ON DELETE CASCADE,
    CONSTRAINT stored_reports_report_types_fk FOREIGN KEY (report_type_id) REFERENCES report_types(report_type_id) ON DELETE CASCADE
);
-- stored_reports indexes
CREATE INDEX if not exists stored_reports_report_status_id_idx ON stored_reports USING btree (report_status_id, completed_at);
CREATE INDEX if not exists stored_reports_report_type_id_idx ON stored_reports USING btree (report_type_id);
-- stored_reports column comments
COMMENT ON COLUMN stored_reports.report_id IS 'primary key';
COMMENT ON COLUMN stored_reports.created_at IS 'creation timestamp';
COMMENT ON COLUMN stored_reports.report_parameters IS 'parameters from report creation request';
COMMENT ON COLUMN stored_reports.report_type_id IS 'report type to distinct the report contents';
COMMENT ON COLUMN stored_reports.report_status_id IS 'report status  (created|completed|failed|...)';
COMMENT ON COLUMN stored_reports.completed_at IS 'completion timestamp';

-- report_data report data "rows"
CREATE TABLE if not exists report_data (
    report_row_id bigserial NOT NULL,
    report_id bigint NOT NULL,
    report_row json NULL,
    CONSTRAINT report_data_pk PRIMARY KEY (report_row_id),
    CONSTRAINT report_data_report_fk FOREIGN KEY (report_id) REFERENCES stored_reports(report_id) ON DELETE CASCADE
);
-- report_data indexes
CREATE INDEX if not exists report_data_report_id_idx ON report_data (report_id);
-- report_data column comments
COMMENT ON COLUMN report_data.report_row_id IS 'primary key';
COMMENT ON COLUMN report_data.report_id IS 'reference to report';
COMMENT ON COLUMN report_data.report_row IS 'row data (JSON)';

-- report_service_operations service operations from backend to make report
CREATE TABLE if not exists report_service_operations (
    report_operation_id bigserial NOT NULL,
    report_id int8 NOT NULL,
    operation_title varchar NULL,
    operation_path varchar NOT NULL,
    operation_path_re varchar NULL,
    operation_method varchar NOT NULL,
    operation_status varchar NULL,
    operation_hit_count int DEFAULT 0 NOT NULL,
    CONSTRAINT report_service_operations_pk PRIMARY KEY (report_operation_id),
    CONSTRAINT report_service_operations_stored_reports_fk FOREIGN KEY (report_id) REFERENCES stored_reports(report_id)
);
-- report_service_operations indexes
CREATE INDEX if not exists report_service_operations_operation_path_idx ON report_service_operations USING btree (operation_path);
CREATE INDEX if not exists report_service_operations_report_id_idx ON report_service_operations USING btree (report_id);
-- report_service_operations column comments
COMMENT ON COLUMN report_service_operations.report_operation_id IS 'primary key';
COMMENT ON COLUMN report_service_operations.report_id IS 'reference to report';
COMMENT ON COLUMN report_service_operations.operation_title IS 'operation''s title';
COMMENT ON COLUMN report_service_operations.operation_path IS 'operation''s request path';
COMMENT ON COLUMN report_service_operations.operation_path_re IS 'operation''s request path regular expression to deal with substitutions';
COMMENT ON COLUMN report_service_operations.operation_method IS 'operation''s method (GET/POST/...)';
COMMENT ON COLUMN report_service_operations.operation_status IS 'operation''s status against the capture data (Captured/Not captured/...)';
COMMENT ON COLUMN report_service_operations.operation_hit_count IS 'how often operation occurred in the dump';

-- report_affected_rows references to affected rows
CREATE TABLE if not exists report_affected_rows (
    report_id bigint NOT NULL,
    reference_id bigint NOT NULL,
    reference_type int NOT NULL,
    hit_count int NOT NULL,
    CONSTRAINT report_affected_rows_pk PRIMARY KEY (report_id,reference_id,reference_type),
    CONSTRAINT report_affected_rows_stored_reports_fk FOREIGN KEY (report_id) REFERENCES stored_reports(report_id)
);
-- report_affected_rows column comments
COMMENT ON COLUMN report_affected_rows.report_id IS 'reference to report';
COMMENT ON COLUMN report_affected_rows.reference_id IS 'reference to ip_packets of report_service_operations';
COMMENT ON COLUMN report_affected_rows.reference_type IS 'when 1 then ip_packets, when 2 then report_service_operations';

-- service_addresses
CREATE TABLE if not exists service_addresses (
    address_id bigserial NOT NULL,
    ip_address varchar NOT NULL,
    service_name varchar NULL,
    service_version varchar NULL,
    capture_id varchar NOT NULL,
    CONSTRAINT service_address_pk PRIMARY KEY (address_id)
);
-- service_addresses indexes
CREATE INDEX if not exists service_address_ip_idx ON service_addresses USING btree (ip_address);
CREATE INDEX if not exists service_address_name_idx ON service_addresses USING btree (service_name, service_version);
CREATE INDEX if not exists service_addresses_capture_id_idx ON service_addresses USING btree (capture_id);
-- service_addresses column comments
COMMENT ON COLUMN service_addresses.address_id IS 'primary key';
COMMENT ON COLUMN service_addresses.ip_address IS 'a captured IP address';
COMMENT ON COLUMN service_addresses.service_name IS 'a service name for this address or empty when no service bound on';
COMMENT ON COLUMN service_addresses.service_version IS 'a service version caught from the cloud';
COMMENT ON COLUMN service_addresses.capture_id IS 'a capture id where the row is valid';

-- service_packets captured packets
CREATE TABLE if not exists service_packets (
    packet_id bigserial NOT NULL,
    source_id int8 NOT NULL,
    source_port int4 NOT NULL,
    dest_id int8 NOT NULL,
    dest_port int4 NOT NULL,
    seq_no int8 NOT NULL,
    ack_no int8 NOT NULL,
    "time_stamp" timestamptz NOT NULL,
    body text NULL,
    capture_id varchar NOT NULL,
    request_path text NULL,
    request_method text NULL,
    CONSTRAINT service_packets_pk PRIMARY KEY (packet_id),
    CONSTRAINT packets_source_svc_fk FOREIGN KEY (source_id) REFERENCES service_addresses(address_id),
    CONSTRAINT packets_dest_svc_fk FOREIGN KEY (dest_id) REFERENCES service_addresses(address_id)
);
-- service_packets index
CREATE INDEX if not exists service_packets_address_idx ON service_packets (source_id, dest_id);
-- service_packets column comments
COMMENT ON COLUMN service_packets.packet_id IS 'primary key';
COMMENT ON COLUMN service_packets.source_id IS 'source IP address';
COMMENT ON COLUMN service_packets.source_port IS 'source TCP port';
COMMENT ON COLUMN service_packets.dest_id IS 'dest IP address';
COMMENT ON COLUMN service_packets.dest_port IS 'dest TCP port';
COMMENT ON COLUMN service_packets.seq_no IS 'TCP sequence number';
COMMENT ON COLUMN service_packets.ack_no IS 'TCP acknowlege number';
COMMENT ON COLUMN service_packets."time_stamp" IS 'packet time stamp';
COMMENT ON COLUMN service_packets.body IS 'TCP packet payload';
COMMENT ON COLUMN service_packets.capture_id IS 'capture identifier';
COMMENT ON COLUMN service_packets.request_path IS 'HTTP request path';
COMMENT ON COLUMN service_packets.request_method IS 'HTTP request method';

-- service_packet_headers headers for packet
CREATE TABLE if not exists service_packet_headers (
    packet_id int8 NOT NULL,
    header_id varchar(36) NOT NULL,
    CONSTRAINT service_packet_headers_pk PRIMARY KEY (header_id, packet_id),
    CONSTRAINT service_packet_headers_http_headers_fk FOREIGN KEY (header_id) REFERENCES http_headers(header_id) ON DELETE CASCADE,
    CONSTRAINT service_packet_headers_ip_packets_fk FOREIGN KEY (packet_id) REFERENCES service_packets(packet_id) ON DELETE CASCADE
);
-- service_packet_headers column comments
COMMENT ON COLUMN service_packet_headers.packet_id IS 'reference to service packets';
COMMENT ON COLUMN service_packet_headers.header_id IS 'reference to service packet headers';

CREATE TABLE if not exists report_service_operations2 (
    report_id int8 NULL,
    src_peer varchar NULL,
    dst_peer varchar NULL,
    operation_title varchar NULL,
    operation_path varchar NULL,
    operation_method varchar NULL,
    operation_status varchar NULL,
    hit_count int8 NULL,
    CONSTRAINT report_service_operations2_stored_reports_fk FOREIGN KEY (report_id) REFERENCES stored_reports(report_id) ON DELETE CASCADE
);
CREATE INDEX if not exists report_service_operations2_report_id_idx ON report_service_operations2 USING btree (report_id);
