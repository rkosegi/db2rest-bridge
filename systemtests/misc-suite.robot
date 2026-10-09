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
Resource                variables/common.robot
Resource                libraries/mysql.robot
Resource                libraries/db2rest.robot
Suite Setup             Setup
Suite Teardown          Teardown
Documentation           Miscellaneous test suite

*** Variables ***
${DDL_CREATE}           systemtests/data/create-geo.sql
${DDL_DROP}             systemtests/data/drop-geo.sql

*** Keywords ***
Setup
    [Documentation]             Setup this suite
    Builtin.Log                 Setting up suite
    Builtin.Log Variables
    OperatingSystem.Remove Files     .cache/systemtests/*.log   .cache/systemtests/*.err    ${APP_CONFIG}
    Builtin.Log                 Creating config
    ${dsn}                      Db2rest.Make DSN    ${MYSQL_APP_USER}    ${MYSQL_APP_PASS}     ${MYSQL_HOST}   ${MYSQL_DB}
    Db2rest.Write Config        ${APP_CONFIG}  demo  ${dsn}
    ${continentsConfig}         Builtin.Create Dictionary   page_size_limit=${5}
    ${countriesConfig}          Builtin.Create Dictionary   page_size_limit=${25}
    ${entitiesConfig}           Builtin.Create Dictionary   continent=${continentsConfig}   country=${countriesConfig}
    ${demoConfig}               Builtin.Create Dictionary   entities=${entitiesConfig}    page_size_limit=${10}
    ${backends}                 Builtin.Create Dictionary   demo=${demoConfig}
    ${extra}                    Builtin.Create Dictionary   backends=${backends}
    Db2rest.Merge To Config     ${APP_CONFIG}   ${extra}
    Db2rest.Init Session
    Builtin.Log                 Starting server
    Process.Start Process       go run pkg/cmd/main.go --config ${APP_CONFIG}
    ...                         cwd=.    alias=Server   shell=True
    ...                         stdout=.cache/systemtests/app.log  stderr=.cache/systemtests/app.err
    MySQL.Write Config          ${MYSQL_CNF_FILE}     ${MYSQL_HOST}   ${MYSQL_DDL_USER}     ${MYSQL_DDL_PASS}   ${MYSQL_DB}
    MySQL.Run client            ${MYSQL_CNF_FILE}     ${DDL_DROP}     drop
    MySQL.Run client            ${MYSQL_CNF_FILE}     ${DDL_CREATE}     create
    Builtin.Sleep               1s

Teardown
    [Documentation]             Tear down this suite
    MySQL.Run client            ${MYSQL_CNF_FILE}     ${DDL_DROP}   delete
    Process.Terminate All Processes    kill=True
    RequestsLibrary.Delete All Sessions
    Builtin.Log Variables


*** Test Cases ***
Check server version
    [Documentation]         Check server version
    Builtin.Log             Checking server version
    ${version}              Builtin.Wait Until Keyword Succeeds     10   1 sec   db2rest.Get Server Version
    Builtin.Log             ${version}


Verify Page Size Limit Apply For Entity Override
    ${items}                List Simple     demo    country
    ${size}                 Builtin.Get Length    ${items['data']}
    Builtin.Should Be Equal As Integers      ${size}   25

Verify Page Size Limit Apply For Backend Override
    [Documentation]
    ${items}                List Simple     demo    city
    ${size}                 Builtin.Get Length    ${items['data']}
    Builtin.Should Be Equal As Integers      ${size}   10
