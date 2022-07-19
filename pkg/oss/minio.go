/*
 * @Author: lqc
 * @Date: 2021-11-23 11:26:24
 * @Description: minio客户端初始化
 */

package oss

import (
	"context"
	"io"
	"time"

	"github.com/ashin-l/go-demo/pkg/logger"
	"github.com/ashin-l/go-demo/pkg/option"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pkg/errors"
	"github.com/rs/xid"
)

var (
	c          *minio.Client
	endpoint   string
	bucketName string
	prefix     string
	ctx        = context.Background()
)

func Init(opt *option.Options) {
	bucketName = opt.Minio.Bucketname
	prefix = opt.Minio.Prefix + "/" + opt.Minio.Bucketname
	endpoint = opt.Minio.Endpoint
	accessKeyID := opt.Minio.Accesskeyid
	secretAccessKey := opt.Minio.Secretaccesskey
	useSSL := opt.Minio.Usessl

	// Initialize minio client object.
	var err error
	c, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		logger.Logger().Fatal(err)
	}

	exists, errBucketExists := c.BucketExists(ctx, bucketName)
	if errBucketExists == nil && exists {
		logger.Logger().Infof("We already own %s\n", bucketName)
	} else {
		logger.Logger().Fatal(err)
	}
}

func genXidName(dir, deviceId string) string {
	id := xid.New().String()
	return dir + "/" + deviceId + "/" + time.Now().Format("20060102") + "/" + id + ".jpg"
}

func genTimeName(dir, deviceId string) string {
	now := time.Now()
	id := now.Format("20060102150405")
	return dir + "/" + deviceId + "/" + time.Now().Format("20060102") + "/" + id + ".jpg"
}

func PutObject(reader io.Reader, objectSize int64, dir string, deviceId string) (string, string, error) {
	objectName := genTimeName(dir, deviceId)
	_, err := c.PutObject(ctx, bucketName, objectName, reader, objectSize, minio.PutObjectOptions{ContentType: "image/jpg"})
	if err != nil {
		return "", "", errors.Wrap(err, "上传文件报错")
	}

	return prefix + "/" + objectName, objectName, nil
}

func RemoveObject(objectName string) {
	if err := c.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{}); err != nil {
		logger.Logger().Warn(err)
	}
}
