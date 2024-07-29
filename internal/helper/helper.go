package helper

import (
	"Ultra-learn/internal/dto"
	"Ultra-learn/internal/errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/iancoleman/strcase"
)

func GenerateUserId() string {
	return uuid.New().String()
}

func UploadImage(file any, user_id string, c *gin.Context) (string, error) {
	cldApiKey := os.Getenv("CLOUDINARY_API_KEY")
	cldApiSecret := os.Getenv("CLOUDINARY_API_SECRET")
	//fileTags := ".jpg,.png,.jpeg"
	cld, cerr := cloudinary.NewFromURL(fmt.Sprintf("cloudinary://%v:%v@dubinx", cldApiKey, cldApiSecret))

	if cerr != nil {
		return "", cerr

	}
	result, err := cld.Upload.Upload(c, file, uploader.UploadParams{
		PublicID: user_id,
		// Split the tags by comma
		//Tags: strings.Split(",", fileTags),
	})

	if err != nil {
		return "", err
	}
	return result.SecureURL, nil
}

func ParseRequestBody(c *gin.Context, req interface{}) any {
	log.Println("Parsing request body")
	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		if validationErrs, ok := bindErr.(validator.ValidationErrors); ok {
			// Handle validation errors
			var res []errors.FieldError
			for _, e := range validationErrs {
				// Extract the field name and the error message
				fieldName := strings.Split(e.Namespace(), ".")[1]
				fieldName = strcase.ToLowerCamel(fieldName)
				errorMessage := e.ActualTag()
				// Translate each error one at a time
				res = append(res, errors.FieldError{Field: fieldName, Message: errorMessage})
			}
			c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": res})

		} else {
			// Handle other errors (like invalid JSON)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		}
		return bindErr
	}
	return nil
}

func ValidateRequest(c *gin.Context, req interface{}) any {
	log.Println("Validating request body")
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errors.ApiError{
			Message:    "Invalid request",
			Error:      err.Error(),
			StatusCode: http.StatusBadRequest,
		})
		return err
	}
	return nil
}

func BuildSuccessResponse(c *gin.Context, status int, message string, data any) {
	rd := dto.ApiSuccessResponse{
		Message:    message,
		StatusCode: status,
		Data:       data,
	}
	c.JSON(status, rd)
}

func BuildErrorResponse(c *gin.Context, status int, message string, err error) {
	rd := errors.ApiError{
		Message:    message,
		StatusCode: status,
		Error:      err.Error(),
	}
	c.JSON(status, rd)
}
