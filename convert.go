package goerror

import "encoding/json"

// Map object of the error
func Map(err error) map[string]interface{} {
	emptyCtx := Context{}
	var errMap = make(map[string]interface{})
	if customErr, ok := err.(*GoError); ok && customErr.Context != emptyCtx {
		errMap["field"] = customErr.Context.Key
		errMap["message"] = customErr.Context.Value
		errMap["type"] = GetType(customErr)
		return errMap
	}
	errMap["title"] = GetTitle(err)
	errMap["message"] = err.Error()
	errMap["type"] = GetType(err)
	errMap["display"] = GetDisplay(err)

	// check if trace exists in the error
	if trace := GetTrace(err); trace != nil {
		errMap["trace"] = GetTrace(err)
	}

	// if original error exists, add it to the map
	if originalErr := GetOriginalError(err); originalErr != nil && originalErr != err {
		errMap["original_error"] = originalErr.Error()
	}

	return errMap
}

// JSON convert error object into bytes json
func JSON(err error) []byte {
	jsonStr, _ := json.Marshal(Map(err))
	return jsonStr
}
