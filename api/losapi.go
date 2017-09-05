package api

import (
	"os"
	"fmt"
	"io"
	"crypto/md5"
	"time"
	"errors"
	"strings"
	"io/ioutil"
	"net/http"
	"gopkg.in/iris.v4"
	"encoding/base64"
	. "github.com/journeymidnight/yig-iam/api/datatype"
	"github.com/journeymidnight/yig-iam/helper"
        "github.com/journeymidnight/yig-iam/db"
	"github.com/minio/minio-go"	
)


var minioClient *minio.Client

type GetS3Domain struct {
	S3Domain string `json:"s3Domain"`
}

func LosApiHandler(c *iris.Context) {
	query := c.Get("queryRequest").(QueryRequest)
	switch query.Action {
		case ACTION_LOS_GETS3DOMAIN:
			getS3Domain(c, query)
		case ACTION_LOS_LISTBUCKETS:
			listBuckets(c, query)
		case ACTION_LOS_DELETEBUCKET:
			deleteBucket(c, query)
		case ACTION_LOS_GETBUCKETSTATS:
			getBucketStats(c, query)
		case ACTION_LOS_CREATEBUCKET:
			createBucket(c, query)
		case ACTION_LOS_PUTCORS:
			putCors(c, query)

		default:
			c.JSON(iris.StatusOK, QueryResponse{RetCode:0,Message:"unsupport action",Data:""})
			return
		}
}

func getS3Domain(c *iris.Context, query QueryRequest) {
        tokenRecord := c.Get("token").(TokenRecord)
        if helper.Enforcer.Enforce(tokenRecord.UserName, API_LOS_GetS3Domain, ACT_ACCESS) != true {
                c.JSON(iris.StatusOK, QueryResponse{RetCode: 4030, Message: "You do not have permission to perform", Data: query})
                return
        }

        if query.Endpoint == "" {
                c.JSON(iris.StatusOK, QueryResponse{RetCode: 4000, Message: "please provide an endpoint to verify", Data: query})
                return
        }

        services, err := db.ListSerivceRecords()
        if err != nil {
                c.JSON(iris.StatusOK, QueryResponse{RetCode: 4030, Message: "failed to get all suppport services", Data: query})
                return
        }

        found := false
        for _, s := range services {
                if s.Endpoint == query.Endpoint {
                        found = true
                        break
                }
        }

        if found == true {
                var data GetS3Domain
                data.S3Domain = "http://" + query.Endpoint
                c.JSON(iris.StatusOK, QueryResponse{RetCode: 0, Message: "", Data: data})
                return
        } else {
                c.JSON(iris.StatusOK, QueryResponse{RetCode: 4030, Message: "your s3domain is not supported", Data: query})
                return
        }

}

func listBuckets(c *iris.Context, query QueryRequest) {
	ak := c.RequestHeader("X-Le-Key")
	sk := c.RequestHeader("X-Le-Secret")
	if ak == "" || sk == "" {
		c.JSON(iris.StatusOK, QueryResponse{RetCode:4101,Message:"invalid accessKey and secretKay",Data:""})
		return
	}
        tokenRecord := c.Get("token").(TokenRecord)
        if helper.Enforcer.Enforce(tokenRecord.UserName, API_LOS_GetS3Domain, ACT_ACCESS) != true {
                c.JSON(iris.StatusOK, QueryResponse{RetCode: 4030, Message: "You do not have permission to perform", Data: query})
                return
        }

        if query.Endpoint == "" {
                c.JSON(iris.StatusOK, QueryResponse{RetCode: 4000, Message: "please provide an endpoint to verify", Data: query})
                return
        }

	client, err := minio.NewV2(query.Endpoint, ak, sk, false)
	if err != nil {
		c.JSON(iris.StatusOK, QueryResponse{RetCode:4000,Message:err.Error(),Data:""})
		return
	}
	BucketsInfo, err := client.ListBuckets()
	if err != nil {
		helper.Logger.Println(5, "listbucket error", err.Error(), ak, sk)
		c.JSON(iris.StatusOK, QueryResponse{RetCode:4101,Message:err.Error(),Data:""})
	} else {
		c.JSON(iris.StatusOK, QueryResponse{RetCode:0,Message:"ok",Data:BucketsInfo})
	}

}

