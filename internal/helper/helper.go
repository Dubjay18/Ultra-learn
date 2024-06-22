package helper

import (
	"fmt"
	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"os"
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
