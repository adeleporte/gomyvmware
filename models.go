package gomyvmware

import "net/http"

type MyClient struct {
	Client      *http.Client
	XSRF_Cookie http.Cookie
}

type ProductsAtoZ struct {
	ProductCategoryList []struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		ProductList []struct {
			Name    string `json:"name"`
			Actions []struct {
				LinkName string `json:"linkname"`
				Target   string `json:"target"`
				OrderId  int    `json:"orderId"`
			} `json:"actions"`
		} `json:"productList"`
	} `json:"productCategoryList"`
}

type DLGListParams struct {
	Category string
	Product  string
	Version  string
	DlgType  string
}

type DLGListResults struct {
	DlgEditionsLists []struct {
		OrderId int    `json:"orderId"`
		Name    string `json:"name"`
		DlgList []struct {
			Name             string `json:"name"`
			Code             string `json:"code"`
			ReleaseDate      string `json:"releaseDate"`
			ProductID        string `json:"productId"`
			ReleasePackageId string `json:"releasePackageId"`
			OrderID          int    `json:"orderId"`
		} `json:"dlgList"`
	} `json:"dlgEditionsLists"`
}

type GetDLGParams struct {
	DownloadGroup string
	ProductID     int
}

type GetDLGHeaderResults struct {
	Product struct {
		ID               string `json:"id"`
		Name             string `json:"name"`
		ReleasePackageId string `json:"releasePackageId"`
	} `json:"product"`
	Dlg struct {
		Name        string `json:"name"`
		Code        string `json:"code"`
		ReleaseDate string `json:"releaseDate"`
		TagID       int    `json:"tagId"`
		Type        string `json:"type"`
	} `json:"dlg"`
}

type DownloadFile struct {
	UUID           string `json:"uuid"`
	FileName       string `json:"fileName"`
	Sha1checksum   string `json:"sha1checksum"`
	Sha256checksum string `json:"sha256checksum"`
	Md5checksum    string `json:"md5checksum"`
	Build          string `json:"build"`
	ReleaseDate    string `json:"releaseDate"`
	FileType       string `json:"fileType"`
	Description    string `json:"description"`
	FileSize       string `json:"fileSize"`
	Title          string `json:"title"`
	Version        string `json:"version"`
	Status         string `json:"status"`
}

type GetDLGDetailsResults struct {
	DownloadFiles []DownloadFile `json:"downloadFiles"`
}

type GetDownloadBody struct {
	Locale        string `json:"locale"`
	DownloadGroup string `json:"downloadGroup"`
	ProductId     string `json:"productId"`
	Md5checksum   string `json:"md5checksum"`
	TagID         int    `json:"tagId"`
	UUID          string `json:"uUId"`
	DlgType       string `json:"dlgType"`
	ProductFamily string `json:"productFamily"`
	ReleaseDate   string `json:"releaseDate"`
	DlgVersion    string `json:"dlgVersion"`
	IsBetaFlow    bool   `json:"isBetaFlow"`
}

type GetDownloadResults struct {
	DownloadURL string `json:"downloadURL"`
	FileName    string `json:"fileName"`
}

type WriteCounter struct {
	Total uint64
}
