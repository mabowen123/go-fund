package main

func main() {
	r := LoadRouter()
	r.Run(":9090")
}
