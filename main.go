package gomyvmware

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
	"golang.org/x/net/publicsuffix"
)

func NewClient(username, password string) (MyClient, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	jar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})

	client := &http.Client{
		Jar: jar,
	}

	client.Get("https://my.vmware.com/web/vmware/login")

	data := url.Values{}
	data.Set("username", username)
	data.Set("password", password)

	r, err := http.NewRequest("POST", "https://auth.vmware.com/oam/server/auth_cred_submit?Auth-AppID=WMVMWR", strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatal(err)
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	res, _ := client.Do(r)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	body_text := fmt.Sprintf("%s", body)

	body_text = strings.ReplaceAll(body_text, "\n", "%0D%0A")
	body_text = strings.ReplaceAll(body_text, "\r", "")
	body_text = strings.ReplaceAll(body_text, "+", "%2B")

	re := regexp.MustCompile(`NAME="SAMLResponse" VALUE="(.+)"`)
	saml_response := re.FindString(body_text)
	saml_response = strings.TrimPrefix(saml_response, `NAME="SAMLResponse" VALUE=`)
	saml_response = fmt.Sprintf("SAMLResponse=%s", saml_response)
	saml_response = strings.ReplaceAll(saml_response, `"`, "")

	r2, err := http.NewRequest("POST", "https://my.vmware.com/vmwauth/saml/SSO", strings.NewReader(saml_response))
	if err != nil {
		log.Fatal(err)
	}
	r2.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r2.Header.Add("Content-Length", strconv.Itoa(len(saml_response)))

	res2, _ := client.Do(r2)
	if err != nil {
		log.Fatal(err)
	}

	u, err := url.Parse("https://my.vmware.com")

	XSRF_token := http.Cookie{}
	for _, cookie := range client.Jar.Cookies(u) {
		if cookie.Name == "XSRF-TOKEN" {
			XSRF_token = *cookie
			//fmt.Printf("  %s: %s\n", cookie.Name, cookie.Value)
		}
	}

	defer res2.Body.Close()

	return MyClient{
		Client:      client,
		XSRF_Cookie: XSRF_token,
	}, nil

}

