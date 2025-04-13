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

drop SEQUENCE if exists ip_peers_seq CASCADE;
drop SEQUENCE if exists ip_packets_seq CASCADE;
DROP TABLE IF EXISTS packets CASCADE;
DROP TABLE IF EXISTS packet_headers CASCADE;
DROP TABLE IF EXISTS http_headers CASCADE;
DROP TABLE IF EXISTS ip_peers CASCADE;
-- packets sequence
CREATE SEQUENCE IF NOT EXISTS ip_peers_seq
    INCREMENT BY 1
    MINVALUE 1
    START 1
    CACHE 1
    NO CYCLE;

-- ip_peers definition
CREATE TABLE IF NOT EXISTS ip_peers (
      peer_id int8 DEFAULT nextval('ip_peers_seq') NOT NULL,
      source_ip varchar(16) NOT NULL,
      source_port INT NOT NULL,
      dest_ip varchar(16) NOT NULL,
      dest_port INT,
      service_name varchar(256),
      capture_id varchar(36),
      CONSTRAINT ip_peers_pk PRIMARY KEY (peer_id));

CREATE INDEX IF NOT EXISTS ip_peers_capture_id_idx ON ip_peers (capture_id);
CREATE INDEX IF NOT EXISTS ip_peers_ServiceName_idx ON ip_peers (service_name);
CREATE INDEX IF NOT EXISTS ip_peers_source_ip_idx ON ip_peers (source_ip,dest_ip);
CREATE INDEX IF NOT EXISTS ip_peers_dest_ip_idx ON ip_peers (dest_ip,source_ip);

-- http_headers definition

CREATE TABLE IF NOT EXISTS http_headers (
      header_id varchar(36) NOT NULL,
      name varchar(256) NOT NULL,
      value TEXT,
      CONSTRAINT http_headers_pk PRIMARY KEY (header_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS http_headers_name_idx ON http_headers (Name,Value);

-- packets sequence
CREATE SEQUENCE IF NOT EXISTS ip_packets_seq
    INCREMENT BY 1
    MINVALUE 1
    START 1
    CACHE 1
    NO CYCLE;

-- packets definition
CREATE TABLE IF NOT EXISTS ip_packets (
       packet_id int8 DEFAULT nextval('ip_packets_seq') NOT NULL,
       peer_id int8 NOT NULL,
       seq_no int8 NOT NULL,
       ack_no int8 NOT NULL,
       time_stamp timestamp with time zone NOT NULL,
       body TEXT,
       capture_id varchar(36) NOT NULL,
       request_uri TEXT,
       request_method TEXT,
       CONSTRAINT ip_packets_pk PRIMARY KEY (packet_id),
       CONSTRAINT packets_ip_peers_fk FOREIGN KEY (peer_id) REFERENCES ip_peers(peer_id)
);

CREATE INDEX IF NOT EXISTS packets_peer_id1_idx ON ip_packets (peer_id,seq_no,ack_no);
CREATE INDEX IF NOT EXISTS packets_peer_id2_idx ON ip_packets (peer_id,ack_no,seq_no);
CREATE INDEX IF NOT EXISTS packets_capture_id_idx ON ip_packets (capture_id);

-- packet_headers definition

CREATE TABLE IF NOT EXISTS packet_headers (
    packet_id BIGINT NOT NULL,
    header_id TEXT NOT NULL,
    CONSTRAINT packet_headers_pk PRIMARY KEY (header_id,packet_id),
    CONSTRAINT packet_headers_http_headers_fk FOREIGN KEY (header_id) REFERENCES http_headers(header_id) ON DELETE CASCADE,
    CONSTRAINT packet_headers_packets_fk FOREIGN KEY (packet_id) REFERENCES ip_packets(packet_id) ON DELETE CASCADE
);

create table if not exists capture_metadata (
    capture_id varchar(36) NOT NULL,
    capture_metadata json not null,
    CONSTRAINT capture_metadata_pk PRIMARY KEY (capture_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS capture_metadata_uidx ON capture_metadata (capture_id);