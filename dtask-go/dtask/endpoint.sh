#!/bin/bash

curl -X POST http://localhost:8082/open/api/dtask/task/list -d '{ "pagingVo" : { "limit" : 10, "page": 1}}' -H 'content-type:application/json' -H 'id:1' -H 'username:zhuangyongj' -H 'userno:A123123' -H 'role:admin'

curl -X POST http://localhost:8082/open/api/dtask/task/list -d '{ "jobName" : "Delete", "pagingVo" : { "limit" : 10, "page": 1}}' -H 'content-type:application/json' -H 'id:1' -H 'username:zhuangyongj' -H 'userno:A123123' -H 'role:admin'

curl -X POST http://localhost:8082/open/api/dtask/task/history -d '{ "jobName" : "", "pagingVo" : { "limit" : 10, "page": 1}}' -H 'content-type:application/json' -H 'id:1' -H 'username:zhuangyongj' -H 'userno:A123123' -H 'role:admin'

curl -X POST http://localhost:8082/open/api/dtask/task/trigger -d '{ "id" : 1 }' -H 'content-type:application/json' -H 'id:1' -H 'username:zhuangyongj' -H 'userno:A123123' -H 'role:admin'