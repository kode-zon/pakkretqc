package almsdk

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"
)

func NewALMError(resp *http.Response) error {
	var almerr ALMError
	var err error
	almerr.Body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	almerr.Code = resp.StatusCode
	return &almerr
}

type ALMError struct {
	Body []byte
	Code int
}

func (a *ALMError) Error() string {
	return fmt.Sprintf("%d: %s", a.Code, string(a.Body))
}

type Client struct {
	client *http.Client
	config *ClientOptions
	token  string
}

type ClientOptions struct {
	Endpoint string
}

func New(opt *ClientOptions) *Client {
	cookieJar, _ := cookiejar.New(nil)
	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &Client{
		client: &http.Client{
			Transport:     transCfg,
			CheckRedirect: http.DefaultClient.CheckRedirect,
			Jar:           cookieJar,
		},
		config: opt,
	}
}

func join(endpoint string, pathname ...string) *url.URL {
	u, err := url.Parse(endpoint)
	if err != nil {
		panic(fmt.Sprintf("invalid url for alm endpoint %s: %+v", endpoint, err))
	}
	u.Path = path.Join("qcbin/api", path.Join(pathname...))
	return u
}

func joinRest(endpoint string, pathname ...string) *url.URL {
	u, err := url.Parse(endpoint)
	if err != nil {
		panic(fmt.Sprintf("invalid url for alm endpoint %s: %+v", endpoint, err))
	}
	u.Path = path.Join("qcbin/rest", path.Join(pathname...))
	return u
}

type sessionCookieContext struct{}

func (c *Client) setTokenToRequest(ctx context.Context, req *http.Request) {
	if token, ok := ctx.Value(sessionCookieContext{}).(string); ok {
		req.Header.Set("Authorization", "Basic "+token)
		c.Authenticate(ctx, token)
	}
}

var InvalidCredential = fmt.Errorf("invalid credential")

func (c *Client) Authenticate(ctx context.Context, authtoken string) error {
	req, err := http.NewRequest("POST", join(c.config.Endpoint, "authentication/sign-in").String(), nil)
	req.Header.Add("Content-Type", "application/xml")
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", authtoken))

	log.Printf("Authenticate :: URL :: %v", req.URL)

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		respb, _ := ioutil.ReadAll(resp.Body)
		if resp.StatusCode == 401 {
			return InvalidCredential
		}
		return fmt.Errorf("login return %d: %s", resp.StatusCode, string(respb))
	}
	c.token = authtoken
	c.client.Jar.SetCookies(join(c.config.Endpoint), resp.Cookies())
	return nil
}

type Projects struct {
	Name string `json:"name"`
}

func (c *Client) Projects(ctx context.Context, domain string) ([]Projects, error) {
	var req, err = http.NewRequest("GET", join(c.config.Endpoint, "domains", domain, "projects").String(), nil)
	req.Header.Set("Accept", "application/json")
	c.setTokenToRequest(ctx, req)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 300 {
		message, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected error: /domains/%s/projects return %d status\n%s", domain, resp.StatusCode, string(message))
	}

	type respbody struct {
		Results []Projects `json:"results"`
	}
	var body respbody
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return nil, err
	}
	return body.Results, nil
}

type Domains struct {
	Name     string        `json:"name"`
	Projects []interface{} `json:"projects"`
}

func (c *Client) Domains(ctx context.Context) ([]*Domains, error) {
	var req, _ = http.NewRequest("GET", join(c.config.Endpoint, "domains").String(), nil)
	req.Header.Set("Accept", "application/json")
	c.setTokenToRequest(ctx, req)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 300 {
		return nil, fmt.Errorf("unexpected error: /domains return %d status", resp.StatusCode)
	}

	type respbody struct {
		Results []*Domains `json:"results"`
	}
	var body respbody
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return nil, err
	}
	return body.Results, nil
}

