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
