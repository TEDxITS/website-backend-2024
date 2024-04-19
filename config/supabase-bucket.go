package config

import (
	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type (
	SupabaseBucket struct {
		Session   http.Client
		BucketID  string
		BucketURL string
		token     string
	}
)

// https://github.com/supabase-community/storage-go/tree/dev
// https://github.com/adityarizkyramadhan/supabase-storage-uploader/tree/master
func SetUpSupabaseBucket() *SupabaseBucket {
	key := os.Getenv("SUPABASE_API_KEY")
	bucketID := os.Getenv("SUPABASE_BUCKET_ID")
	url := os.Getenv("SUPABASE_PROJECT_URL") + "/storage/v1/object/" + bucketID + "/"

	return &SupabaseBucket{
		Session:   http.Client{},
		BucketID:  bucketID,
		BucketURL: url,
		token:     key,
	}
}

func (b *SupabaseBucket) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+b.token)

	return b.Session.Do(req)
}
