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
Library    OperatingSystem
Library    String
Library    yaml
Library    Collections
Library    JSONLibrary
Library    RequestsLibrary

*** Keywords ***
Init Session
    [Documentation]     Create requests session
    ${headers}          Create Dictionary   Accept=application/json     Content-Type=application/json
    ${session}          RequestsLibrary.Create Session   client     http://127.0.0.1:22001/api/v1   headers=${headers}
    RETURN              ${session}

Wait For Server
    [Documentation]     Check server version. This waits until server is available for 10 seconds at most
    Builtin.Log         Checking server version
    ${version}          Builtin.Wait Until Keyword Succeeds     10   1 sec   db2rest.Get Server Version
    Builtin.Log         ${version}

Make DSN
    [Documentation]     Builds a DSN from args
    [Arguments]         ${user}    ${pass}     ${host}   ${db}
    ${result}           String.Format String    {}:{}@tcp({})/{}?parseTime=true     ${user}     ${pass}     ${host}     ${db}
    RETURN              ${result}

Create Config
    [Documentation]     Creates base configuration object for backend
    [Arguments]         ${backend}  ${dsn}
    ${logging}          Builtin.Create Dictionary   level=debug
    ${backendObj}       Builtin.Create Dictionary   create=${True}  read=${True}  update=${True}  delete=${True}  dsn=${dsn}
    ${backends}         Builtin.Create Dictionary   ${backend}=${backendObj}
    ${result}           Builtin.Create Dictionary   logging=${logging}      backends=${backends}
    RETURN              ${result}

Write To Yaml
    [Documentation]     Writes object to the YAML file
    [Arguments]         ${path}     ${dict}
    ${str}              Evaluate    yaml.safe_dump(${dict}, default_flow_style=False)    modules=yaml
    OperatingSystem.Create File     ${path}    ${str}

Write Config
    [Documentation]     Writes configuration to specified file
    [Arguments]         ${path}  ${backend}  ${dsn}
    ${body}             Create Config   ${backend}  ${dsn}
    JSONLibrary.Dump Json To File   ${path}     ${body}

Merge To Config
    [Documentation]     Merges given dictionary to existing configuration file
    [Arguments]         ${path}     ${dict}
    ${json}             JSONLibrary.Load Json From File   ${path}
    ${json}             Merge Recursive    ${json}    ${dict}
    JSONLibrary.Dump Json To File   ${path}     ${json}

Merge Recursive
    [Documentation]     Merges 2 dictionaries recursively
    [Arguments]         ${d1}    ${d2}
    &{result}           Collections.Copy Dictionary    ${d1}

    FOR    ${key}    IN    @{d2.keys()}
        ${key_exists}=    Run Keyword And Return Status    Dictionary Should Contain Key    ${result}    ${key}
        IF    ${key_exists}
            ${is_dict1}=    Builtin.Evaluate    isinstance($result.get($key), dict)
            ${is_dict2}=    Builtin.Evaluate    isinstance($d2.get($key), dict)
            IF    ${is_dict1} and ${is_dict2}
                ${nested_merged}=    Merge Recursive    ${result['${key}']}    ${d2['${key}']}
                Collections.Set To Dictionary    ${result}    ${key}=${nested_merged}
            ELSE
                Collections.Set To Dictionary    ${result}    ${key}=${d2['${key}']}
            END
        ELSE
            Collections.Set To Dictionary    ${result}    ${key}=${d2['${key}']}
        END
    END
    RETURN    ${result}

List Configured Backends
    ${response}         REST Get        /backends
    RequestsLibrary.Status Should Be    200    ${response}
    RETURN              ${response.json()}

Get Server Version
    [Documentation]     Retrieve server version
    ${response}         REST Get        /version
    RequestsLibrary.Status Should Be    200    ${response}
    RETURN              ${response.json()}

Create Item
    [Documentation]     Create single item in backend
    [Arguments]         ${backend}     ${entity}   ${body}
    ${response}         Create Item And Expect Status    ${backend}    ${entity}   ${body}    201
    RETURN              ${response.json()}

