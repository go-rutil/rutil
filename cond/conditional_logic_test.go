package cond

import (
	"testing"
)

func TestIf(t *testing.T) {
	// String example
	t.Run("String test", func(t *testing.T) {
		name := "Alice"
		greeting := If(name == "Alice", "Hello Alice!", "Hello stranger!")
		expected := "Hello Alice!"
		if greeting != expected {
			t.Errorf("Expected %s, got %s", expected, greeting)
		}
	})

	// Integer example
	t.Run("Integer test", func(t *testing.T) {
		x := 10
		result := If(x > 5, x*2, x/2)
		expected := 20
		if result != expected {
			t.Errorf("Expected %d, got %d", expected, result)
		}
	})

	// Boolean example
	t.Run("Boolean test", func(t *testing.T) {
		isActive := true
		status := If(isActive, "Online", "Offline")
		expected := "Online"
		if status != expected {
			t.Errorf("Expected %s, got %s", expected, status)
		}
	})

	// Float example
	t.Run("Float test", func(t *testing.T) {
		temperature := 25.5
		category := If(temperature > 30.0, "Hot", "Mild")
		expected := "Mild"
		if category != expected {
			t.Errorf("Expected %s, got %s", expected, category)
		}
	})

	// Struct example
	t.Run("Struct test", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}

		age := 17
		person := If(age >= 18,
			Person{Name: "Adult", Age: age},
			Person{Name: "Minor", Age: age})
		expected := Person{Name: "Minor", Age: 17}
		if person != expected {
			t.Errorf("Expected %+v, got %+v", expected, person)
		}
	})

	// Pointer example - useful for optional values
	t.Run("Pointer test", func(t *testing.T) {
		var ptr *int
		defaultValue := 42
		value := If(ptr != nil, ptr, &defaultValue)
		expected := 42
		if *value != expected {
			t.Errorf("Expected %d, got %d", expected, *value)
		}
	})

	// Additional test cases for edge conditions
	t.Run("False condition string test", func(t *testing.T) {
		name := "Bob"
		greeting := If(name == "Alice", "Hello Alice!", "Hello stranger!")
		expected := "Hello stranger!"
		if greeting != expected {
			t.Errorf("Expected %s, got %s", expected, greeting)
		}
	})

	t.Run("False condition integer test", func(t *testing.T) {
		x := 3
		result := If(x > 5, x*2, x/2)
		expected := 1 // integer division
		if result != expected {
			t.Errorf("Expected %d, got %d", expected, result)
		}
	})

	t.Run("False condition boolean test", func(t *testing.T) {
		isActive := false
		status := If(isActive, "Online", "Offline")
		expected := "Offline"
		if status != expected {
			t.Errorf("Expected %s, got %s", expected, status)
		}
	})
}
