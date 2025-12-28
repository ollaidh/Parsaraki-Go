package main

var CONFIG Config

func init() {
	CONFIG = loadConfig()
}

func main() {
	pingBot()

}
