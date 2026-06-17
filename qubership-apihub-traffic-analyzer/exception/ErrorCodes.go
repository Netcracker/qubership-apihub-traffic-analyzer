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
