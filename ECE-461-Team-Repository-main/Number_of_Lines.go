package main
import (
	"fmt"
	"io/ioutil"
	"strings"
	"os"
	"strconv"
)

func main(){
	files, err := ioutil.ReadDir("./")
	if err != nil {
		panic(err)
	}
	file, err := os.Create("output_lines.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	for _, f:= range files {
		fmt.Println(f.Name())
		content, err := ioutil.ReadFile(f.Name())
		content1 := string(content)
		if err != nil {
			fmt.Println("Error reading file:", err)
			continue
		}
		lines := strings.Split(content1,"\n")
		nonEmpty := []string{}
		for _, str := range lines{
			ex := ([]rune(str))
			ex1 := 13
			if len(ex) != 0{
				ex1 = int(ex[0])
			}
			if ex1 != 13 || len(ex) != 1{
				nonEmpty = append(nonEmpty,str)
			}
		}
		_, err = file.WriteString(f.Name())
		if err != nil {
			fmt.Println(err)
		}
		string_len := strconv.Itoa(len(nonEmpty)) + " \n"
		_, err = file.WriteString(string_len)
		if err != nil {
			fmt.Println(err)
		}
		// fmt.Println(len(nonEmpty))
	}
}

