package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"github.com/ton-developer-program/internal/response"
)

// GenerateUUID generates a random UUID
func generateUUID() string {
	return uuid.New().String()
}

func (app *application) uploadImageHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form data with a specified max file size
	r.ParseMultipartForm(10 << 20) // 10 MB

	// Access the file from the request form data
	file, handler, err := r.FormFile("file")
	if err != nil {
		app.serverError(w, r, err) 
		 app.logger.Error(err,nil)
		return
	}
	defer file.Close()

	// Generate a unique filename
	filename := handler.Filename
	ext := filepath.Ext(filename)
	filename = fmt.Sprintf("%s%s", generateUUID(), ext)

	// Create an AWS session
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(app.config.AWS.AWSRegion),
		Credentials: credentials.NewStaticCredentials(app.config.AWS.AWSAccessKeyID, app.config.AWS.AWSSecretAccessKey, ""),
	})
	if err != nil {
		app.serverError(w, r, err) 
		 app.logger.Error(err,nil)
		return
	}

	// Upload the file to S3
	_, err = s3.New(sess).PutObject(&s3.PutObjectInput{
		Bucket: aws.String(app.config.AWS.AWSBucket),
		Key:    aws.String(filename),
		Body:   file,
		ACL:    aws.String("public-read"),
	})
	if err != nil {
		app.serverError(w, r, err) 
		 app.logger.Error(err,nil)
		return
	}

	// Send a success response
	response.JSON(w, http.StatusOK, map[string]string{
		"url": fmt.Sprintf("https://%s.s3-%s.amazonaws.com/%s", app.config.AWS.AWSBucket, app.config.AWS.AWSRegion, filename),
	})
}

func (app *application) uploadCSVHandler(w http.ResponseWriter, r *http.Request) {
	// csv file must be not more than 80 rows

	// Parse the multipart form data with a specified max file size
	r.ParseMultipartForm(10 << 20) // 10 MB

	// Access the file from the request form data
	file, handler, err := r.FormFile("file")
	if err != nil {
		app.serverError(w, r, err) 
		 app.logger.Error(err,nil)
		return
	}
	defer file.Close()

	read := csv.NewReader(file)
	records, err := read.ReadAll()
	if err != nil {
		app.serverError(w, r, err) 
		 app.logger.Error(err,nil)
		return
	}

	// Reset the file pointer to the beginning of the file
	file.Seek(0, 0)

	if len(records) > 80 {
		response.JSON(w, http.StatusBadRequest, map[string]string{
			"message": "csv file must be not more than 80 rows",
		})
		return
	}
	// Generate a unique filename
	filename := handler.Filename
	ext := filepath.Ext(filename)
	filename = fmt.Sprintf("%s%s", generateUUID(), ext)

	// Create an AWS session
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(app.config.AWS.AWSRegion),
		Credentials: credentials.NewStaticCredentials(app.config.AWS.AWSAccessKeyID, app.config.AWS.AWSSecretAccessKey, ""),
	})
	if err != nil {
		app.serverError(w, r, err) 
		 app.logger.Error(err,nil)
		return
	}

	// Upload the file to S3
	_, err = s3.New(sess).PutObject(&s3.PutObjectInput{
		Bucket: aws.String(app.config.AWS.AWSBucket),
		Key:    aws.String(filename),
		Body:   file,
		ACL:    aws.String("public-read"),
	})
	if err != nil {
		app.serverError(w, r, err) 
		 app.logger.Error(err,nil)
		return
	}

	
	// Send a success response
	response.JSON(w, http.StatusOK, map[string]string{
		"url": fmt.Sprintf("https://%s.s3-%s.amazonaws.com/%s", app.config.AWS.AWSBucket, app.config.AWS.AWSRegion, filename),
	})
}
