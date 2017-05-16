package test
import (
	"testing"
	"net/http"
	. "legitlab.letv.cn/yig/iam/api/datatype"
	"github.com/bmizerany/assert"
	"encoding/json"
	"github.com/bitly/go-simplejson"
	"strings"
	"io/ioutil"
	"fmt"
)
var root_token string
var account_token string
var user_token string
var accountId string
var projectId string
var accessKey string
var accessSecret string

const (
	USER_NAME = "test1"
	USER_PASSWORD = "123456"
	ACCOUNT_NAME = "15579423@qq.com"
	ACCOUNT_PASSWORD = "123456"
	SERVICE_NAME = "S3"
)

func Test_ConnectService(t *testing.T) {
	client := &http.Client{}
	var query QueryRequest
	query.Action = ACTION_ConnectService
	query.UserName = "root"
	query.Password = "admin"
	body, err := json.Marshal(query)
	if err != nil {
		t.Error("Test_ConnectService error", err)
	}
	req, _ := http.NewRequest("POST", "http://127.0.0.1:8888/", strings.NewReader(string(body)))
	response, err := client.Do(req)
	if err != nil {
		t.Error("Test_ConnectService error", err)
	}
	if response.StatusCode != 200 {
		t.Error("Test_ConnectService error", err)
	}
	res_body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Error("Test_ConnectService error", err)
	}
	fmt.Println("res_body", string(res_body))
	js, _ := simplejson.NewJson(res_body)
	retCode, err := js.Get("retCode").Int()
	if retCode != 0 {
		t.Error("Test_ConnectService response not 0 :", retCode)
	} else {
		root_token, err = js.Get("data").Get("token").String()
		assert.Equal(t, nil, err)
	}
}

func Test_CreateAccount(t *testing.T) {
	client := &http.Client{}
	var query QueryRequest
	query.Action = ACTION_CreateAccount
	query.Token = root_token
	query.UserName = ACCOUNT_NAME
	query.Password = ACCOUNT_PASSWORD
	body, err := json.Marshal(query)
	if err != nil {
		t.Error("Test_CreateAccount error", err)
	}
	req, _ := http.NewRequest("POST", "http://127.0.0.1:8888/", strings.NewReader(string(body)))
	response, err := client.Do(req)
	if err != nil {
		t.Error("Test_CreateAccount error", err)
	}
	if response.StatusCode != 200 {
		t.Error("Test_CreateAccount error", err)
	}
	res_body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Error("Test_CreateAccount error", err)
	}

	js, _ := simplejson.NewJson(res_body)
	retCode, err := js.Get("retCode").Int()
	if retCode != 0 {
		t.Error("Test_ConnectService response not 0 :", retCode)
	}
}

func Test_ListAccounts(t *testing.T) {
	client := &http.Client{}
	var query QueryRequest
	query.Action = ACTION_ListAccounts
	query.Token = root_token
	body, err := json.Marshal(query)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}
	req, _ := http.NewRequest("POST", "http://127.0.0.1:8888/", strings.NewReader(string(body)))
	response, err := client.Do(req)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}
	if response.StatusCode != 200 {
		t.Error("Test_ListAccounts error", err)
	}
	res_body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}

	js, err := simplejson.NewJson(res_body)
	assert.NotEqual(t, nil, js)
	assert.Equal(t, nil, err)
	retCode, err := js.Get("retCode").Int()
	assert.Equal(t, nil, err)
	assert.Equal(t, 0, retCode)
	fmt.Println("response", string(res_body))
	accountId, err = js.Get("data").GetIndex(0).Get("AccountId").String()
	assert.NotEqual(t, nil, accountId)
	fmt.Println("accountId=", accountId)

}

