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
Documentation           CLI suite

*** Variables ***
${SUITE_ID}             cli
${DDL_CREATE}           systemtests/data/create-geo.sql
${DDL_DROP}             systemtests/data/drop-geo.sql

*** Keywords ***
Setup
    [Documentation]             Setup this suite
    Builtin.Log                 Setting up suite
    Builtin.Set Suite Variable  ${APP_CONFIG}   .cache/systemtests/config-${SUITE_ID}.yaml
    Builtin.Log Variables
    OperatingSystem.Remove Files    ${APP_CONFIG}
    Builtin.Log                 Creating config
    ${dsn}                      Db2rest.Make DSN    ${MYSQL_APP_USER}    ${MYSQL_APP_PASS}     ${MYSQL_HOST}   ${MYSQL_DB}
    ${defaultClientProfile}     Builtin.Create Dictionary   url=http://127.0.0.1:22001/api/v1
    ${clientProfiles}           Builtin.Create Dictionary   default=${defaultClientProfile}
    ${extra}                    Builtin.Create Dictionary   active_client_profile=default    client_profiles=${clientProfiles}
    ${baseConfig}               Db2rest.Create Config   demo  ${dsn}
    ${config}                   Db2rest.Merge Recursive     ${baseConfig}   ${extra}
    Db2rest.Write To Yaml       ${APP_CONFIG}   ${config}
    Db2rest.Init Session
    Builtin.Log                 Starting server ${SUITE_ID}
    Process.Start Process       go run pkg/cli/main.go server --log-level debug --config ${APP_CONFIG}
    ...                         cwd=.    alias=Server   shell=True
    ...                         stdout=.cache/systemtests/app-${SUITE_ID}.log  stderr=.cache/systemtests/app-${SUITE_ID}.err
    MySQL.Write Config          ${MYSQL_CNF_FILE}     ${MYSQL_HOST}   ${MYSQL_DDL_USER}     ${MYSQL_DDL_PASS}   ${MYSQL_DB}
    MySQL.Run client            ${MYSQL_CNF_FILE}     ${DDL_DROP}     drop
    MySQL.Run client            ${MYSQL_CNF_FILE}     ${DDL_CREATE}     create
    Builtin.Sleep               1s
    Db2rest.Wait For Server

Teardown
    [Documentation]             Tear down this suite
    MySQL.Run client            ${MYSQL_CNF_FILE}     ${DDL_DROP}   delete
    Process.Terminate All Processes    kill=True
    RequestsLibrary.Delete All Sessions
    Builtin.Log Variables


*** Test Cases ***
Show version
    Db2rest.Run CLI And Expect Success     show-version    version

List Backends
    Db2rest.Run CLI And Expect Success     list-backends    client backends
