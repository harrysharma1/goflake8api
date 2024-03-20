package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	md "github.com/JohannesKaufmann/html-to-markdown"
)

type Body struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Content string   `json:"content"`
	Links   []string `json:"links"`
}

func api(enable bool, code string)  {
	client := &http.Client{}
	endpoint := "https://www.flake8rules.com/api/rules/"+code+"/"
	req, err := http.NewRequest("GET", endpoint, nil)

	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	if !enable{
		normalPrint(bodyText)
	}else{
		prettyPrint(bodyText)
	}
}


func normalPrint(text [] uint8) {
	fmt.Printf("%s\n",text)
}

func prettyPrint(text[] uint8) {
	body := Body{}
	err := json.Unmarshal(text, &body)

	if err != nil {
		fmt.Println(err)
		return
	}
	converter := md.NewConverter("", true, nil)

	fmt.Println(body.Code)
	fmt.Println(body.Message)

	markdown, err:= converter.ConvertString(body.Content)

	if err!=nil{
		log.Fatal(err)
	}
	fmt.Println(markdown)
	for _,i := range body.Links{
		fmt.Println(i)
	}
}

func main()  {
	
    errorCmd := flag.NewFlagSet("error", flag.ExitOnError)
    errorPrettyEnable := errorCmd.Bool("p",false, "Enable pretty print")
    errorCode := errorCmd.String("code","", "Error code for flake8")

    if len(os.Args) < 2{
        fmt.Println("expected 'error' subcommand")
        os.Exit(1)
    }

    switch os.Args[1]{
        case "error":
                errorCmd.Parse(os.Args[2:])
                api(*errorPrettyEnable, *errorCode)
        default:
           fmt.Println("expected 'error' subcommand")
            os.Exit(1)


	}
}

