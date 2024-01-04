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
	errMap["message"] = err.Error()
	errMap["type"] = GetType(err)
	return errMap
}

// JSON convert error object into bytes json
func JSON(err error) []byte {
	jsonStr, _ := json.Marshal(Map(err))
	return jsonStr
}
