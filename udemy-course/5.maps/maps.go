package main

import "fmt"

func main() {
	colors := map[string]string{
		"red":   "#ff0000",
		"green": "#00ff00",
		"white": "#ffffff",
	}
	fmt.Println(colors)
	printMap(colors)

	var emptyColors map[string]string
	fmt.Println(emptyColors)

	anotherMap := make(map[string]string)
	fmt.Println(anotherMap)
	anotherMap["white"] = "#ffffff"
	fmt.Println(anotherMap)
	delete(anotherMap, "white")
	fmt.Println(anotherMap)
}

func printMap(m map[string]string) {
	for key, value := range m {
		fmt.Printf("%+v = %+v \n", key, value)
	}
}
