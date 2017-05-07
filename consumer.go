package main

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	resp, err := http.Get("http://localhost:3000/tax/" + os.Args[1:][0]+ "/" + os.Args[1:][1])

	if (err != nil)	{
		fmt.Println(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	fmt.Print(string(body));
}
