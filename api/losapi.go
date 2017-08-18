package api

import (
	"log"
	"crypto/md5"
	"time"
	"strings"
	"io/ioutil"
	"net/http"
	"gopkg.in/iris.v4"
	"encoding/base64"
	. "github.com/journeymidnight/yig-iam/api/datatype"
	"github.com/journeymidnight/yig-iam/helper"
	"github.com/minio/minio-go"	
)


var slog  *log.Logger
var minioClient *minio.Client

type GetS3Domain struct {
	S3Domain string `json:"s3Domain"`
}

func LosApiHandler(c *iris.Context) {
	query := c.Get("queryRequest").(QueryRequest)
	switch query.Action {
		case ACTION_LOS_GETS3DOMAIN:
			getS3Domain(c)
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

func getS3Domain(c *iris.Context) {
	var data GetS3Domain
	data.S3Domain = helper.CONFIG.S3Domain
	c.JSON(iris.StatusOK, QueryResponse{RetCode:0,Message:"",Data:data})
}

func listBuckets(c *iris.Context, query QueryRequest) {
	ak := c.RequestHeader("X-Le-Key")
	sk := c.RequestHeader("X-Le-Secret")
	if ak == "" || sk == "" {
		c.JSON(iris.StatusOK, QueryResponse{RetCode:4101,Message:"invalid accessKey and secretKay",Data:""})
		return
	}
	//client, err := minio.NewV2(helper.CONFIG.S3Domain, ak, sk, helper.CONFIG.UseSSL)
	client, err := minio.NewV2(helper.CONFIG.S3Domain, ak, sk, false)
//	slog.Println("listbuckets:", cfg.S3Domain,ak,sk,cfg.UseSSL)
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
	client, err := minio.NewV2(helper.CONFIG.S3Domain, ak, sk, false)
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
	Client := &http.Client{Timeout: time.Second * 5}
	url := "http://" + helper.CONFIG.S3Domain + "/admin/bucket?format=json&bucket=" + query.Bucket +"&stats=False" +"uid=" + query.ProjectId
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
//	slog.Println(response.StatusCode)
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
	client, err := minio.NewV2(helper.CONFIG.S3Domain, ak, sk, false)
	if err != nil {
		slog.Fatalln(err)
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
	url := "http://" + helper.CONFIG.S3Domain + "/" + query.Bucket + "?cors"
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