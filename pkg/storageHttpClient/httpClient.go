package storageHttpClient

import "github.com/von1994/CSIFramework/driver/util"

type StorageHttpClient struct {
	Endponit string
	ReqMethod  string
	ReqTimeout int
	Scheme     string
	Domain string
	SecretID string
	SecretKey string
}

func NewStorageHttpClient(storageUrl, endpoint string) *StorageHttpClient{
	secretID,secretKey := util.GetSercet()
	return &StorageHttpClient{
		Endponit: endpoint,
		Domain: storageUrl,
		SecretID: secretID,
		SecretKey: secretKey,
	}
}