type Defect struct {
	DevComments  string   `json:"dev-comments"`
	Description  string   `json:"description"`
	LastModified *ALMTime `json:"last-modified"`
	CreationTime *ALMTime `json:"creation-time"`
	Status       string   `json:"status"`
	StatusExt    string   `json:"user-46"`
	Owner        string   `json:"owner"`
	Severity     string   `json:"severity"`
	DetectedBy   string   `json:"detected-by"`
	Name         string   `json:"name"`
	ID           int      `json:"id"`
}

const almTimeLayout = "2006-01-02 15:04:05"
const almDateLayout = "2006-01-02"

type ALMTime struct {
	t time.Time
}

func (a *ALMTime) UnmarshalJSON(b []byte) error {
	var err error
	var layout = almTimeLayout
	var data = strings.Trim(string(b), "\"")
	if len(strings.Split(data, " ")) == 1 {
		layout = almDateLayout
	}

	if len(strings.Split(data, ":")) == 2 {
		data = data + ":00"
	}

	if a.t, err = time.Parse(layout, data); err != nil {
		return err
	}
	return nil
}

func (a *ALMTime) Time() time.Time {
	return a.t
}

func (a *ALMTime) MarshalJSON() ([]byte, error) {
	return a.t.MarshalJSON()
}

func (c *Client) Defect(ctx context.Context, domain, project, id string, orignReq *http.Request) (*Defect, error) {

	log.Printf("Defect :: orignReq.Method :: %+v", orignReq.Method)
	//  https://stackoverflow.com/questions/38099501/updating-qc-alm-defect-comments-section-using-rest-api
	// if ctx == "PUT" {
	// 	var jsonStr = []byte(`{ "Fields": [{ "Name": "dev-comments", "values": [{ "value": "htmltext <div>\"string\"</div>  " }] }] }`)
	// 	var req, _ = http.NewRequest("GET", join(c.config.Endpoint, "domains", domain, "projects", project, "defects", id).String(), reqBody)
	// }
	q := orignReq.URL.Query()
	dbg := q.Get("debug")
	log.Printf("Defect :: dbg = %v", dbg)

	var req, _ = http.NewRequest("GET", join(c.config.Endpoint, "domains", domain, "projects", project, "defects", id).String(), nil)
	req.Header.Set("Accept", "application/json")
	c.setTokenToRequest(ctx, req)

	log.Printf("Defect :: URL :: %v", req.URL)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode > 300 {
		message, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected error: /domains/%s/projects/%s/defects return %d status\n%s", domain, project, resp.StatusCode, string(message))
	}
	var deflect Defect
	//err = json.NewDecoder(resp.Body).Decode(&deflect)

	bodyBuffer, _ := ioutil.ReadAll(resp.Body)
	err = json.NewDecoder(bytes.NewBuffer(bodyBuffer)).Decode(&deflect)

	if dbg != "" {
		log.Printf("Defect :: debug response body :: %v", string(bodyBuffer))
	}

	return &deflect, nil
}

func (c *Client) Defects(ctx context.Context, domain, project string, limit, offset int, orderFlag string, qq string) ([]*Defect, int, error) {

	var req, _ = http.NewRequest("GET", join(c.config.Endpoint, "domains", domain, "projects", project, "defects").String(), nil)
	req.Header.Set("Accept", "application/json")
	q := req.URL.Query()
	q.Add("order-by", orderFlag)
	q.Add("limit", strconv.Itoa(limit))
	if qq != "" {
		if strings.HasPrefix(qq, "{") {
			q.Add("query", qq)
		} else {
			q.Add("query", fmt.Sprintf("\"%s\"", qq))
		}
	}
	req.URL.RawQuery = q.Encode()
	c.setTokenToRequest(ctx, req)
	//fmt.Println(req.URL.String())
	log.Printf("Defects :: URL :: %v", req.URL)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 300 {
		message, _ := ioutil.ReadAll(resp.Body)
		return nil, 0, fmt.Errorf("unexpected error: /domains/%s/projects/%s/defects return %d status\n%s", domain, project, resp.StatusCode, string(message))
	}

	type respbody struct {
		Data []*Defect `json:"data"`
	}
	var body respbody
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return nil, 0, err
	}

	return body.Data, 0, nil

}

