package view

type MinioStorageCreds struct {
	BucketName      string // minio bucket name
	IsActive        bool   // a flag indicates full-fledged interface
	Endpoint        string // minio endpoint address
	Crt             string // minio certificate
	AccessKeyId     string // minio access key ID
	SecretAccessKey string // secret minio access key
	ProductionMode  bool   // production mode flag
	WorkDir         string // local working directory to store intermediate files
}
