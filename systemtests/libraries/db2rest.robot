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
Library    Collections
Library    JSONLibrary
Library    RequestsLibrary

*** Keywords ***
Init Session
    [Documentation]     Create requests session
    ${headers}          Create Dictionary   Accept=application/json     Content-Type=application/json
    ${session}          RequestsLibrary.Create Session   client     http://127.0.0.1:22001/api/v1   headers=${headers}
    RETURN              ${session}

Make DSN
    [Documentation]     Builds a DSN from args
    [Arguments]         ${user}    ${pass}     ${host}   ${db}
    ${result}           String.Format String    {}:{}@tcp({})/{}?parseTime=true     ${user}     ${pass}     ${host}     ${db}
    RETURN              ${result}

Write Config
    [Documentation]     Writes configuration to specified YAML file
    [Arguments]         ${path}  ${backend}  ${dsn}
    ${logging}          Builtin.Create Dictionary   level=debug
    ${backendObj}       Builtin.Create Dictionary   create=${True}  read=${True}  update=${True}  delete=${True}  dsn=${dsn}
    ${backends}         Builtin.Create Dictionary   ${backend}=${backendObj}
    ${body}             Builtin.Create Dictionary   logging=${logging}      backends=${backends}
    JSONLibrary.Dump Json To File   ${path}     ${body}

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
    ${uri}              String.Format String    /{}/{}  ${backend}  ${entity}
    ${response}         REST Post       ${uri}  ${body}
    RequestsLibrary.Status Should Be    201    ${response}
    RETURN              ${response.json()}

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
    ${uri}              String.Format String    /{}/{}/{}    ${backend}  ${entity}    ${id}
    ${response}         REST Delete    ${uri}
    RequestsLibrary.Status Should Be    204    ${response}
    RETURN              ${response}

REST Get
    [Documentation]     Make a GET call and return response as a dictionary
    [Arguments]         ${uri}
    ${result}           RequestsLibrary.Get On Session    alias=client    url=${uri}
    RETURN              ${result}

REST Post
    [Documentation]     Make a POST call and return response as a dictionary
    [Arguments]         ${uri}  ${body}
    ${result}           RequestsLibrary.Post On Session    alias=client    url=${uri}  json=${body}
    RETURN              ${result}

REST Put
    [Documentation]     Make a PUT call and return response as a dictionary
    [Arguments]         ${uri}  ${body}
    ${result}           RequestsLibrary.Put On Session    alias=client    url=${uri}  json=${body}
    RETURN              ${result}

REST Delete
    [Documentation]     Make a DELETE call and return response as a dictionary
    [Arguments]         ${uri}
    ${result}           RequestsLibrary.Delete On Session    alias=client    url=${uri}
    RETURN              ${result}