func Test_DescribeAccount(t *testing.T) {
	client := &http.Client{}
	var query QueryRequest
	query.Action = ACTION_DescribeAccount
	query.Token = root_token
	query.AccountId = accountId
	body, err := json.Marshal(query)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}
	req, _ := http.NewRequest("POST", "http://127.0.0.1:8888/", strings.NewReader(string(body)))
	response, err := client.Do(req)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}
	if response.StatusCode != 200 {
		t.Error("Test_ListAccounts error", err)
	}
	res_body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}

	js, _ := simplejson.NewJson(res_body)
	retCode, err := js.Get("RetCode").Int()
	if retCode != 0 {
		t.Error("Test_ListAccounts response not 0 :", retCode)
	} else {
		data_map, err := js.Map()
		if err != nil {
			t.Error("Test_ConnectService parase body :", err)
		}
		fmt.Println("data:", data_map)
	}
}

func Test_AccountLogin(t *testing.T) {
	client := &http.Client{}
	var query QueryRequest
	query.Action = ACTION_ConnectService
	query.UserName = ACCOUNT_NAME
	query.Password = ACCOUNT_PASSWORD
	body, err := json.Marshal(query)
	if err != nil {
		t.Error("Test_ConnectService error", err)
	}
	req, _ := http.NewRequest("POST", "http://127.0.0.1:8888/", strings.NewReader(string(body)))
	response, err := client.Do(req)
	if err != nil {
		t.Error("Test_ConnectService error", err)
	}
	if response.StatusCode != 200 {
		t.Error("Test_ConnectService error", err)
	}
	res_body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Error("Test_ConnectService error", err)
	}

	js, _ := simplejson.NewJson(res_body)
	retCode, err := js.Get("retCode").Int()
	if retCode != 0 {
		t.Error("Test_ConnectService response not 0 :", retCode)
	} else {
		account_token, err = js.Get("data").Get("token").String()
		if err != nil {
			t.Error("Test_ConnectService parase body :", err)
		}
	}
	fmt.Println("account login success")
}

func Test_CreateUser(t *testing.T) {
	client := &http.Client{}
	var query QueryRequest
	query.Action = ACTION_CreateUser
	query.Token = account_token
	query.UserName = USER_NAME
	query.Password = USER_PASSWORD
	body, err := json.Marshal(query)
	if err != nil {
		t.Error("Test_CreateUser error", err)
	}
	req, _ := http.NewRequest("POST", "http://127.0.0.1:8888/", strings.NewReader(string(body)))
	response, err := client.Do(req)
	if err != nil {
		t.Error("Test_CreateUser error", err)
	}
	if response.StatusCode != 200 {
		t.Error("Test_CreateUser error", err)
	}
	res_body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Error("Test_CreateUser error", err)
	}

	js, _ := simplejson.NewJson(res_body)
	retCode, err := js.Get("retCode").Int()
	if retCode != 0 {
		t.Error("Test_CreateUser response not 0 :", retCode)
	}
	fmt.Println("user create success")
}

func Test_UserLogin(t *testing.T) {
	client := &http.Client{}
	var query QueryRequest
	query.Action = ACTION_ConnectService
	query.UserName = accountId + ":" + USER_NAME
	fmt.Println("User login:", query.UserName)
	query.Password = USER_PASSWORD
	body, err := json.Marshal(query)
	if err != nil {
		t.Error("Test_ConnectService error", err)
	}
	req, _ := http.NewRequest("POST", "http://127.0.0.1:8888/", strings.NewReader(string(body)))
	response, err := client.Do(req)
	if err != nil {
		t.Error("Test_ConnectService error", err)
	}
	if response.StatusCode != 200 {
		t.Error("Test_ConnectService error", err)
	}
	res_body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Error("Test_ConnectService error", err)
	}

	js, _ := simplejson.NewJson(res_body)
	retCode, err := js.Get("retCode").Int()
	if retCode != 0 {
		t.Error("Test_ConnectService response not 0 :", retCode)
	} else {
		user_token, err = js.Get("data").Get("token").String()
		if err != nil {
			t.Error("Test_ConnectService parase body :", err)
		}
	}

	fmt.Println("user login success", string(res_body))
}

