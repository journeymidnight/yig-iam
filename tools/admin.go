package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	. "github.com/journeymidnight/yig-iam/api/datatype"
	"io/ioutil"
	"net/http"
	"os"
)

var client = &http.Client{}

type Config struct {
	RequestUrl string
	AdminKey   string
}

var config Config

var endPoint *string
var token *string

func printHelp() {
	fmt.Println("Usage: admin <commands> [options...] ")
	fmt.Println("Commands: login|createaccount|createuser|listaccounts|createproject|listkeys")
	fmt.Println("Options:")
	fmt.Println(" -e, --endpoint   Specify endpoint of yig-iam")
	fmt.Println(" -u, --user      Specify user name to login")
	fmt.Println(" -p, --password   Specify password to login")
	fmt.Println(" -t, --token   Specify token to send request with")
	fmt.Println(" -m, --email   Specify email")
	fmt.Println(" -d, --displayname   Specify display name")
}

func isParaEmpty(p string) bool {
	if p == "" {
		fmt.Printf("Bad usage, Try admin")
		return true
	} else {
		return false
	}
}

func prettyPrint(body []byte) {
	if body == nil || len(body) == 0{
		return
	}
	var prettyJSON bytes.Buffer
	error := json.Indent(&prettyJSON, body, "", "\t")
	if error != nil {
		fmt.Println("bad json format", error.Error())
		return
	}
	fmt.Printf(string(prettyJSON.Bytes()))
}
func login(user , password string) {
	if isParaEmpty(user) ||  isParaEmpty(password){
		return
	}

	query := &QueryRequest{}
	query.UserName = user
	query.Password = password

	blob, _ := json.Marshal(query)

	url := *endPoint + "/api/v1/login"
	request, err := http.NewRequest("POST", url, bytes.NewReader(blob))
	if err != nil {
		fmt.Println("create request failed", err)
		return
	}

	response, err := client.Do(request)
	if err != nil {
		fmt.Println("send request failed", err)
		return
	}
	if response.StatusCode != 200 {
		fmt.Println("login failed as status != 200", response.StatusCode)
		return
	}

	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	prettyPrint(body)
}

func createAccount(user , password, email, displayName string) {
	if isParaEmpty(user) ||  isParaEmpty(password){
		return
	}

	query := &QueryRequest{}
	query.UserName = user
	query.Password = password
	query.Email = email
	query.DisplayName = displayName

	blob, _ := json.Marshal(query)

	url := *endPoint + "/api/v1/account/create"
	request, _ := http.NewRequest("POST", url, bytes.NewReader(blob))
	request.Header.Set("Token", *token)
	response, _ := client.Do(request)
	if response.StatusCode != 200 {
		fmt.Println("createAccount failed as status != 200", response.StatusCode)
		return
	}

	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	prettyPrint(body)
}

func createUser(user , password, email, displayName string) {
	if isParaEmpty(user) ||  isParaEmpty(password){
		return
	}

	query := &QueryRequest{}
	query.UserName = user
	query.Password = password
	query.Email = email
	query.DisplayName = displayName

	blob, _ := json.Marshal(query)

	url := *endPoint + "/api/v1/user/create"
	request, _ := http.NewRequest("POST", url, bytes.NewReader(blob))
	request.Header.Set("Token", *token)
	response, _ := client.Do(request)
	if response.StatusCode != 200 {
		fmt.Println("createUser failed as status != 200", response.StatusCode)
		return
	}

	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	prettyPrint(body)
}

func listAccount() {
	url := *endPoint + "/api/v1/account/list"
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Set("Token", *token)
	response, _ := client.Do(request)
	if response.StatusCode != 200 {
		fmt.Println("listAccount failed as status != 200", response.StatusCode)
		return
	}

	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	prettyPrint(body)
}

func createProject(projectName, description string) {
	if isParaEmpty(projectName) {
		return
	}

	query := &QueryRequest{}
	query.ProjectName = projectName
	query.Description = description

	blob, _ := json.Marshal(query)

	url := *endPoint + "/api/v1/project/create"
	request, _ := http.NewRequest("POST", url, bytes.NewReader(blob))
	request.Header.Set("Token", *token)
	response, _ := client.Do(request)
	if response.StatusCode != 200 {
		fmt.Println("createProject failed as status != 200", response.StatusCode)
		return
	}

	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	prettyPrint(body)
}

func listkeys() {
	url := *endPoint + "/api/v1/keys/list"
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Set("Token", *token)
	response, _ := client.Do(request)
	if response.StatusCode != 200 {
		fmt.Println("listkeys failed as status != 200", response.StatusCode)
		return
	}

	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	prettyPrint(body)

}

func main() {
	if len(os.Args) <= 1 {
		printHelp()
		return
	}
	mySet := flag.NewFlagSet("", flag.ExitOnError)
	user := mySet.String("u", "", "user")
	password := mySet.String("p", "", "password")
	endPoint = mySet.String("e", "http://127.0.0.1:8888", "endpoint")
	token = mySet.String("t", "", "token")
	email := mySet.String("m", "", "email")
	displayName := mySet.String("d", "", "display name")
	projectName := mySet.String("pn", "", "project name")
	description := mySet.String("pd", "", "project description")
	mySet.Parse(os.Args[2:])
	fmt.Println("command:", os.Args[1], "user:", *user, "password:", *password)
	fmt.Println("endPoint:", *endPoint, "token:", *token)
	switch os.Args[1] {
	case "login":
		login(*user, *password)
	case "createaccount":
		createAccount(*user, *password, *email, *displayName)
	case "createuser":
		createUser(*user, *password, *email, *displayName)
	case "listaccounts":
		listAccount()
	case "createproject":
		createProject(*projectName, *description)
	case "listkeys":
		listkeys()

	default:
		printHelp()
		return
	}
}
