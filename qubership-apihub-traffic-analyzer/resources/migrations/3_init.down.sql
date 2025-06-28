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

do
$$
    declare
        rec_cnt integer;
        rec_st varchar;
        statuses varchar[] := array ['created', 'ready', 'in progress', 'failed'];
    begin
        rec_st := 'service operations';
        select count(report_type_id) into rec_cnt from report_types where report_type=rec_st;
        if rec_cnt<1 then
            insert into report_types (report_type) values (rec_st);
        else
            raise info 'type record exists for %', rec_st;
        end if;
        foreach  rec_st in array statuses
            loop
                select count(report_status_id) into rec_cnt from report_status where report_status=rec_st;
                if rec_cnt<1 then
                    insert into report_status (report_status) values(rec_st);
                else
                    raise info 'status record exists for %', rec_st;
                end if;
            end loop;
    end;
$$