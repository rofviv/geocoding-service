package reponses

const (
	OK             = "OK"
	MISSING_PARAMS = "MISSING_PARAMS"
	FAILED         = "FAILED"
	DENIED         = "DENIED"
	UNKNOWN        = "UNKNOWN"
	INVALID_DATA   = "IVALID_DATA"
	ZERO_RESULTS   = "ZERO_RESULTS"

	OK_MESSAGE             = "request sent successfully"
	MISSING_PARAMS_MESSAGE = "you must submit all required fields"
	EMPTY_FIELD_MESSAGE    = "you cannot send empty values"
	FAILED_MESSAGE         = "invalid format json"
	DENIED_MESSAGE         = "access denied"
	UNKNOWN_MESSAGE        = "uknowk error server"
	INVALID_DATA_MESSAGE   = "The key value is invalid"
)