func deleteBucket(c *iris.Context, query QueryRequest) {
	ak := c.RequestHeader("X-Le-Key")
	sk := c.RequestHeader("X-Le-Secret")
	if ak == "" || sk == "" {
		c.JSON(iris.StatusOK, QueryResponse{RetCode:4101,Message:"invalid accessKey and secretKay",Data:""})
		return
	}

        tokenRecord := c.Get("token").(TokenRecord)
        if helper.Enforcer.Enforce(tokenRecord.UserName, API_LOS_GetS3Domain, ACT_ACCESS) != true {
                c.JSON(iris.StatusOK, QueryResponse{RetCode: 4030, Message: "You do not have permission to perform", Data: query})
                return
        }

        if query.Endpoint == "" {
                c.JSON(iris.StatusOK, QueryResponse{RetCode: 4000, Message: "please provide an endpoint to verify", Data: query})
                return
        }
	client, err := minio.NewV2(query.Endpoint, ak, sk, false)
	if err != nil {
		helper.Logger.Fatalln(5, err)
	}
	err = client.RemoveBucket(query.Bucket)
	if err != nil {
		helper.Logger.Println(5, "deleteBucket error", err.Error(), ak, sk)
		c.JSON(iris.StatusOK, QueryResponse{RetCode:4101,Message:err.Error(),Data:""})
	} else {
		c.JSON(iris.StatusOK, QueryResponse{RetCode:0,Message:"bucket removed",Data:""})
	}

	return
}

func getBucketStats(c *iris.Context, query QueryRequest) {
	if query.Bucket == "" {
		c.JSON(iris.StatusOK, QueryResponse{RetCode:4000,Message:"please provide a bucketname", Data:""})
		return
	}

	Client := &http.Client{Timeout: time.Second * 5}
	url := "http://" + query.Endpoint + "/admin/bucket?format=json&bucket=" + query.Bucket +"&stats=False"
	method := "GET"

//	slog.Println("new request to s3:" + url)
	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		helper.Logger.Println(5, err.Error())
		c.JSON(iris.StatusOK, QueryResponse{RetCode:4000,Message:err.Error(), Data:""})
		return
	}
	request_p := helper.SignV2(*request, helper.CONFIG.AccessKey, helper.CONFIG.SecretKey)
	response, err := Client.Do(request_p)
	if err != nil {
		helper.Logger.Println(5, err.Error())
		c.JSON(iris.StatusOK, QueryResponse{RetCode:4000,Message:err.Error(), Data:""})
		return
	}
	defer response.Body.Close()
//	stdout := os.Stdout
//	_, err = io.Copy(stdout, response.Body)
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		helper.Logger.Println(5, err.Error())
		c.JSON(iris.StatusOK, QueryResponse{RetCode:4000,Message:err.Error(),Data:""})
		return
	}

	if response.StatusCode != 200 {
		c.JSON(iris.StatusOK, QueryResponse{RetCode:4000,Message:"getUsageByNow response.StatusCode != 200", Data:""})
		return
	}
	c.JSON(iris.StatusOK, QueryResponse{RetCode:0,Message:"ok",Data:string(body)})
	return

}

func createBucket(c *iris.Context, query QueryRequest) {
	ak := c.RequestHeader("X-Le-Key")
	sk := c.RequestHeader("X-Le-Secret")
	if ak == "" || sk == "" {
		c.JSON(iris.StatusOK, QueryResponse{RetCode:4101,Message:"invalid accessKey and secretKay",Data:""})
		return
	}

        tokenRecord := c.Get("token").(TokenRecord)
        if helper.Enforcer.Enforce(tokenRecord.UserName, API_LOS_GetS3Domain, ACT_ACCESS) != true {
                c.JSON(iris.StatusOK, QueryResponse{RetCode: 4030, Message: "You do not have permission to perform", Data: query})
                return
        }

        if query.Endpoint == "" {
                c.JSON(iris.StatusOK, QueryResponse{RetCode: 4000, Message: "please provide an endpoint to verify", Data: query})
                return
        }

	client, err := minio.NewV2(query.Endpoint, ak, sk, false)
	if err != nil {
		helper.Logger.Fatalln(0, err)
	}
	err = client.MakeBucket(query.Bucket, "")
	if err != nil {
		exists, err := minioClient.BucketExists(query.Bucket)
		if err == nil && exists {
			c.JSON(iris.StatusOK, QueryResponse{RetCode:4000,Message:"bucket existd",Data:""})
		} else {
			helper.Logger.Println(5, "createBucket error", err.Error(), ak, sk)
			c.JSON(iris.StatusOK, QueryResponse{RetCode:4101,Message:err.Error(),Data:""})
		}
		return
	}
	c.JSON(iris.StatusOK, QueryResponse{RetCode:0,Message:"bucket created",Data:""})
	return
}

