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

-- migrating data
-- split ip_peers to ip addresses
insert into ip_addresses (ip_address, service_name)
select distinct ipa, ip3.service_name from  ip_addresses ia right join (
    select distinct source_ip as ipa, service_name from ip_peers ip inner join
                                                        (select distinct dest_ip as ipa, service_name sn2 from ip_peers) ip2 on ip.source_ip=ip2.ipa
) AS ip3 on ip3.ipa=ip_address where ia.service_name is null;

-- create temporary table for service names
drop table if exists tmp_service_update;
create temporary table tmp_service_update as
    (select distinct ia.address_id, ipa, ip3.service_name from  ip_addresses ia join (
        select distinct source_ip as ipa, service_name from ip_peers ip inner join
                                                            (select distinct dest_ip as ipa, service_name sn2 from ip_peers) ip2 on ip.source_ip=ip2.ipa
    ) AS ip3 on ip3.ipa=ip_address where ia.service_name != ip3.service_name);

-- update service names
update ip_addresses set service_name=t.service_name
from (select address_id, service_name from tmp_service_update) as t
where ip_addresses.address_id=t.address_id;

drop table if exists tmp_service_update;

-- create temporary table to convert ip_peers into ip_peers2
drop table if exists tmp_address_update;
create temporary table tmp_address_update as
    (select
         ia.address_id as source_id,
         ip.source_port,
         ia2.address_id as dest_id,
         ip.dest_port,
         --case when coalesce(length(ia.service_name) ,0)=0 then ia2.service_name else ia.service_name end as service_name,
         ip.capture_id
     from
         ip_peers ip
             join ip_addresses ia on
             ia.ip_address = ip.source_ip
                 and ip.service_name = ia.service_name
             join ip_addresses ia2 on
             ia2.ip_address = ip.source_ip
                 and ip.service_name = ia2.service_name);

-- update existing records
update ip_peers2 set capture_id=t.capture_id
from (select p2.peer_id, tt.source_id, tt.source_port, tt.dest_id, tt.dest_port, tt.capture_id
      from tmp_address_update tt
               join ip_peers2 p2 on p2.source_id=tt.source_id and
                                    p2.source_port=tt.source_port and
                                    p2.dest_id=tt.dest_id and
                                    p2.dest_port=tt.dest_port) as t
where ip_peers2.peer_id=t.peer_id;

-- insert new records
insert into ip_peers2 (source_id, source_port, dest_id, dest_port, capture_id)
select tt.source_id, tt.source_port, tt.dest_id, tt.dest_port, tt.capture_id
from   tmp_address_update tt
           join ip_peers2 p2 on p2.source_id=tt.source_id and
                                p2.source_port=tt.source_port and
                                p2.dest_id=tt.dest_id and
                                p2.dest_port=tt.dest_port
where p2.peer_id is null;
-- drop table
drop table if exists tmp_address_update;

-- drop obsolete tables
--drop table if exists ip_peers;
-- complete migration
--alter table ip_peers2 rename to ip_peers;