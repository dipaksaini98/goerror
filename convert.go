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
		var traceObjects []map[string]interface{}
		for _, traceErr := range trace {
			traceObj := map[string]interface{}{
				"title":   GetTitle(traceErr),
				"message": traceErr.Error(),
				"type":    GetType(traceErr),
				"display": GetDisplay(traceErr),
			}
			// if original error exists, add it to the trace object
			if originalErr := GetOriginalError(traceErr); originalErr != nil && originalErr != traceErr {
				traceObj["original_error"] = originalErr.Error()
			}
			traceObjects = append(traceObjects, traceObj)
		}
		errMap["trace"] = traceObjects
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