type Attachment struct {
	Type         string      `json:"type"`
	LastModified string      `json:"last-modified"`
	VcCurVer     interface{} `json:"vc-cur-ver"`
	VcUserName   interface{} `json:"vc-user-name"`
	Name         string      `json:"name"`
	FileSize     int         `json:"file-size"`
	RefSubtype   int         `json:"ref-subtype"`
	Description  interface{} `json:"description"`
	ID           int         `json:"id"`
	RefType      string      `json:"ref-type"`
	Entity       struct {
		ID   int    `json:"id"`
		Type string `json:"type"`
	} `json:"entity"`
}

func IsNetworkError(err error) {
}

func (c *Client) Attachments(ctx context.Context, domain, project string, query string, limit, offset int) ([]*Attachment, error) {
	var req, _ = http.NewRequest("GET", join(c.config.Endpoint, "domains", domain, "projects", project, "attachments").String(), nil)
	q := req.URL.Query()
	log.Printf("Attachments :: URL :: %v", req.URL)
	q.Add("query", fmt.Sprintf("\"%s\"", query))
	if limit <= 0 {
		limit = 20
	}
	if offset <= 0 {
		offset = 0
	}
	q.Add("limit", strconv.Itoa(limit))
	q.Add("offset", strconv.Itoa(offset))
	req.URL.RawQuery = q.Encode()
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, NewALMError(resp)
	}
	type respbody struct {
		Data []*Attachment `json:"data"`
	}
	var body respbody
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return nil, err
	}
	return body.Data, nil
}

func (c *Client) Attachment(ctx context.Context, domain, project string, id string, w io.Writer) error {
	var req, _ = http.NewRequest("GET", join(c.config.Endpoint, "domains", domain, "projects", project, "attachments", id).String(), nil)
	log.Printf("Attachment :: URL :: %v", req.URL)
	c.setTokenToRequest(ctx, req)
	req.Header.Set("Accept", "application/octet-stream")
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	mediaType, params, _ := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	log.Printf("mediaType = %v", mediaType)
	if strings.HasPrefix(mediaType, "multipart") {
		mr := multipart.NewReader(resp.Body, params["boundary"])
		for part, err := mr.NextPart(); err == nil; part, err = mr.NextPart() {
			if part.FormName() == "content" {
				_, err = io.Copy(w, part)
				if err != nil {
					return err
				}
			}
		}
	} else {
		_, err = io.Copy(w, resp.Body)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) ListItems(ctx context.Context, domain, project string, w io.Writer) error {
	var req, _ = http.NewRequest("GET", join(c.config.Endpoint, "domains", domain, "projects", project, "list-items").String(), nil)
	log.Printf("ListItems :: URL :: %v", req.URL)
	c.setTokenToRequest(ctx, req)

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return NewALMError(resp)
	}

	bodyBuffer, _ := ioutil.ReadAll(resp.Body)
	// if dbg != "" {
	// 	log.Printf("ListItems :: debug response body :: %v", string(bodyBuffer))
	// }
	_, err = io.Copy(w, bytes.NewBuffer(bodyBuffer))
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) TestInstances(ctx context.Context, domain, project, id string, w io.Writer) error {
	var req, _ = http.NewRequest("GET", join(c.config.Endpoint, "domains", domain, "projects", project, "test-instances", id).String(), nil)
	log.Printf("TestInstances :: URL :: %v", req.URL)
	c.setTokenToRequest(ctx, req)

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return NewALMError(resp)
	}

	bodyBuffer, _ := ioutil.ReadAll(resp.Body)
	// if dbg != "" {
	// 	log.Printf("ListItems :: debug response body :: %v", string(bodyBuffer))
	// }
	_, err = io.Copy(w, bytes.NewBuffer(bodyBuffer))
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Metadata(ctx context.Context, domain, project, collection_name string, w io.Writer) (string, error) {
	var req, _ = http.NewRequest("GET", join(c.config.Endpoint, "domains", domain, "projects", project, collection_name, "$metadata", "fields").String(), nil)
	log.Printf("Metadata :: URL :: %v", req.URL)
	c.setTokenToRequest(ctx, req)
	req.Header.Set("Accept", "application/json, */*")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", NewALMError(resp)
	}

	bodyBuffer, _ := ioutil.ReadAll(resp.Body)
	// if dbg != "" {
	// 	log.Printf("Defect :: debug response body :: %v", string(bodyBuffer))
	// }
	_, err = io.Copy(w, bytes.NewBuffer(bodyBuffer))
	if err != nil {
		return "", err
	}

	return "", nil
}

