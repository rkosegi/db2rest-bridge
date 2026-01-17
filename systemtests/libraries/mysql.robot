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
Library    Process
Library    String
Library    OperatingSystem

*** Keywords ***
Write Config
    [Documentation]         Creates .my.cnf config file for client
    [Arguments]             ${path}  ${host}    ${user}   ${password}   ${db}     ${port}=3306
    ${content}              Builtin.Catenate    SEPARATOR=\n
    ...    [client]
    ...    host=${host}
    ...    database=${db}
    ...    port=3306
    ...    user=${user}
    ...    protocol=TCP
    ...    password=${password}

    OperatingSystem.Create File     ${path}     ${content}
    Process.Run Process     chmod 0600 ${path}
    ...                     shell=True

Run client
    [Documentation]         Run a DDL file against database
    [Arguments]             ${config_file}      ${ddl_file}     ${alias}
    ${cmd} =                Builtin.Catenate    mysql   --defaults-file=${config_file}  --verbose      --execute="SOURCE ${ddl_file};"
    Builtin.Log             About to invoke ${cmd}
    Process.Run Process     ${cmd}      shell=True    cwd=.    alias=mysql
    ...                     stdout=.cache/systemtests/mysql-${alias}.log
    ...                     stderr=.cache/systemtests/mysql-${alias}.err
    ${result}               Process.Get Process Result  mysql
    Builtin.Should Be Equal	${result.rc}    ${0}
