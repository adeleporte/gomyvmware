# My-Vmware Go Client

This client features a way to list and download files from my-vmware website.


## Contents

- [My-VMware Go client](#my-vmware-go-client)
  - [Contents](#contents)
  - [Installation](#installation)
  - [Quick start](#quick-start)


## Installation

To install My-VMware go client, you need to install Go and set your Go workspace first.

1. The first need [Go](https://golang.org/) installed (**version 1.12+ is required**), then you can use the below Go command to install Gin.

```sh
$ go get -u github.com/adeleporte/go-my-vmware
```

2. Import it in your code:

```go
import "github.com/adeleporte/go-my-vmware"
```


## Quick start

```sh
# assume the following codes in example.go file
$ cat example.go
```

```go
package main

func main() {

	client, err := NewClient("test@vmware.com", "changeme")
	if err != nil {
		log.Fatal(err)
	}

	results, err := GetProductsAtoZ(client)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(results.ProductCategoryList)
	for _, product := range results.ProductCategoryList[0].ProductList {
		log.Println(product.Name)
	}

	params := DLGListParams{
		Category: "networking_security",
		Product:  "vmware_nsx_t_data_center",
		Version:  "3_x",
		DlgType:  "PRODUCT_BINARY",
	}

	_, err = GetRelatedDLGList(client, params)
	if err != nil {
		log.Fatal(err)
	}
	//log.Printf("%+v", results2)

	dlg_params := GetDLGParams{
		DownloadGroup: "NSX-T-30110",
		ProductID:     982,
	}
	headers, err := GetDLGHeader(client, dlg_params)
	if err != nil {
		log.Fatal(err)
	}
	//log.Printf("headers: %+v\n", headers)

	details, err := GetDLGDetails(client, dlg_params)
	if err != nil {
		log.Fatal(err)
	}
	//log.Printf("detais: %+v\n", details)

	file, err := details.Filter("nsx-unified-appliance-.*.ova")
	if err != nil {
		log.Fatal(err)
	}
	//log.Printf("file: %+v\n", file)

	_, err = GetEulaAccept(client, headers.Dlg.Code)
	if err != nil {
		log.Fatal(err)
	}

	body := GetDownloadBody{
		Locale:        "en_US",
		DownloadGroup: headers.Dlg.Code,
		ProductId:     headers.Product.ID,
		Md5checksum:   file.Md5checksum,
		TagID:         headers.Dlg.TagID,
		UUID:          file.UUID,
		DlgType:       headers.Dlg.Type,
		ProductFamily: headers.Product.Name,
		ReleaseDate:   file.ReleaseDate,
		DlgVersion:    file.Version,
		IsBetaFlow:    false,
	}
	download_link, err := GetDownload(client, body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("**************************************")
	log.Printf("File Name: %s\n", download_link.FileName)
	//log.Printf("Download link: %s\n", download_link.DownloadURL)
	log.Println("Downloading...")
	Download(client, download_link)
	log.Println("Downloaded...")

}
```

```
# run example.go
$ go run example.go
```
