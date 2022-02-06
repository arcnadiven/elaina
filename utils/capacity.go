package utils

import "k8s.io/apimachinery/pkg/api/resource"

func BytesToQuantity(bytes int64) string {
	return resource.NewQuantity(bytes, resource.BinarySI).String()
}

func QuantityToBytes(source string) int64 {
	quantity, err := resource.ParseQuantity(source)
	if err != nil {
		return 0
	}
	return quantity.Value()
}