Create Item And Expect Status
    [Documentation]     Create single item in backend
    [Arguments]         ${backend}     ${entity}   ${body}    ${status}
    ${uri}              String.Format String    /{}/{}  ${backend}  ${entity}
    ${response}         REST Post       ${uri}  ${body}
    RequestsLibrary.Status Should Be    ${status}    ${response}
    RETURN              ${response}

Bulk Operation
    [Documentation]     Create/Update/Delete multiple items at once
    [Arguments]         ${backend}     ${entity}   ${mode}  ${items}
    ${uri}              String.Format String    /{}/{}/bulk  ${backend}  ${entity}
    &{obj}              Builtin.Create Dictionary   mode=${mode}      objects=${items}
    ${response}         REST Post       ${uri}  ${obj}
    RequestsLibrary.Status Should Be    200    ${response}
    RETURN              ${response}

List Items With Predicate
    [Documentation]     List all items that matches simple predicate (key OP value)
    [Arguments]         ${backend}      ${entity}   ${key}  ${value}    ${op}==    ${offset}=0    ${size}=20
    &{filter}           Builtin.Create Dictionary   name=${key}     val=${value}    op=${op}
    &{filter}           Builtin.Create Dictionary   simple=${filter}
    ${filterStr}        Convert Json To String    ${filter}
    ${uri}              String.Format String    /{}/{}?page-offset={}&page-size={}&filter={}    ${backend}  ${entity}
    ...    ${offset}    ${size}    ${filterStr}
    ${response}         REST Get    ${uri}
    RequestsLibrary.Status Should Be    200    ${response}
    RETURN              ${response.json()}

List Simple
    [Documentation]     List items without further restrictions
    [Arguments]         ${backend}      ${entity}
    ${uri}              String.Format String    /{}/{}    ${backend}  ${entity}
    ${response}         REST Get    ${uri}
    RequestsLibrary.Status Should Be    200    ${response}
    RETURN              ${response.json()}

Get Item By Id
    [Documentation]     Retrieves item by ID
    [Arguments]         ${backend}    ${entity}    ${id}
    ${uri}              String.Format String    /{}/{}/{}    ${backend}  ${entity}    ${id}
    ${response}         REST Get    ${uri}
    RequestsLibrary.Status Should Be    200    ${response}
    RETURN              ${response.json()}

Update Item By Id
    [Documentation]     Update item by ID
    [Arguments]         ${backend}    ${entity}    ${id}    ${body}
    ${uri}              String.Format String    /{}/{}/{}    ${backend}  ${entity}    ${id}
    ${response}         REST Put    ${uri}    ${body}
    RequestsLibrary.Status Should Be    202    ${response}
    RETURN              ${response.json()}

Delete Item By Id
    [Documentation]     Delete item by ID
    [Arguments]         ${backend}    ${entity}    ${id}
    ${response}         Delete Item And Expect Status    ${backend}  ${entity}    ${id}    204
    RETURN              ${response}

Delete Item And Expect Status
    [Documentation]     Delete item by ID and expect status
    [Arguments]         ${backend}    ${entity}    ${id}    ${status}
    ${uri}              String.Format String    /{}/{}/{}    ${backend}  ${entity}    ${id}
    ${response}         REST Delete    ${uri}
    RequestsLibrary.Status Should Be    ${status}    ${response}
    RETURN              ${response}

REST Get
    [Documentation]     Make a GET call and return response as a dictionary
    [Arguments]         ${uri}
    ${result}           RequestsLibrary.Get On Session    alias=client    url=${uri}
    RETURN              ${result}

REST Post
    [Documentation]     Make a POST call and return response as a dictionary
    [Arguments]         ${uri}  ${body}
    ${result}           RequestsLibrary.Post On Session    alias=client    url=${uri}  json=${body}    expected_status=anything
    RETURN              ${result}

REST Put
    [Documentation]     Make a PUT call and return response as a dictionary
    [Arguments]         ${uri}  ${body}
    ${result}           RequestsLibrary.Put On Session    alias=client    url=${uri}  json=${body}
    RETURN              ${result}

REST Delete
    [Documentation]     Make a DELETE call and return response as a dictionary
    [Arguments]         ${uri}
    ${result}           RequestsLibrary.Delete On Session    alias=client    url=${uri}    expected_status=anything
    RETURN              ${result}
