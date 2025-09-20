package tablewriter

// AttributeGetter is a function type that retrieves a field value from an instance
type AttributeGetter func(fieldName string, instance any) (string, error)
