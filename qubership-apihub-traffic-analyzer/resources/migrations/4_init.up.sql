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

create table if not exists LOAD_PACKET_JOBS
(
    JOB_ID    integer not null,
    CREATED   TIMESTAMP,
    STARTED   timestamp,
    FINISHED  timestamp,
    JOB_TTL   timestamp,
    INSTANCE_ID varchar,
    CAPTURE_ID varchar,
    CONSTRAINT LOAD_PACKET_JOBS_PK PRIMARY KEY (JOB_ID)
);
CREATE INDEX if not exists LOAD_PACKET_JOBS_IDX ON LOAD_PACKET_JOBS USING btree (STARTED, FINISHED);
comment on table LOAD_PACKET_JOBS is 'Packet data loading jobs';
COMMENT ON COLUMN LOAD_PACKET_JOBS.JOB_ID IS 'primary key';
COMMENT ON COLUMN LOAD_PACKET_JOBS.CREATED IS 'when the record created';
COMMENT ON COLUMN LOAD_PACKET_JOBS.STARTED IS 'when the job taken';
COMMENT ON COLUMN LOAD_PACKET_JOBS.FINISHED IS 'when the job finished';
COMMENT ON COLUMN LOAD_PACKET_JOBS.JOB_TTL IS 'last TTL for the job';
COMMENT ON COLUMN LOAD_PACKET_JOBS.CAPTURE_ID IS 'capture id file name part';
COMMENT ON COLUMN LOAD_PACKET_JOBS.INSTANCE_ID IS 'instance id file name part';
insert into report_types (report_type_id, report_type) values (1, 'service operations') on conflict do nothing;
insert into report_status (report_status_id, report_status) values(1, 'created') on conflict do nothing;
insert into report_status (report_status_id, report_status) values(2, 'ready') on conflict do nothing;
insert into report_status (report_status_id, report_status) values(3, 'in progress') on conflict do nothing;
insert into report_status (report_status_id, report_status) values(4, 'failed') on conflict do nothing;