/*








#                   888      d8b                 d88P                        888
#                   888      Y8P                d88P                         888
#                   888                        d88P                          888
#  .d88888  .d8888b 88888b.  888 88888b.      d88P  888d888 .d88b.  .d8888b  888888
# d88" 888 d88P"    888 "88b 888 888 "88b    d88P   888P"  d8P  Y8b 88K      888
# 888  888 888      888  888 888 888  888   d88P    888    88888888 "Y8888b. 888
# Y88b 888 Y88b.    888 d88P 888 888  888  d88P     888    Y8b.          X88 Y88b.
#  "Y88888  "Y8888P 88888P"  888 888  888 d88P      888     "Y8888   88888P'  "Y888
#      888
#      888
#      888

// ALM  12.5x REST API Reference (Technical Preview)
// https://admhelp.microfocus.com/alm/en/12.55/api_refs/REST_TECH_PREVIEW/ALM_REST_API_TP.html#REST_API_Tech_Preview/REST/entity.html%3FTocPath%3DResource%2520Reference%7CCustomization%2520Resources%7C_____1
*/

func (c *Client) FieldsCollection(ctx context.Context, domain, project, entity_name string, orignReq *http.Request, w io.Writer) (string, error) {
	// ALM  12.5x REST API Reference (Technical Preview)
	var req, _ = http.NewRequest("GET", joinRest(c.config.Endpoint, "domains", domain, "projects", project, "customization", "entities", entity_name, "fields").String(), nil)
	q := orignReq.URL.Query()
	req.URL.RawQuery = q.Encode()
	c.setTokenToRequest(ctx, req)

	log.Printf("FieldsCollection :: URL :: %v", req.URL)
	c.setTokenToRequest(ctx, req)
	req.Header.Set("Accept", "application/json, */*")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", NewALMError(resp)
	}

	bodyBuffer, _ := ioutil.ReadAll(resp.Body)
	// if dbg != "" {
	// 	log.Printf("Defect :: debug response body :: %v", string(bodyBuffer))
	// }
	_, err = io.Copy(w, bytes.NewBuffer(bodyBuffer))
	if err != nil {
		return "", err
	}

	return "", nil
}

func (c *Client) ListsRelatedToEntity(ctx context.Context, domain, project, entity_name string, orignReq *http.Request, w io.Writer) (string, error) {
	// ALM  12.5x REST API Reference (Technical Preview)
	var req, _ = http.NewRequest("GET", joinRest(c.config.Endpoint, "domains", domain, "projects", project, "customization", "entities", entity_name, "lists").String(), nil)
	q := orignReq.URL.Query()
	req.URL.RawQuery = q.Encode()
	c.setTokenToRequest(ctx, req)

	log.Printf("FieldsCollection :: URL :: %v", req.URL)
	c.setTokenToRequest(ctx, req)
	req.Header.Set("Accept", "application/json, */*")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", NewALMError(resp)
	}

	bodyBuffer, _ := ioutil.ReadAll(resp.Body)
	// if dbg != "" {
	// 	log.Printf("Defect :: debug response body :: %v", string(bodyBuffer))
	// }
	_, err = io.Copy(w, bytes.NewBuffer(bodyBuffer))
	if err != nil {
		return "", err
	}

	return "", nil
}