func Test_DescribeUser_From_User(t *testing.T) {
	client := &http.Client{}
	var query QueryRequest
	query.Action = ACTION_DescribeUser
	query.Token = user_token
	body, err := json.Marshal(query)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}
	req, _ := http.NewRequest("POST", "http://127.0.0.1:8888/", strings.NewReader(string(body)))
	response, err := client.Do(req)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}
	if response.StatusCode != 200 {
		t.Error("Test_ListAccounts error", err)
	}
	res_body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}

	js, _ := simplejson.NewJson(res_body)
	retCode, err := js.Get("retCode").Int()
	assert.Equal(t, 0, retCode)
	fmt.Println("user describe user success", string(res_body))
}

func Test_DescribeUser_From_Account(t *testing.T) {
	client := &http.Client{}
	var query QueryRequest
	query.Action = ACTION_DescribeUser
	query.Token = account_token
	query.UserName = USER_NAME
	body, err := json.Marshal(query)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}
	req, _ := http.NewRequest("POST", "http://127.0.0.1:8888/", strings.NewReader(string(body)))
	response, err := client.Do(req)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}
	if response.StatusCode != 200 {
		t.Error("Test_ListAccounts error", err)
	}
	res_body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}

	js, _ := simplejson.NewJson(res_body)
	retCode, err := js.Get("retCode").Int()
	if retCode != 0 {
		t.Error("Test_ListAccounts response not 0 :", retCode)
	}
	fmt.Println("account describe user success")
}

func Test_ListUsers(t *testing.T) {
	client := &http.Client{}
	var query QueryRequest
	query.Action = ACTION_ListUsers
	query.Token = account_token
	body, err := json.Marshal(query)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}
	req, _ := http.NewRequest("POST", "http://127.0.0.1:8888/", strings.NewReader(string(body)))
	response, err := client.Do(req)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}
	if response.StatusCode != 200 {
		t.Error("Test_ListAccounts error", err)
	}
	res_body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}

	js, _ := simplejson.NewJson(res_body)
	retCode, err := js.Get("retCode").Int()
	if retCode != 0 {
		t.Error("Test_ListAccounts response not 0 :", retCode)
	}
	userName, err := js.Get("data").GetIndex(0).Get("UserName").String()
	assert.NotEqual(t, nil, userName)
	fmt.Println("userName=", userName)
	fmt.Println("list user success")
}

func Test_CreateProject(t *testing.T) {
	client := &http.Client{}
	var query QueryRequest
	query.Action = ACTION_CreateProject
	query.Token = account_token
	query.ProjectName = "testproject"
	query.Description = "thisistestproject"
	body, err := json.Marshal(query)
	if err != nil {
		t.Error("Test_CreateProject error", err)
	}
	fmt.Println("Test_CreateProject body:", string(body))
	req, err := http.NewRequest("POST", "http://127.0.0.1:8888/", strings.NewReader(string(body)))
	assert.Equal(t, nil, err)
	response, err := client.Do(req)
	if err != nil {
		t.Error("Test_CreateProject error", err)
	}
	if response.StatusCode != 200 {
		t.Error("Test_CreateProject error", err)
	}
	res_body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Error("Test_CreateProject error", err)
	}

	js, _ := simplejson.NewJson(res_body)
	retCode, err := js.Get("retCode").Int()
	if retCode != 0 {
		t.Error("Test_CreateProject response not 0 :", retCode)
	}
	assert.Equal(t, 0, retCode)
	fmt.Println("create project success")
}

func Test_ListProjects(t *testing.T) {
	client := &http.Client{}
	var query QueryRequest
	query.Action = ACTION_ListProjects
	query.Token = account_token
	body, err := json.Marshal(query)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}
	req, _ := http.NewRequest("POST", "http://127.0.0.1:8888/", strings.NewReader(string(body)))
	response, err := client.Do(req)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}
	if response.StatusCode != 200 {
		t.Error("Test_ListAccounts error", err)
	}
	res_body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}
	fmt.Println("res_body=", string(res_body))
	js, _ := simplejson.NewJson(res_body)
	retCode, err := js.Get("retCode").Int()
	if retCode != 0 {
		t.Error("Test_ListAccounts response not 0 :", retCode)
	}
	projectId, err = js.Get("data").GetIndex(0).Get("projectId").String()
	assert.NotEqual(t, nil, projectId)
	fmt.Println("ProjectId=", projectId)
	fmt.Println("list project success")
}

