package payment

import (
	"fmt"

	lru "github.com/hashicorp/golang-lru/v2"
)

func InitializeCache() (*lru.Cache[string, string], error) {
	cache, err := lru.New[string, string](128)
	if err != nil {
		fmt.Println("Error initializing cache:", err)
		return nil, err
	}

	// cache.Add(1, "Alice")
	// cache.Add(2, "Bob")

	// if val, ok := cache.Get(1); ok {
	// 	fmt.Println("Found:", val) // Output: Found: Alice
	// }

	return cache, nil
}