func (c *Client) PermissionCollection(ctx context.Context, domain, project, entity_name string, orignReq *http.Request, w io.Writer) (string, error) {
	// ALM  12.5x REST API Reference (Technical Preview)
	var req, _ = http.NewRequest("GET", joinRest(c.config.Endpoint, "domains", domain, "projects", project, "customization", "entities", entity_name, "permissions").String(), nil)
	q := orignReq.URL.Query()
	req.URL.RawQuery = q.Encode()
	c.setTokenToRequest(ctx, req)

	log.Printf("FieldsCollection :: URL :: %v", req.URL)
	c.setTokenToRequest(ctx, req)
	req.Header.Set("Accept", "application/json, */*")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", NewALMError(resp)
	}

	bodyBuffer, _ := ioutil.ReadAll(resp.Body)
	// if dbg != "" {
	// 	log.Printf("Defect :: debug response body :: %v", string(bodyBuffer))
	// }
	_, err = io.Copy(w, bytes.NewBuffer(bodyBuffer))
	if err != nil {
		return "", err
	}

	return "", nil
}
func (c *Client) RelationsCollection(ctx context.Context, domain, project, entity_name string, orignReq *http.Request, w io.Writer) (string, error) {
	// ALM  12.5x REST API Reference (Technical Preview)
	var req, _ = http.NewRequest("GET", joinRest(c.config.Endpoint, "domains", domain, "projects", project, "customization", "entities", entity_name, "relations").String(), nil)
	q := orignReq.URL.Query()
	req.URL.RawQuery = q.Encode()
	c.setTokenToRequest(ctx, req)

	log.Printf("FieldsCollection :: URL :: %v", req.URL)
	c.setTokenToRequest(ctx, req)
	req.Header.Set("Accept", "application/json, */*")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", NewALMError(resp)
	}

	bodyBuffer, _ := ioutil.ReadAll(resp.Body)
	// if dbg != "" {
	// 	log.Printf("Defect :: debug response body :: %v", string(bodyBuffer))
	// }
	_, err = io.Copy(w, bytes.NewBuffer(bodyBuffer))
	if err != nil {
		return "", err
	}

	return "", nil
}

func (c *Client) TypesCollection(ctx context.Context, domain, project, entity_name string, orignReq *http.Request, w io.Writer) (string, error) {
	// ALM  12.5x REST API Reference (Technical Preview)
	var req, _ = http.NewRequest("GET", joinRest(c.config.Endpoint, "domains", domain, "projects", project, "customization", "entities", entity_name, "types").String(), nil)
	q := orignReq.URL.Query()
	req.URL.RawQuery = q.Encode()
	c.setTokenToRequest(ctx, req)

	log.Printf("FieldsCollection :: URL :: %v", req.URL)
	c.setTokenToRequest(ctx, req)
	req.Header.Set("Accept", "application/json, */*")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", NewALMError(resp)
	}

	bodyBuffer, _ := ioutil.ReadAll(resp.Body)
	// if dbg != "" {
	// 	log.Printf("Defect :: debug response body :: %v", string(bodyBuffer))
	// }
	_, err = io.Copy(w, bytes.NewBuffer(bodyBuffer))
	if err != nil {
		return "", err
	}

	return "", nil
}

func (c *Client) TypeSubtypesFieldCollection(ctx context.Context, domain, project, entity_name, subtype_id string, orignReq *http.Request, w io.Writer) (string, error) {
	// ALM  12.5x REST API Reference (Technical Preview)
	var req, _ = http.NewRequest("GET", joinRest(c.config.Endpoint, "domains", domain, "projects", project, "customization", "entities", entity_name, "types", subtype_id, "fields").String(), nil)
	q := orignReq.URL.Query()
	req.URL.RawQuery = q.Encode()
	c.setTokenToRequest(ctx, req)

	log.Printf("FieldsCollection :: URL :: %v", req.URL)
	c.setTokenToRequest(ctx, req)
	req.Header.Set("Accept", "application/json, */*")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", NewALMError(resp)
	}

	bodyBuffer, _ := ioutil.ReadAll(resp.Body)
	// if dbg != "" {
	// 	log.Printf("Defect :: debug response body :: %v", string(bodyBuffer))
	// }
	_, err = io.Copy(w, bytes.NewBuffer(bodyBuffer))
	if err != nil {
		return "", err
	}

	return "", nil
}