func Test_DescribeProject(t *testing.T) {
	client := &http.Client{}
	var query QueryRequest
	query.Action = ACTION_DescribeProject
	query.Token = account_token
	query.ProjectId = projectId
	body, err := json.Marshal(query)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}
	req, _ := http.NewRequest("POST", "http://127.0.0.1:8888/", strings.NewReader(string(body)))
	response, err := client.Do(req)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}
	if response.StatusCode != 200 {
		t.Error("Test_ListAccounts error", err)
	}
	res_body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}

	js, _ := simplejson.NewJson(res_body)
	retCode, err := js.Get("retCode").Int()
	if retCode != 0 {
		t.Error("Test_ListAccounts response not 0 :", retCode)
	}
	data, err := js.Get("data").Map()
	fmt.Println("project=", data)
	fmt.Println("describe project success")
}

func Test_LinkUserWithProject(t *testing.T) {
	client := &http.Client{}
	var query QueryRequest
	query.Action = ACTION_LinkUserWithProject
	query.Token = account_token
	query.ProjectId = projectId
	query.UserName = USER_NAME
	body, err := json.Marshal(query)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}
	req, _ := http.NewRequest("POST", "http://127.0.0.1:8888/", strings.NewReader(string(body)))
	response, err := client.Do(req)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}
	if response.StatusCode != 200 {
		t.Error("Test_ListAccounts error", err)
	}
	res_body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}
	fmt.Println("res_body=", string(res_body))
	js, _ := simplejson.NewJson(res_body)
	retCode, err := js.Get("retCode").Int()
	if retCode != 0 {
		t.Error("Test_ListAccounts response not 0 :", retCode)
	}
	fmt.Println("link user project success")
}

func Test_ListProjectByUser_USER_TOKEN(t *testing.T) {
	client := &http.Client{}
	var query QueryRequest
	query.Action = ACTION_ListProjectByUser
	query.Token = user_token
	body, err := json.Marshal(query)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}
	req, _ := http.NewRequest("POST", "http://127.0.0.1:8888/", strings.NewReader(string(body)))
	response, err := client.Do(req)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}
	if response.StatusCode != 200 {
		t.Error("Test_ListAccounts error", err)
	}
	res_body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}
	fmt.Println("res_body=", string(res_body))
	js, _ := simplejson.NewJson(res_body)
	retCode, err := js.Get("retCode").Int()
	if retCode != 0 {
		t.Error("Test_ListAccounts response not 0 :", retCode)
	}
	ProjectId, err := js.Get("data").GetIndex(0).Get("projectId").String()
	assert.NotEqual(t, nil, ProjectId)
	fmt.Println("ProjectName=", ProjectId)
	fmt.Println("list project by user from user success")
}

func Test_ListProjectByUser_ACCCOUNT_TOKEN(t *testing.T) {
	client := &http.Client{}
	var query QueryRequest
	query.Action = ACTION_ListProjectByUser
	query.Token = account_token
	query.UserName = USER_NAME
	body, err := json.Marshal(query)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}
	req, _ := http.NewRequest("POST", "http://127.0.0.1:8888/", strings.NewReader(string(body)))
	response, err := client.Do(req)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}
	if response.StatusCode != 200 {
		t.Error("Test_ListAccounts error", err)
	}
	res_body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}
	fmt.Println("res_body=", string(res_body))
	js, _ := simplejson.NewJson(res_body)
	retCode, err := js.Get("retCode").Int()
	if retCode != 0 {
		t.Error("Test_ListAccounts response not 0 :", retCode)
	}
	ProjectId, err := js.Get("data").GetIndex(0).Get("projectId").String()
	assert.NotEqual(t, nil, ProjectId)
	fmt.Println("projectId=", ProjectId)
	fmt.Println("list project by user from account success")
}

