package util

import (
	"fmt"
	"reflect"

	"github.com/nam2184/mymy/util"
)

type SortOrder string

const (
    ASC  SortOrder = "ASC"
    DESC SortOrder = "DESC"
)

func GetSkipLimit(query_params map[string]interface{}) (int, int, error) {
	defaultSkip := 0
	defaultLimit := 10

	// Get "skip" value from query_params
	skip, ok := util.GetIntFromInterface(query_params["skip"])
	if !ok {
		skip = defaultSkip
	}

	// Get "limit" value from query_params
	limit, ok := util.GetIntFromInterface(query_params["limit"])
	if !ok {
		limit = defaultLimit
	}

	// Optional: Validate skip and limit if needed
	if skip < 0 || limit <= 0 {
		return 0, 0, fmt.Errorf("invalid values: 'skip' must be >= 0 and 'limit' must be > 0")
	}
	return skip, limit, nil
}

func ValidateQueries[T any](query_params map[string]interface{}) map[string]interface{} {
    // Get the type of the generic struct
    if len(query_params) == 0 {
      return query_params
    }
    t := reflect.TypeOf((*T)(nil)).Elem()
    // Collect allowed keys based on struct's JSON tags
    allowedKeys := make(map[string]struct{})
    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        
        // Get the JSON tag, if available
        jsonTag := field.Tag.Get("json")
        if jsonTag != "" && jsonTag != "-" {
            allowedKeys[jsonTag] = struct{}{}
        }
    }
    
    // Filter query_params, keeping only keys that exist in allowedKeys
    filteredParams := make(map[string]interface{})
    for key, value := range query_params {
        if _, exists := allowedKeys[key]; exists {
            filteredParams[key] = value
        }
    }
    
    return filteredParams
}

func GetSortBy(query_params map[string]interface{}) (string, string) {

    sort_by, ok := query_params["sort_by"].(string)
    if !ok {
        sort_by = "id"
    }
    
    order, ok := query_params["order"].(string)
    if !ok {
        order = string(DESC) 
    }
    
    return sort_by, string(order)
}
