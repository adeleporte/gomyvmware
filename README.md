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
$ go get -u github.com/adeleporte/gomyvmware
```

2. Import it in your code:

```go
import "github.com/adeleporte/gomyvmware"
```


## Quick start

```sh
# assume the following codes in example.go file
$ cat example.go
```

```go
package main

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	vmw "github.com/adeleporte/gomyvmware"
)

func main() {

	w := tabwriter.NewWriter(os.Stdout, 20, 8, 1, ' ', tabwriter.Debug)

	fmt.Println("Authenticating...")
	client, err := vmw.NewClient("test@vmware.com", "changeme")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Authenticated...")

	fmt.Println("...")
	fmt.Println("Get AtoZ Products...")
	results, err := vmw.GetProductsAtoZ(client)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("...")
	fmt.Fprintf(w, "%s\t%s\t\n", "list.Name", "action.LinkName")
	for _, category := range results.ProductCategoryList {
		for _, list := range category.ProductList {
			for _, action := range list.Actions {
				fmt.Fprintf(w, "%s\t%s\t\n", list.Name, action.LinkName)
			}
		}

	}
	w.Flush()

	params := vmw.DLGListParams{
		Category: "networking_security",
		Product:  "vmware_nsx_t_data_center",
		Version:  "3_x",
		DlgType:  "PRODUCT_BINARY",
	}

	fmt.Println("...")
	fmt.Println("Get NSX-T 3.0 info...")
	nsx_list, err := vmw.GetRelatedDLGList(client, params)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", "Category Name", "Product Name", "Group", "ProductID")
	for _, nsx := range nsx_list.DlgEditionsLists {
		for _, nsx_dlg := range nsx.DlgList {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", nsx.Name, nsx_dlg.Name, nsx_dlg.Code, nsx_dlg.ProductID)
		}
	}
	w.Flush()
	//log.Printf("%+v", results2)

	fmt.Println("...")
	fmt.Println("Get details on a specific NSX-T release...")
	dlg_params := vmw.GetDLGParams{
		DownloadGroup: "NSX-T-30110",
		ProductID:     982,
	}
	headers, err := vmw.GetDLGHeader(client, dlg_params)
	if err != nil {
		log.Fatal(err)
	}
	//log.Printf("headers: %+v\n", headers)

	details, err := vmw.GetDLGDetails(client, dlg_params)
	if err != nil {
		log.Fatal(err)
	}
	//log.Printf("detais: %+v\n", details)

	file, err := details.Filter("nsx-unified-appliance-.*.ova")
	if err != nil {
		log.Fatal(err)
	}
	//log.Printf("file: %+v\n", file)

	_, err = vmw.GetEulaAccept(client, headers.Dlg.Code)
	if err != nil {
		log.Fatal(err)
	}

	body := vmw.GetDownloadBody{
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
	download_link, err := vmw.GetDownload(client, body)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("File Name: %s\n", download_link.FileName)
	log.Printf("Download URL: %s\n", download_link.DownloadURL)

	log.Println("**************************************")
	log.Println("Downloading...")
	vmw.Download(client, download_link)
	log.Println("Downloaded...")

}
```

```
# run example.go
deleporte@adeleporte-a02 test-govmware % go run main.go
Authenticating...
Authenticated...
...
Get AtoZ Products...
...
list.Name                                                             |action.LinkName     |
Advanced Monitoring for VMware Horizon by ControlUp                   |Download Product    |
Advanced Monitoring for VMware Horizon by ControlUp                   |Drivers & Tools     |
Downloads                                                             |Download Product    |
Downloads                                                             |Drivers & Tools     |
VMWare NSX Intelligence                                               |Download Product    |
VMWare NSX Intelligence                                               |Drivers & Tools     |
VMmark                                                                |Download Product    |
VMmark                                                                |Drivers & Tools     |
VMware App Volumes                                                    |Download Product    |
VMware App Volumes                                                    |Drivers & Tools     |
VMware App Volumes                                                    |Download Trial      |
VMware AppDefense Plugin                                              |Download Product    |
VMware AppDefense Plugin                                              |Drivers & Tools     |
VMware Blockchain                                                     |Download Product    |
VMware Blockchain                                                     |Drivers & Tools     |
VMware Carbon Black Cloud Workload                                    |Download Product    |
VMware Carbon Black Cloud Workload                                    |Drivers & Tools     |
VMware Cloud Director                                                 |Download Product    |
VMware Cloud Director                                                 |Drivers & Tools     |
VMware Cloud Director App Launchpad                                   |Download Product    |
VMware Cloud Director App Launchpad                                   |Drivers & Tools     |
VMware Cloud Director Availability                                    |Download Product    |
VMware Cloud Director Availability                                    |Drivers & Tools     |
VMware Cloud Director Object Storage Extension                        |Download Product    |
VMware Cloud Director Object Storage Extension                        |Drivers & Tools     |
VMware Cloud Foundation                                               |Download Product    |
VMware Cloud Foundation                                               |Drivers & Tools     |
VMware Cloud Services                                                 |Download Product    |
VMware Cloud Services                                                 |Drivers & Tools     |
VMware Cloud on AWS                                                   |Download Product    |
VMware Cloud on AWS                                                   |Drivers & Tools     |
VMware Dynamic Environment Manager                                    |Download Product    |
VMware Dynamic Environment Manager                                    |Drivers & Tools     |
VMware Dynamic Environment Manager                                    |Download Trial      |
VMware Enterprise PKS                                                 |Download Product    |
VMware Enterprise PKS                                                 |Drivers & Tools     |
VMware Essential PKS                                                  |Download Product    |
VMware Essential PKS                                                  |Drivers & Tools     |
VMware Fusion                                                         |Download Product    |
VMware Fusion                                                         |Drivers & Tools     |
VMware Fusion                                                         |Download Trial      |
VMware HCX                                                            |Download Product    |
VMware HCX                                                            |Drivers & Tools     |
VMware Horizon                                                        |Download Product    |
VMware Horizon                                                        |Drivers & Tools     |
VMware Horizon                                                        |Download Trial      |
VMware Horizon (with View)                                            |Download Product    |
VMware Horizon (with View)                                            |Drivers & Tools     |
VMware Horizon Air                                                    |Download Product    |
VMware Horizon Air                                                    |Drivers & Tools     |
VMware Horizon Apps                                                   |Download Product    |
VMware Horizon Apps                                                   |Drivers & Tools     |
VMware Horizon Clients                                                |Download Product    |
VMware Horizon Clients                                                |Drivers & Tools     |
VMware Horizon DaaS                                                   |Download Product    |
VMware Horizon DaaS                                                   |Drivers & Tools     |
VMware Horizon Service                                                |Download Product    |
VMware Horizon Service                                                |Drivers & Tools     |
VMware Integrated OpenStack                                           |Download Product    |
VMware Integrated OpenStack                                           |Drivers & Tools     |
VMware Integrated OpenStack                                           |Download Trial      |
VMware Internet of things                                             |Download Product    |
VMware Internet of things                                             |Drivers & Tools     |
VMware Internet of things                                             |Download Trial      |
VMware Mirage                                                         |Download Product    |
VMware Mirage                                                         |Drivers & Tools     |
VMware Mirage                                                         |Download Trial      |
VMware NSX Advanced Load Balancer                                     |Download Product    |
VMware NSX Advanced Load Balancer                                     |Drivers & Tools     |
VMware NSX Cloud                                                      |Download Product    |
VMware NSX Cloud                                                      |Drivers & Tools     |
VMware NSX Data Center for vSphere                                    |Download Product    |
VMware NSX Data Center for vSphere                                    |Drivers & Tools     |
VMware NSX Security                                                   |Download Product    |
VMware NSX Security                                                   |Drivers & Tools     |
VMware NSX-T Data Center                                              |Download Product    |
VMware NSX-T Data Center                                              |Drivers & Tools     |
VMware NSX-T Data Center                                              |Download Trial      |
VMware Photon Platform                                                |Download Product    |
VMware Photon Platform                                                |Drivers & Tools     |
VMware SD-WAN                                                         |Download Product    |
VMware SD-WAN                                                         |Drivers & Tools     |
VMware Site Recovery Manager                                          |Download Product    |
VMware Site Recovery Manager                                          |Drivers & Tools     |
VMware Site Recovery Manager                                          |Download Trial      |
VMware Skyline Collector                                              |Download Product    |
VMware Skyline Collector                                              |Drivers & Tools     |
VMware Smart Assurance                                                |Download Product    |
VMware Smart Assurance                                                |Drivers & Tools     |
VMware Smart Experience                                               |Download Product    |
VMware Smart Experience                                               |Drivers & Tools     |
VMware Tanzu Advanced Edition                                         |Download Product    |
VMware Tanzu Kubernetes Grid                                          |Download Product    |
VMware Tanzu Kubernetes Grid                                          |Drivers & Tools     |
VMware Tanzu Kubernetes Grid Integrated Edition                       |Download Product    |
VMware Tanzu Kubernetes Grid Integrated Edition                       |Drivers & Tools     |
VMware Tanzu Toolkit for Kubernetes                                   |Download Product    |
VMware Tanzu Toolkit for Kubernetes                                   |Drivers & Tools     |
VMware Telco Cloud Automation                                         |Download Product    |
VMware Telco Cloud Automation                                         |Drivers & Tools     |
VMware Telco Cloud Automation                                         |Download Trial      |
VMware Telco Cloud Infrastructure - Cloud Director Edition            |Download Product    |
VMware Telco Cloud Infrastructure - Cloud Director Edition            |Drivers & Tools     |
VMware Telco Cloud Infrastructure - OpenStack Edition                 |Download Product    |
VMware Telco Cloud Infrastructure - OpenStack Edition                 |Drivers & Tools     |
VMware Telco Cloud Operations                                         |Download Product    |
VMware Telco Cloud Operations                                         |Drivers & Tools     |
VMware Telco Cloud Platform-5G Edition                                |Download Product    |
VMware Telco Cloud Platform-5G Edition                                |Drivers & Tools     |
VMware ThinApp                                                        |Download Product    |
VMware ThinApp                                                        |Drivers & Tools     |
VMware ThinApp                                                        |Download Trial      |
VMware Tools                                                          |Download Product    |
VMware Tools                                                          |Drivers & Tools     |
VMware TrustPoint                                                     |Download Product    |
VMware TrustPoint                                                     |Drivers & Tools     |
VMware Unified Access Gateway                                         |Download Product    |
VMware Unified Access Gateway                                         |Drivers & Tools     |
VMware Validated Design for Software-Defined Data Center              |Download Product    |
VMware Validated Design for Software-Defined Data Center              |Drivers & Tools     |
VMware View Planner                                                   |Download Product    |
VMware View Planner                                                   |Drivers & Tools     |
VMware Workspace                                                      |Download Product    |
VMware Workspace                                                      |Drivers & Tools     |
VMware Workspace ONE                                                  |Download Product    |
VMware Workspace ONE                                                  |Drivers & Tools     |
VMware Workspace ONE Access (VIDM)                                    |Download Product    |
VMware Workspace ONE Access (VIDM)                                    |Drivers & Tools     |
VMware Workspace ONE Access (VIDM)                                    |Download Trial      |
VMware Workstation Player                                             |Download Product    |
VMware Workstation Pro                                                |Download Product    |
VMware Workstation Pro                                                |Drivers & Tools     |
VMware vCenter Converter Standalone                                   |Download Product    |
VMware vCenter Converter Standalone                                   |Drivers & Tools     |
VMware vCloud Air                                                     |Download Product    |
VMware vCloud Air                                                     |Drivers & Tools     |
VMware vCloud Availability for Cloud-to-Cloud DR                      |Download Product    |
VMware vCloud Availability for Cloud-to-Cloud DR                      |Drivers & Tools     |
VMware vCloud NFV                                                     |Download Product    |
VMware vCloud NFV                                                     |Drivers & Tools     |
VMware vCloud NFV OpenStack Edition                                   |Download Product    |
VMware vCloud NFV OpenStack Edition                                   |Drivers & Tools     |
VMware vCloud Suite                                                   |Download Product    |
VMware vCloud Suite                                                   |Drivers & Tools     |
VMware vCloud Suite Platinum                                          |Download Product    |
VMware vCloud Suite Platinum                                          |Drivers & Tools     |
VMware vCloud Suite Platinum                                          |Buy Now             |
VMware vCloud Suite Subscription                                      |Download Product    |
VMware vCloud Suite Subscription                                      |Drivers & Tools     |
VMware vCloud Usage Meter                                             |Download Product    |
VMware vCloud Usage Meter                                             |Drivers & Tools     |
VMware vRealize Automation                                            |Download Product    |
VMware vRealize Automation                                            |Drivers & Tools     |
VMware vRealize Automation                                            |Download Trial      |
VMware vRealize Business for Cloud                                    |Download Product    |
VMware vRealize Business for Cloud                                    |Drivers & Tools     |
VMware vRealize Cloud Universal                                       |Download Product    |
VMware vRealize Cloud Universal                                       |Drivers & Tools     |
VMware vRealize Code Stream                                           |Download Product    |
VMware vRealize Code Stream                                           |Drivers & Tools     |
VMware vRealize Code Stream                                           |Download Trial      |
VMware vRealize Configuration Manager                                 |Download Product    |
VMware vRealize Configuration Manager                                 |Drivers & Tools     |
VMware vRealize Hyperic                                               |Download Product    |
VMware vRealize Hyperic                                               |Drivers & Tools     |
VMware vRealize Log Insight                                           |Download Product    |
VMware vRealize Log Insight                                           |Drivers & Tools     |
VMware vRealize Log Insight                                           |Download Trial      |
VMware vRealize Network Insight                                       |Download Product    |
VMware vRealize Network Insight                                       |Drivers & Tools     |
VMware vRealize Network Insight                                       |Download Trial      |
VMware vRealize Operations                                            |Download Product    |
VMware vRealize Operations                                            |Drivers & Tools     |
VMware vRealize Operations                                            |Download Trial      |
VMware vRealize Operations Insight                                    |Download Product    |
VMware vRealize Operations Insight                                    |Drivers & Tools     |
VMware vRealize Operations Management Pack for Care Systems Analytics |Download Product    |
VMware vRealize Operations Management Pack for Care Systems Analytics |Drivers & Tools     |
VMware vRealize Operations Management Pack for MEDITECH               |Download Product    |
VMware vRealize Operations Management Pack for MEDITECH               |Drivers & Tools     |
VMware vRealize Operations for Horizon and Published Applications     |Download Product    |
VMware vRealize Operations for Horizon and Published Applications     |Drivers & Tools     |
VMware vRealize Operations for Horizon and Published Applications     |Download Trial      |
VMware vRealize Operations for IBM Power Systems                      |Download Product    |
VMware vRealize Operations for IBM Power Systems                      |Drivers & Tools     |
VMware vRealize Suite                                                 |Download Product    |
VMware vRealize Suite                                                 |Drivers & Tools     |
VMware vRealize True Visibility Management Packs                      |Download Product    |
VMware vRealize True Visibility Management Packs                      |Drivers & Tools     |
VMware vRealize True Visibility Suite                                 |Download Product    |
VMware vRealize True Visibility Suite                                 |Drivers & Tools     |
VMware vSAN                                                           |Download Product    |
VMware vSAN                                                           |Drivers & Tools     |
VMware vSAN                                                           |Download Trial      |
VMware vSphere                                                        |Download Product    |
VMware vSphere                                                        |Drivers & Tools     |
VMware vSphere                                                        |Download Trial      |
VMware vSphere Bitfusion                                              |Download Product    |
VMware vSphere Bitfusion                                              |Drivers & Tools     |
VMware vSphere Data Protection Advanced                               |Download Product    |
VMware vSphere Data Protection Advanced                               |Drivers & Tools     |
VMware vSphere Hypervisor (ESXi)                                      |Download Product    |
VMware vSphere Hypervisor (ESXi)                                      |Drivers & Tools     |
VMware vSphere Integrated Containers                                  |Download Product    |
VMware vSphere Integrated Containers                                  |Drivers & Tools     |
VMware vSphere Storage Appliance                                      |Download Product    |
VMware vSphere Storage Appliance                                      |Drivers & Tools     |
...
Get NSX-T 3.0 info...
Category Name                          |Product Name                                                              |Group               |ProductID
VMware Network Management add-on       |VMware vRealize Network Insight 6.1.0                                     |VRNI-610            |982
VMware NSX Data Center Enterprise Plus |VMware HCX                                                                |HCX_353             |982
VMware NSX Data Center Enterprise Plus |VMware Identity Manager 3.3.3 (for vRA, vRops, vRLI, vRB, vRNI, NSX only) |VIDM_ONPREM_3330    |982
VMware NSX Data Center Enterprise Plus |VMware vRealize Log Insight 8.2.0 for NSX                                 |VRLI-820-NSX        |982
VMware NSX Data Center Enterprise Plus |VMware NSX-T Data Center 3.1.0                                            |NSX-T-310           |982
VMware NSX Data Center Enterprise Plus |VMware NSX Intelligence 1.2.0                                             |NSX-INTLG-120       |982
VMware NSX Data Center Enterprise Plus |VMware vRealize Network Insight 6.1.0                                     |VRNI-610            |982
VMware NSX Data Center Advanced        |VMware Identity Manager 3.3.3 (for vRA, vRops, vRLI, vRB, vRNI, NSX only) |VIDM_ONPREM_3330    |982
VMware NSX Data Center Advanced        |VMware vRealize Log Insight 8.2.0 for NSX                                 |VRLI-820-NSX        |982
VMware NSX Data Center Advanced        |VMware NSX-T Data Center 3.1.0                                            |NSX-T-310           |982
VMware NSX Data Center Professional    |VMware Identity Manager 3.3.3 (for vRA, vRops, vRLI, vRB, vRNI, NSX only) |VIDM_ONPREM_3330    |982
VMware NSX Data Center Professional    |VMware vRealize Log Insight 8.2.0 for NSX                                 |VRLI-820-NSX        |982
VMware NSX Data Center Professional    |VMware NSX-T Data Center 3.1.0                                            |NSX-T-310           |982
VMware NSX Data Center Standard        |VMware Identity Manager 3.3.3 (for vRA, vRops, vRLI, vRB, vRNI, NSX only) |VIDM_ONPREM_3330    |982
VMware NSX Data Center Standard        |VMware vRealize Log Insight 8.2.0 for NSX                                 |VRLI-820-NSX        |982
VMware NSX Data Center Standard        |VMware NSX-T Data Center 3.1.0                                            |NSX-T-310           |982
...
Get details on a specific NSX-T release...
2021/01/27 11:37:32 File Name: nsx-unified-appliance-3.0.1.1.0.16556500.ova
2021/01/27 11:37:32 Download URL: https://download2.vmware.com/software/NST-T30110/nsx-unified-appliance-3.0.1.1.0.16556500.ova?HashKey=e4334cfdce6ac5177d95a8ad1c000e13&params=%7B%22custnumber%22%3A%22dyVwKmV3ZXR0dA%3D%3D%22%2C%22sourcefilesize%22%3A%2211.02+GB%22%2C%22dlgcode%22%3A%22NSX-T-30110%22%2C%22languagecode%22%3A%22en%22%2C%22source%22%3A%22DOWNLOADS%22%2C%22downloadtype%22%3A%22manual%22%2C%22eula%22%3A%22Y%22%2C%22downloaduuid%22%3A%22992efa02-8344-475f-8efe-5ad7dd9070d3%22%2C%22purchased%22%3A%22Y%22%2C%22dlgtype%22%3A%22Product+Binaries%22%2C%22productversion%22%3A%223.0.1.1%22%2C%22productfamily%22%3A%22VMware+NSX-T+Data+Center%22%7D&AuthKey=1611744752_51d930e17cdc71a423b8ea3c71449521
2021/01/27 11:37:32 **************************************
2021/01/27 11:37:32 Downloading...
Downloading... 12 GB complete2021/01/27 12:23:00 Downloaded...
adeleporte@adeleporte-a02 test-govmware %
```
