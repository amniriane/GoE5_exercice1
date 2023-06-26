//package main
//import "fmt"






package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Entrez votre nom : ")
	nom, _ := reader.ReadString('\n')
	fmt.Printf("Bonjour, %s", nom)
}


//func main() {


	
//	var val string
//	var val1 int
//	var val2 float32
//	val3 := "pouet"
//	var val4 bool

//	fmt.Println(val)
//	fmt.Println(val1)
//	fmt.Println(val2)
//	fmt.Println(val3)
//	fmt.Println(val4)
//}