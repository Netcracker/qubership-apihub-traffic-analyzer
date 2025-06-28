// Copyright 2024-2025 NetCracker Technology Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package exception

const EmptyParameter = "8"
const EmptyParameterMsg = "Parameter $param should not be empty"

const InvalidParameterValue = "9"
const InvalidParameterValueMsg = "Value '$value' is not allowed for parameter $param"

const BadRequestBody = "10"
const BadRequestBodyMsg = "Failed to decode body"

const ContentIdNotFound = "40"
const ContentIdNotFoundMsg = "Content with id $contentId not found in branch $branch for project $projectId"

const ApiKeyNotFound = "83"
const ApiKeyNotFoundMsg = "Api key for user $user and integration $integration not found"

const NoApihubAccess = "200"
const NoApihubAccessMsg = "No access to Apihub with code: $code. Probably incorrect configuration: api key."

const ReportGenerationTimeOut = "20000"
const ReportGenerationTimeOutMsg = "Report generation takes more time than expected"
