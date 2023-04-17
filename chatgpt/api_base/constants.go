package api_base

const (
    ErrorTypeInsufficientQuota  = "insufficient_quota"
    ErrorTypeInvalidRequest     = "invalid_request_error"
    ErrorTypeBillingNotActive   = "billing_not_active"
    ErrorTypeAccessTerminated   = "access_terminated"
    ErrorTypeAccountDeactivated = "account_deactivated"
    ErrorTypeInvalidApiKey      = "invalid_api_key"
)

/*

"error": {
"message": "You requested a model that is not compatible with this engine. Please contact us through our help center at help.openai.com for further questions.",
"type": "invalid_request_error",
"param": "model",
"code": null
}

"error": {
"message": "You exceeded your current quota, please check your plan and billing details.",
"type": "insufficient_quota",
"param": null,
"code": null
}

"error": {
"message": "Your account is not active, please check your billing details on our website.",
"type": "billing_not_active",
"param": null,
"code": null
}

"error": {
"message": "Your access was terminated due to violation of our policies, please check your email for more information. If you believe this is in error and would like to appeal, please contact support@openai.com.",
"type": "access_terminated",
"param": null,
"code": null
}

"error": {
"message": "Incorrect API key provided: sk-5cV63*************************************RkTE. You can find your API key at https://platform.openai.com/account/api-keys.",
"type": "invalid_request_error",
"param": null,
"code": "invalid_api_key"
}

"error": {
"message": "This key is associated with a deactivated account. If you feel this is an error, contact us through our help center at help.openai.com.",
"type": "invalid_request_error",
"param": null,
"code": "account_deactivated"
}

*/