func putCors(c *iris.Context, query QueryRequest) {
	ak := c.RequestHeader("X-Le-Key")
	sk := c.RequestHeader("X-Le-Secret")
	if ak == "" || sk == "" {
		c.JSON(iris.StatusOK, QueryResponse{RetCode:4101,Message:"invalid accessKey and secretKay",Data:""})
		return
	}

        tokenRecord := c.Get("token").(TokenRecord)
        if helper.Enforcer.Enforce(tokenRecord.UserName, API_LOS_GetS3Domain, ACT_ACCESS) != true {
                c.JSON(iris.StatusOK, QueryResponse{RetCode: 4030, Message: "You do not have permission to perform", Data: query})
                return
        }

        if query.Endpoint == "" {
                c.JSON(iris.StatusOK, QueryResponse{RetCode: 4000, Message: "please provide an endpoint to verify", Data: query})
                return
        }

	cors := `
	<CORSConfiguration>
		<CORSRule>
			<AllowedMethod>PUT</AllowedMethod>
			<AllowedMethod>GET</AllowedMethod>
			<AllowedMethod>POST</AllowedMethod>
			<AllowedMethod>DELETE</AllowedMethod>
			<AllowedOrigin>*</AllowedOrigin>
			<AllowedHeader>*</AllowedHeader>
			<ExposeHeader>x-amz-acl</ExposeHeader>
			<ExposeHeader>ETag</ExposeHeader>
		</CORSRule>
	</CORSConfiguration>`
	md5sum := md5.Sum([]byte(cors))
	str := base64.StdEncoding.EncodeToString(md5sum[:])
	client := &http.Client{}
	url := "http://" + query.Endpoint + "/" + query.Bucket + "?cors"
//	slog.Println("put cors url:", url)
	request, err := http.NewRequest("PUT", url, strings.NewReader(cors))
	if err != nil {
		c.JSON(iris.StatusOK, QueryResponse{RetCode:4000,Message:"make new request failed",Data:""})
		return
	}
	request.Header.Set("Content-type", "text/xml")
	request.Header.Set("Content-MD5", str)
	request_p := helper.SignV2(*request, ak, sk)
	response, err := client.Do(request_p)
	if err != nil {
		c.Text(500, err.Error())
		return
	}
//	body, _ := ioutil.ReadAll(response.Body)
//	slog.Println("put cors success", string(body))
	defer response.Body.Close()
	switch status := response.StatusCode; status {
	case 200:
		c.JSON(iris.StatusOK, QueryResponse{RetCode:0,Message:"put cors success",Data:""})
	case 403:
		c.JSON(iris.StatusOK, QueryResponse{RetCode:4101,Message:"Access Denied",Data:""})
	default:
		c.JSON(iris.StatusOK, QueryResponse{RetCode:4000,Message:"Unknown Error",Data:""})
	}
	return
}


func s3CreateKey(uid, accessKey, secretKey string) error {
	Client := &http.Client{Timeout: time.Second * 5}
	url := "http://" + helper.CONFIG.S3Domain + "/admin/user?format=json&uid=" + uid +"&display-name=" + uid +
		"&access-key=" + accessKey + "&secret-key=" + secretKey
	method := "PUT"

	helper.Logger.Println(5, "new request to s3:" + url)
	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		helper.Logger.Println(5, err.Error())
		return err
	}
	request_p := helper.SignV2(*request, helper.CONFIG.AccessKey, helper.CONFIG.SecretKey)
	response, err := Client.Do(request_p)
	if err != nil {
		helper.Logger.Println(5, err.Error())
		return err
	}
	defer response.Body.Close()
//	stdout := os.Stdout
//	_, err = io.Copy(stdout, response.Body)
	helper.Logger.Println(5, response.StatusCode)
	if response.StatusCode != 200 {
		return errors.New("create key return code not 200")
	}
	return nil
}

func s3DeleteKey(uid, accessKey, secretKey string) error {
	Client := &http.Client{Timeout: time.Second * 5}
	url := "http://" + helper.CONFIG.S3Domain + "/admin/user?key&format=json&uid=" + uid +
		"&access-key=" + accessKey + "&secret-key=" + secretKey
	method := "DELETE"

	helper.Logger.Println(5, "new request to s3:" + url)
	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		helper.Logger.Println(5, err.Error())
		return err
	}
	request_p := helper.SignV2(*request, helper.CONFIG.AccessKey, helper.CONFIG.SecretKey)
	response, err := Client.Do(request_p)
	if err != nil {
		helper.Logger.Println(5, err.Error())
		return err
	}
	defer response.Body.Close()
	stdout := os.Stdout
	_, err = io.Copy(stdout, response.Body)
	if err != nil {
		helper.Logger.Println(5, err.Error())
		return err
	}
	helper.Logger.Println(5, response.StatusCode)
	if response.StatusCode != 200 {
		if response.StatusCode == 403 {
			return fmt.Errorf("InvalidAccessKeyId")
		} else {
			return fmt.Errorf("return code not 204")
		}
	}
	return nil

}