func (c *Client) EntityMetadata(ctx context.Context, domain, project, entity_name string, orignReq *http.Request, w io.Writer) (string, error) {
	// ALM  12.5x REST API Reference (Technical Preview)
	var req, _ = http.NewRequest("GET", joinRest(c.config.Endpoint, "domains", domain, "projects", project, "customization", "entities", entity_name).String(), nil)
	q := orignReq.URL.Query()
	req.URL.RawQuery = q.Encode()
	c.setTokenToRequest(ctx, req)

	log.Printf("FieldsCollection :: URL :: %v", req.URL)
	c.setTokenToRequest(ctx, req)
	req.Header.Set("Accept", "application/json, */*")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", NewALMError(resp)
	}

	bodyBuffer, _ := ioutil.ReadAll(resp.Body)
	// if dbg != "" {
	// 	log.Printf("Defect :: debug response body :: %v", string(bodyBuffer))
	// }
	_, err = io.Copy(w, bytes.NewBuffer(bodyBuffer))
	if err != nil {
		return "", err
	}

	return "", nil
}

func (c *Client) EntitiesCollection(ctx context.Context, domain, project string, orignReq *http.Request, w io.Writer) (string, error) {
	// ALM  12.5x REST API Reference (Technical Preview)
	var req, _ = http.NewRequest("GET", joinRest(c.config.Endpoint, "domains", domain, "projects", project, "customization", "entities").String(), nil)
	q := orignReq.URL.Query()
	req.URL.RawQuery = q.Encode()
	c.setTokenToRequest(ctx, req)

	log.Printf("FieldsCollection :: URL :: %v", req.URL)
	c.setTokenToRequest(ctx, req)
	req.Header.Set("Accept", "application/json, */*")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", NewALMError(resp)
	}

	bodyBuffer, _ := ioutil.ReadAll(resp.Body)
	// if dbg != "" {
	// 	log.Printf("Defect :: debug response body :: %v", string(bodyBuffer))
	// }
	_, err = io.Copy(w, bytes.NewBuffer(bodyBuffer))
	if err != nil {
		return "", err
	}

	return "", nil
}

func (c *Client) UsedListsCollection(ctx context.Context, domain, project string, orignReq *http.Request, w io.Writer) (string, error) {
	// ALM  12.5x REST API Reference (Technical Preview)
	var req, _ = http.NewRequest("GET", joinRest(c.config.Endpoint, "domains", domain, "projects", project, "customization", "used-lists").String(), nil)
	q := orignReq.URL.Query()
	req.URL.RawQuery = q.Encode()
	c.setTokenToRequest(ctx, req)

	log.Printf("UsedListsCollection :: URL :: %v", req.URL)
	c.setTokenToRequest(ctx, req)
	req.Header.Set("Accept", "application/json, */*")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", NewALMError(resp)
	}

	bodyBuffer, _ := ioutil.ReadAll(resp.Body)
	// if dbg != "" {
	// 	log.Printf("Defect :: debug response body :: %v", string(bodyBuffer))
	// }
	_, err = io.Copy(w, bytes.NewBuffer(bodyBuffer))
	if err != nil {
		return "", err
	}

	return "", nil
}

func (c *Client) UserCollection(ctx context.Context, domain, project string, orignReq *http.Request, w io.Writer) (string, error) {
	// ALM  12.5x REST API Reference (Technical Preview)
	var req, _ = http.NewRequest("GET", joinRest(c.config.Endpoint, "domains", domain, "projects", project, "customization", "users").String(), nil)
	q := orignReq.URL.Query()
	req.URL.RawQuery = q.Encode()
	c.setTokenToRequest(ctx, req)

	log.Printf("UserCollection :: URL :: %v", req.URL)
	c.setTokenToRequest(ctx, req)
	req.Header.Set("Accept", "application/json, */*")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", NewALMError(resp)
	}

	bodyBuffer, _ := ioutil.ReadAll(resp.Body)
	// if dbg != "" {
	// 	log.Printf("Defect :: debug response body :: %v", string(bodyBuffer))
	// }
	_, err = io.Copy(w, bytes.NewBuffer(bodyBuffer))
	if err != nil {
		return "", err
	}

	return "", nil
}
