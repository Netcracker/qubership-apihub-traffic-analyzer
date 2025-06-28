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

insert into report_types (report_type_id, report_type) values (1, 'service operations') on conflict do nothing;
insert into report_status (report_status_id, report_status) values(1, 'created') on conflict do nothing;
insert into report_status (report_status_id, report_status) values(2, 'ready') on conflict do nothing;
insert into report_status (report_status_id, report_status) values(3, 'in progress') on conflict do nothing;
insert into report_status (report_status_id, report_status) values(4, 'failed') on conflict do nothing;