func Test_ListUserByProject(t *testing.T) {
	client := &http.Client{}
	var query QueryRequest
	query.Action = ACTION_ListUserByProject
	query.Token = account_token
	query.ProjectId = projectId
	body, err := json.Marshal(query)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}
	req, _ := http.NewRequest("POST", "http://127.0.0.1:8888/", strings.NewReader(string(body)))
	response, err := client.Do(req)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}
	if response.StatusCode != 200 {
		t.Error("Test_ListAccounts error", err)
	}
	res_body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}
	fmt.Println("res_body=", string(res_body))
	js, _ := simplejson.NewJson(res_body)
	retCode, err := js.Get("retCode").Int()
	if retCode != 0 {
		t.Error("Test_ListAccounts response not 0 :", retCode)
	}
	UserName, err := js.Get("data").GetIndex(0).Get("userName").String()
	assert.NotEqual(t, nil, UserName)
	fmt.Println("userName=", UserName)
	fmt.Println("list user by project success")
}

func Test_AddProjectService(t *testing.T) {
	client := &http.Client{}
	var query QueryRequest
	query.Action = ACTION_AddProjectService
	query.Token = account_token
	query.ProjectId = projectId
	query.Service = SERVICE_NAME
	body, err := json.Marshal(query)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}
	req, _ := http.NewRequest("POST", "http://127.0.0.1:8888/", strings.NewReader(string(body)))
	response, err := client.Do(req)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}
	if response.StatusCode != 200 {
		t.Error("Test_ListAccounts error", err)
	}
	res_body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}
	fmt.Println("res_body=", string(res_body))
	js, _ := simplejson.NewJson(res_body)
	retCode, err := js.Get("retCode").Int()
	if retCode != 0 {
		t.Error("Test_ListAccounts response not 0 :", retCode)
	}
	fmt.Println("add project service success")
}

func Test_ListServiceByProject_ACCOUNT_TOKEN(t *testing.T) {
	client := &http.Client{}
	var query QueryRequest
	query.Action = ACTION_ListServiceByProject
	query.Token = account_token
	query.ProjectId = projectId
	body, err := json.Marshal(query)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}
	req, _ := http.NewRequest("POST", "http://127.0.0.1:8888/", strings.NewReader(string(body)))
	response, err := client.Do(req)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}
	if response.StatusCode != 200 {
		t.Error("Test_ListAccounts error", err)
	}
	res_body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}
	fmt.Println("res_body=", string(res_body))
	js, _ := simplejson.NewJson(res_body)
	retCode, err := js.Get("retCode").Int()
	if retCode != 0 {
		t.Error("Test_ListAccounts response not 0 :", retCode)
	}
	Service, err := js.Get("data").GetIndex(0).Get("service").String()
	assert.NotEqual(t, nil, Service)
	fmt.Println("Service=", Service)
	fmt.Println("list user by project success")
}

func Test_ListServiceByProject_USER_TOKEN(t *testing.T) {
	client := &http.Client{}
	var query QueryRequest
	query.Action = ACTION_ListServiceByProject
	query.Token = user_token
	query.ProjectId = projectId
	body, err := json.Marshal(query)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}
	req, _ := http.NewRequest("POST", "http://127.0.0.1:8888/", strings.NewReader(string(body)))
	response, err := client.Do(req)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}
	if response.StatusCode != 200 {
		t.Error("Test_ListAccounts error", err)
	}
	res_body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}
	fmt.Println("res_body=", string(res_body))
	js, _ := simplejson.NewJson(res_body)
	retCode, err := js.Get("retCode").Int()
	if retCode != 0 {
		t.Error("Test_ListAccounts response not 0 :", retCode)
	}
	Service, err := js.Get("data").GetIndex(0).Get("service").String()
	assert.NotEqual(t, nil, Service)
	fmt.Println("Service=", Service)
	fmt.Println("list user by project success")
}

