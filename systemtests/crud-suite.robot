# Copyright 2026 Richard Kosegi
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

*** Settings ***
Resource                libraries/mysql.robot
Resource                libraries/db2rest.robot
Suite Setup             Setup
Suite Teardown          Teardown
Documentation           CRUD test suite 1

*** Variables ***
${MYSQL_HOST}           %{MYSQL_HOST=localhost}
${MYSQL_APP_USER}       %{MYSQL_APP_USER=demo}
${MYSQL_APP_PASS}       %{MYSQL_APP_PASS=123456}
${MYSQL_DDL_USER}       %{MYSQL_DDL_USER=root}
${MYSQL_DDL_PASS}       %{MYSQL_DDL_PASS=123456}
${MYSQL_DB}             %{MYSQL_DB=demo}
${MYSQL_PORT}           %{MYSQL_PORT=3306}
${MYSQL_CNF_FILE}       .cache/systemtests/.my.cnf
${APP_CONFIG}           .cache/systemtests/config.yaml
${DDL_CREATE}           systemtests/data/create.sql
${DDL_DROP}             systemtests/data/drop.sql

*** Keywords ***
Setup
    [Documentation]         Setup this suite
    Builtin.Log             Setting up suite
    Builtin.Log Variables
    OperatingSystem.Remove Files     .cache/systemtests/*.log   .cache/systemtests/*.err    ${APP_CONFIG}
    Builtin.Log             Creating config
    ${dsn}                  Db2rest.Make DSN    ${MYSQL_APP_USER}    ${MYSQL_APP_PASS}     ${MYSQL_HOST}   ${MYSQL_DB}
    Db2rest.Write Config    ${APP_CONFIG}  demo  ${dsn}
    Db2rest.Init Session
    Builtin.Log             Starting server
    Process.Start Process   go run pkg/cmd/main.go --config ${APP_CONFIG}
    ...                     cwd=.    alias=Server   shell=True
    ...                     stdout=.cache/systemtests/app.log  stderr=.cache/systemtests/app.err
    MySQL.Write Config      ${MYSQL_CNF_FILE}     ${MYSQL_HOST}   ${MYSQL_DDL_USER}     ${MYSQL_DDL_PASS}   ${MYSQL_DB}
    MySQL.Run client        ${MYSQL_CNF_FILE}     ${DDL_DROP}     drop
    MySQL.Run client        ${MYSQL_CNF_FILE}     ${DDL_CREATE}     create
    Builtin.Sleep           1s

Teardown
    [Documentation]         Tear down this suite
    MySQL.Run client        ${MYSQL_CNF_FILE}     ${DDL_DROP}   delete
    Process.Terminate All Processes    kill=True
    RequestsLibrary.Delete All Sessions
    Builtin.Log Variables

*** Test Cases ***
Check server version
    [Documentation]         Check server version
    Builtin.Log             Checking server version
    ${version}              Builtin.Wait Until Keyword Succeeds     10   1 sec   db2rest.Get Server Version
    Builtin.Log             ${version}

Verify backend exists
    [Documentation]         Check if list of configured backends reported by REST API contains one we need
    ${backends}             db2rest.List Configured Backends
    Builtin.Should Contain  ${backends}     ${MYSQL_DB}

Create/Update/Delete single item
    [Documentation]         Create single item and verify it exists. Then perform update/delete.
    &{emp}                  Builtin.Create Dictionary   name=Alice  salary=500  department=HR
    ${item}                 Db2rest.Create Item     ${MYSQL_DB}     employee    ${emp}
    ${itemId}               Collections.Get From Dictionary    ${item}    id
    ${patch}                Builtin.Create Dictionary   name=Bob  salary=200  department=Management
    Db2rest.Get Item By Id  ${MYSQL_DB}     employee    ${itemId}
    Db2rest.Update Item By Id  ${MYSQL_DB}     employee    ${itemId}    ${patch}
    ${item}                 Db2rest.Get Item By Id  ${MYSQL_DB}     employee    ${itemId}
    Builtin.Should be equal      ${item['name']}   Bob
    Builtin.Should be equal      ${item['department']}   Management
    Db2rest.Delete Item By Id  ${MYSQL_DB}     employee    ${itemId}

Create same item twice should fail
    [Documentation]         Attempt to create item with same primary key should fail for second call
    &{emp}                  Builtin.Create Dictionary   name=Charlie  id=100      salary=500  department=HR
    ${item}                 Db2rest.Create Item And Expect Status     ${MYSQL_DB}     employee    ${emp}    201
    ${item}                 Db2rest.Create Item And Expect Status     ${MYSQL_DB}     employee    ${emp}    409

Create item with missing properties should fail
    [Documentation]         Attempt to create item mandatory (non-null) property should fail
    &{emp}                  Builtin.Create Dictionary   name=Hugo
    ${item}                 Db2rest.Create Item And Expect Status     ${MYSQL_DB}     employee    ${emp}    400

Create item with invalid foreign key ref should fail
    [Documentation]         Attempt to create item with invalid reference to other table should fail
    &{emp}                  Builtin.Create Dictionary   employee_id=999    name=prop1  val=XYZ
    ${item}                 Db2rest.Create Item And Expect Status     ${MYSQL_DB}     employee_property    ${emp}    409


Delete item with existing foreign key reference should fail
    [Documentation]         Attempt to delete item with incoming reference from other table should fail
    &{emp}                  Builtin.Create Dictionary   name=Charlie  id=200      salary=500  department=HR
    ${item}                 Db2rest.Create Item And Expect Status     ${MYSQL_DB}     employee    ${emp}    201
    &{prop}                 Builtin.Create Dictionary   id=500    employee_id=200    name=prop1  val=XYZ
    ${item}                 Db2rest.Create Item And Expect Status     ${MYSQL_DB}     employee_property    ${prop}    201
    ${result}               Db2rest.Delete Item And Expect Status  ${MYSQL_DB}     employee    200    409
    ${result}               Db2rest.Delete Item By Id  ${MYSQL_DB}     employee_property    500
    ${result}               Db2rest.Delete Item By Id  ${MYSQL_DB}     employee    200

Bulk Create/Update/Delete
    [Documentation]         Mutate multiple items at once using bulk operation API
    @{items}                Builtin.Create List
    FOR    ${index}         IN RANGE    0   20
        &{item}             Builtin.Create Dictionary   name=Employee-${index}  salary=100  department=Eng
        Collections.Append To List   ${items}    ${item}
    END
    Db2rest.Bulk Operation  ${MYSQL_DB}     employee    INSERT     ${items}
    ${items}                Db2rest.List Items With Predicate    ${MYSQL_DB}     employee    salary    100
    ${size}                 Builtin.Get Length    ${items['data']}
    Builtin.Should Be Equal As Integers      ${size}   20
    Db2rest.Bulk Operation  ${MYSQL_DB}     employee    DELETE     ${items['data']}
    ${items}                Db2rest.List Items With Predicate    ${MYSQL_DB}     employee    salary    100
    ${size}                 Builtin.Get Length    ${items['data']}
    Builtin.Should Be Equal As Integers      ${size}   0
