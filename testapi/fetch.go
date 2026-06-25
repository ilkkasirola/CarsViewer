package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type Car struct {
	Image string `json:image`
}

const (
	urlPrefix string = "http://"
	urlSuffix string = ".hive.fi:3000/"

	colorRed   string = "\u001B[31m"
	colorReset string = "\033[0m"
)

var url string

func main() {
	if len(os.Args) != 2 {
		fmt.Println(colorRed + "Need computer ID as argument!\neg. c1r1p1" + colorReset)
		os.Exit(1)
	} else {
		url = urlPrefix + os.Args[1] + urlSuffix
		fmt.Println("Fetching data from", url)
	}

	os.Mkdir("api", 0750)

	fetchImages()
	fetchJson()

	fmt.Println("Program executed successfully.")
}

func fetchImages() {
	var cars []Car
	resp, err := http.Get(url + "api/models")
	if err != nil {
		printErr("http.Get ERR:", err)
		os.Exit(1)
	}

	defer resp.Body.Close()
	//fmt.Println(url + "api/models")

	if resp.StatusCode != 200 {
		fmt.Println(resp.StatusCode)
		printErr("Address cannot be reached!", err)
		os.Exit(1)
	}

	err = json.NewDecoder(resp.Body).Decode(&cars)
	if err != nil {
		printErr("json.Decode ERR:", err)
	}

	for i, car := range cars {
		cmd := exec.Command("wget", "-P", "api/img/", url+"api/images/"+car.Image)
		err := cmd.Run()

		if i%(len(cars)/10) == 0 || i == len(cars)-1 {
			progress := fmt.Sprintf("Progress: %.2f%%", float64(i)/float64(len(cars))*100.0)
			fmt.Println(progress)
		}

		if err != nil {
			printErr("cmd.Run ERR:", err)
		}
	}

	fmt.Println("Images saved to api/img/")
}

func fetchJson() {
	m := map[string]string{
		"models":        "",
		"manufacturers": "",
		"categories":    "",
	}

	for key, _ := range m {
		resp, err := http.Get(url + "api/" + key)
		defer resp.Body.Close()

		if err != nil {
			printErr("http.Get ERR:", err)
		}

		buf := new(strings.Builder)
		n, err := io.Copy(buf, resp.Body)
		if err != nil {
			printErr("io.Copy ERR:", err)
			fmt.Println(n)
		}

		m[key] = buf.String()
	}

	var b strings.Builder

	b.WriteString("{")

	b.WriteString(`"carModels": ` + m["models"] + ",")
	b.WriteString(`"manufacturers": ` + m["manufacturers"] + ",")
	b.WriteString(`"categories": ` + m["categories"])

	b.WriteString("}")

	err := os.WriteFile("api/data.json", []byte(b.String()), 0666)
	if err != nil {
		printErr("os.WriteFile ERR::", err)
	} else {
		fmt.Println("Data written to api/data.json")
	}
}

func printErr(msg string, err error) {
	fmt.Println(colorRed+msg, err, colorReset)
}