func Test_CreateAccessKey(t *testing.T) {
	client := &http.Client{}
	var query QueryRequest
	query.Action = ACTION_CreateAccessKey
	query.Token = account_token
	query.ProjectId = projectId
	query.KeyName = "hehe"
	body, err := json.Marshal(query)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}
	req, _ := http.NewRequest("POST", "http://127.0.0.1:8888/", strings.NewReader(string(body)))
	response, err := client.Do(req)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}
	if response.StatusCode != 200 {
		t.Error("Test_ListAccounts error", err)
	}
	res_body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}
	fmt.Println("res_body=", string(res_body))
	js, _ := simplejson.NewJson(res_body)
	retCode, err := js.Get("retCode").Int()
	if retCode != 0 {
		t.Error("Test_ListAccounts response not 0 :", retCode)
	}
	fmt.Println("create ak sk success")
}

func Test_ListAccessKeysByProject(t *testing.T) {
	client := &http.Client{}
	var query QueryRequest
	query.Action = ACTION_ListAccessKeysByProject
	query.Token = account_token
	query.ProjectId = projectId
	body, err := json.Marshal(query)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}
	req, _ := http.NewRequest("POST", "http://127.0.0.1:8888/", strings.NewReader(string(body)))
	response, err := client.Do(req)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}
	if response.StatusCode != 200 {
		t.Error("Test_ListAccounts error", err)
	}
	res_body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}
	fmt.Println("res_body=", string(res_body))
	js, _ := simplejson.NewJson(res_body)
	retCode, err := js.Get("retCode").Int()
	if retCode != 0 {
		t.Error("Test_ListAccounts response not 0 :", retCode)
	}
	accessKey, err = js.Get("data").GetIndex(0).Get("accessKey").String()
	accessSecret, err = js.Get("data").GetIndex(0).Get("accessSecret").String()
	fmt.Println("list ak sk success", accessKey, accessSecret)
}

func Test_DescribeAccessKeys(t *testing.T) {
	client := &http.Client{}
	var query QueryRequest
	query.Action = ACTION_DescribeAccessKeys
	query.AccessKeys = []string{accessKey}
	body, err := json.Marshal(query)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}
	req, _ := http.NewRequest("POST", "http://127.0.0.1:8888/", strings.NewReader(string(body)))
	req.Header.Set("X-Le-Key", "key")
	req.Header.Set("X-Le-Secret", "secret")
	response, err := client.Do(req)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}
	if response.StatusCode != 200 {
		t.Error("Test_ListAccounts error", err)
	}
	res_body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}
	fmt.Println("res_body=", string(res_body))
	js, _ := simplejson.NewJson(res_body)
	retCode, err := js.Get("retCode").Int()
	if retCode != 0 {
		t.Error("Test_ListAccounts response not 0 :", retCode)
	}
	AccessKey, err := js.Get("data").Get("accessKeySet").GetIndex(0).Get("accessKey").String()
	AccessSecret, err := js.Get("data").Get("accessKeySet").GetIndex(0).Get("accessSecret").String()
	fmt.Println("list ak sk success", AccessKey, AccessSecret)
}


func Test_DeleteAccessKey(t *testing.T) {
	client := &http.Client{}
	var query QueryRequest
	query.Action = ACTION_DeleteAccessKey
	query.Token = account_token
	query.ProjectId = projectId
	query.AccessKey = accessKey
	body, err := json.Marshal(query)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}
	req, _ := http.NewRequest("POST", "http://127.0.0.1:8888/", strings.NewReader(string(body)))
	response, err := client.Do(req)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}
	if response.StatusCode != 200 {
		t.Error("Test_ListAccounts error", err)
	}
	res_body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}
	fmt.Println("res_body=", string(res_body))
	js, _ := simplejson.NewJson(res_body)
	retCode, err := js.Get("retCode").Int()
	if retCode != 0 {
		t.Error("Test_ListAccounts response not 0 :", retCode)
	}
	fmt.Println("delete ak sk success")
}


