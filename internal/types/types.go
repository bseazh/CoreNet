package types

type InitReq struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
	Mime string `json:"mime"`
}

type InitResp struct {
	UploadId  string `json:"uploadId"`
	ChunkSize int32  `json:"chunkSize"`
}

type CompleteReq struct {
	UploadId string `json:"uploadId"`
}

type CompleteResp struct {
	FileId string `json:"fileId"`
}

type FileResp struct {
	FileId  string `json:"fileId"`
	Name    string `json:"name"`
	Size    int64  `json:"size"`
	Mime    string `json:"mime"`
	Sha1    string `json:"sha1"`
	Version int32  `json:"version"`
}

type PreviewResp struct {
	Url string `json:"url"`
}

type SearchResp struct {
	Items []SearchItem `json:"items"`
}

type SearchItem struct {
	FileId  string `json:"fileId"`
	Name    string `json:"name"`
	Snippet string `json:"snippet,optional"`
}

type CreateOCRReq struct {
	FileId string `json:"fileId"`
}

type CreateOCRResp struct {
	JobId string `json:"jobId"`
}

type JobResp struct {
	JobId     string `json:"jobId"`
	Status    string `json:"status"`
	ResultURI string `json:"resultURI,optional"`
}