func GetProductsAtoZ(client MyClient) (ProductsAtoZ, error) {
	dat := ProductsAtoZ{}

	r, err := http.NewRequest("GET", "https://my.vmware.com/channel/public/api/v1.0/products/getProductsAtoZ", nil)
	if err != nil {
		return dat, err
	}
	r.Header.Add("X-XSRF-TOKEN", client.XSRF_Cookie.Value)

	res, _ := client.Client.Do(r)
	if err != nil {
		return dat, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err := json.Unmarshal(body, &dat); err != nil {
		return dat, err
	}

	return dat, nil

}

func GetRelatedDLGList(client MyClient, params DLGListParams) (DLGListResults, error) {
	dat := DLGListResults{}

	r, err := http.NewRequest("GET", fmt.Sprintf("https://my.vmware.com/channel/public/api/v1.0/products/getRelatedDLGList?category=%s&product=%s&version=%s&dlgType=%s", params.Category, params.Product, params.Version, params.DlgType), nil)
	if err != nil {
		return dat, err
	}
	r.Header.Add("X-XSRF-TOKEN", client.XSRF_Cookie.Value)

	res, _ := client.Client.Do(r)
	if err != nil {
		return dat, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err := json.Unmarshal(body, &dat); err != nil {
		return dat, err
	}

	return dat, nil

}

func GetDLGHeader(client MyClient, params GetDLGParams) (GetDLGHeaderResults, error) {
	dat := GetDLGHeaderResults{}

	r, err := http.NewRequest("GET", fmt.Sprintf("https://my.vmware.com/channel/public/api/v1.0/products/getDLGHeader?downloadGroup=%s&productId=%d", params.DownloadGroup, params.ProductID), nil)
	if err != nil {
		return dat, err
	}
	r.Header.Add("X-XSRF-TOKEN", client.XSRF_Cookie.Value)

	res, _ := client.Client.Do(r)
	if err != nil {
		return dat, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err := json.Unmarshal(body, &dat); err != nil {
		return dat, err
	}

	return dat, nil

}

func GetDLGDetails(client MyClient, params GetDLGParams) (GetDLGDetailsResults, error) {
	dat := GetDLGDetailsResults{}

	r, err := http.NewRequest("GET", fmt.Sprintf("https://my.vmware.com/channel/api/v1.0/dlg/details?downloadGroup=%s&productId=%d", params.DownloadGroup, params.ProductID), nil)
	if err != nil {
		return dat, err
	}
	r.Header.Add("X-XSRF-TOKEN", client.XSRF_Cookie.Value)

	res, _ := client.Client.Do(r)
	if err != nil {
		return dat, err
	}

	//fmt.Printf("%s", r)

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err := json.Unmarshal(body, &dat); err != nil {
		return dat, err
	}

	return dat, nil

}

func (r GetDLGDetailsResults) Filter(filter string) (DownloadFile, error) {

	re := regexp.MustCompile(filter)

	for _, file := range r.DownloadFiles {
		if re.Match([]byte(file.FileName)) {
			return file, nil
		}
	}
	return DownloadFile{}, errors.New("File not found")
}

func GetEulaAccept(client MyClient, downloadGroup string) ([]byte, error) {

	r, err := http.NewRequest("GET", fmt.Sprintf("https://my.vmware.com/channel/api/v1.0/dlg/eula/accept?isPrivate=true&downloadGroup=%s", downloadGroup), nil)
	if err != nil {
		return nil, err
	}
	r.Header.Add("X-XSRF-TOKEN", client.XSRF_Cookie.Value)

	res, _ := client.Client.Do(r)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	return body, nil

}

func GetDownload(client MyClient, get_download_body GetDownloadBody) (GetDownloadResults, error) {
	dat := GetDownloadResults{}

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(get_download_body)

	r, err := http.NewRequest("POST", "https://my.vmware.com/channel/api/v1.0/dlg/download", buf)
	if err != nil {
		return dat, err
	}
	r.Header.Add("X-XSRF-TOKEN", client.XSRF_Cookie.Value)
	r.Header.Add("Content-Type", "application/json")

	res, _ := client.Client.Do(r)
	if err != nil {
		return dat, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err := json.Unmarshal(body, &dat); err != nil {
		return dat, err
	}

	return dat, nil

}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

func (wc WriteCounter) PrintProgress() {
	// Clear the line by using a character return to go back to the start and remove
	// the remaining characters by filling it with spaces
	fmt.Printf("\r%s", strings.Repeat(" ", 35))

	// Return again and print current status of download
	// We use the humanize package to print the bytes in a meaningful way (e.g. 10 MB)
	fmt.Printf("\rDownloading... %s complete", humanize.Bytes(wc.Total))
}

func Download(client MyClient, dnl GetDownloadResults) error {

	// Create the file, but give it a tmp file extension, this means we won't overwrite a
	// file until it's downloaded, but we'll remove the tmp extension once downloaded.
	out, err := os.Create(dnl.FileName + ".tmp")
	if err != nil {
		return err
	}

	// Get the data
	resp, err := http.Get(dnl.DownloadURL)
	if err != nil {
		out.Close()
		return err
	}
	defer resp.Body.Close()

	// Create our progress reporter and pass it to be used alongside our writer
	counter := &WriteCounter{}
	if _, err = io.Copy(out, io.TeeReader(resp.Body, counter)); err != nil {
		out.Close()
		return err
	}

	// The progress use the same line so print a new line once it's finished downloading
	fmt.Print("\n")

	// Close the file without defer so it can happen before Rename()
	out.Close()

	if err = os.Rename(dnl.FileName+".tmp", "./"); err != nil {
		return err
	}
	return nil

}