func Test_DelProjectService(t *testing.T) {
	client := &http.Client{}
	var query QueryRequest
	query.Action = ACTION_DelProjectService
	query.Token = account_token
	query.ProjectId = projectId
	query.Service = SERVICE_NAME
	body, err := json.Marshal(query)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}
	req, _ := http.NewRequest("POST", "http://127.0.0.1:8888/", strings.NewReader(string(body)))
	response, err := client.Do(req)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}
	if response.StatusCode != 200 {
		t.Error("Test_ListAccounts error", err)
	}
	res_body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}
	fmt.Println("res_body=", string(res_body))
	js, _ := simplejson.NewJson(res_body)
	retCode, err := js.Get("retCode").Int()
	if retCode != 0 {
		t.Error("Test_ListAccounts response not 0 :", retCode)
	}
	fmt.Println("del project service success")
}

func Test_UnLinkUserWithProject(t *testing.T) {
	client := &http.Client{}
	var query QueryRequest
	query.Action = ACTION_UnLinkUserWithProject
	query.Token = account_token
	body, err := json.Marshal(query)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}
	req, _ := http.NewRequest("POST", "http://127.0.0.1:8888/", strings.NewReader(string(body)))
	response, err := client.Do(req)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}
	if response.StatusCode != 200 {
		t.Error("Test_ListAccounts error", err)
	}
	res_body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}
	fmt.Println("res_body=", string(res_body))
	js, _ := simplejson.NewJson(res_body)
	retCode, err := js.Get("retCode").Int()
	if retCode != 0 {
		t.Error("Test_ListAccounts response not 0 :", retCode)
	}
	fmt.Println("unlink user project success")
}

func Test_DeleteProject(t *testing.T) {
	client := &http.Client{}
	var query QueryRequest
	query.Action = ACTION_DeleteProject
	query.Token = account_token
	query.ProjectId = projectId
	body, err := json.Marshal(query)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}
	req, _ := http.NewRequest("POST", "http://127.0.0.1:8888/", strings.NewReader(string(body)))
	response, err := client.Do(req)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}
	if response.StatusCode != 200 {
		t.Error("Test_ListAccounts error", err)
	}
	res_body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}

	js, _ := simplejson.NewJson(res_body)
	retCode, err := js.Get("retCode").Int()
	if retCode != 0 {
		t.Error("Test_ListAccounts response not 0 :", retCode)
	}

	fmt.Println("delete project success")
}

func Test_DeleteUser(t *testing.T) {
	client := &http.Client{}
	var query QueryRequest
	query.Action = ACTION_DeleteUser
	query.Token = account_token
	query.UserName = USER_NAME
	body, err := json.Marshal(query)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}
	req, _ := http.NewRequest("POST", "http://127.0.0.1:8888/", strings.NewReader(string(body)))
	response, err := client.Do(req)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}
	if response.StatusCode != 200 {
		t.Error("Test_ListAccounts error", err)
	}
	res_body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}

	js, _ := simplejson.NewJson(res_body)
	retCode, err := js.Get("retCode").Int()
	if retCode != 0 {
		t.Error("Test_ListAccounts response not 0 :", retCode)
	}
	fmt.Println("user deleted success")
}

func Test_DeleteAccount(t *testing.T) {
	client := &http.Client{}
	var query QueryRequest
	query.Action = ACTION_DeleteAccount
	query.Token = root_token
	query.AccountId = accountId
	body, err := json.Marshal(query)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}
	req, _ := http.NewRequest("POST", "http://127.0.0.1:8888/", strings.NewReader(string(body)))
	response, err := client.Do(req)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}
	if response.StatusCode != 200 {
		t.Error("Test_ListAccounts error", err)
	}
	res_body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Error("Test_ListAccounts error", err)
	}
	fmt.Println("res_body", string(res_body))
	js, _ := simplejson.NewJson(res_body)
	retCode, err := js.Get("retCode").Int()
	if retCode != 0 {
		t.Error("Test_ListAccounts response not 0 :", retCode)
	}
}
