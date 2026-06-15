package main

import (
	"fmt"
	"os"
)

// Config holds the application configuration.
type Config struct {
	Host string
	Port int
}

// NewConfig creates a new Config with default values.
func NewConfig() *Config {
	return &Config{
		Host: "localhost",
		Port: 8080,
	}
}

// Person represents a person.
type Person struct {
	Name string
	Age  int
}

// Print prints the person's details.
func (p *Person) Print() {
	fmt.Printf("Name: %s, Age: %d\n", p.Name, p.Age)
}

// Greet returns a greeting string.
func Greet(name string) string {
	return fmt.Sprintf("Hello, %s!", name)
}

func main() {
	cfg := NewConfig()
	fmt.Println(cfg.Host)

	person := Person{Name: "Alice", Age: 30}
	person.Print()

	if os.Getenv("DEBUG") != "" {
		fmt.Println("Debug mode enabled")
	}

	fmt.Println(Greet("Bob"))
}
