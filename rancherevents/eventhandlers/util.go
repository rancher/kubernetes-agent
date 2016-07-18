package eventhandlers

func GetStringMap(data map[string]interface{}, keys ...string) map[string]string {
	for _, key := range keys {
		val, ok := data[key]
		if !ok {
			return nil
		}
		mapVal, ok := val.(map[string]interface{})
		if !ok {
			return nil
		}
		data = mapVal
	}

	result := map[string]string{}
	for k, v := range data {
		if s, ok := v.(string); ok {
			result[k] = s
		}
	}

	return result
}
