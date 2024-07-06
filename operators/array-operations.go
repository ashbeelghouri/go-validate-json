package operators

func ArrayOfObjToObj(i interface{}, attr map[string]interface{}) *interface{} {
	arr := i.([]interface{})
	uniqueValueKey, ok := attr["unique_string_key"].(string)
	if !ok {
		return nil
	}
	newMap := make(map[string]interface{})
	for _, o := range arr {
		obj, ok := o.(map[string]interface{})
		if !ok {
			return nil
		}
		uk, exists := obj[uniqueValueKey].(string)
		if !exists {
			return nil
		}
		newMap[uk] = obj
	}
	var results interface{} = newMap
	return &results
}
