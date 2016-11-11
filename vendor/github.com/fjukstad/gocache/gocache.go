package gocache

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

// Time until cache entry is invalidated. Default 24h.
var cacheInvalidation = "24h"

type Entry struct {
	Response *http.Response
	Content  string
}

func Get(url string) (resp *http.Response, err error) {

	resp, err = getFromCache(url)

	if err != nil {
		return getFromWeb(url)
	}

	return
}

// Set time until cache entry is invalidated such as "300h", "1.5s" or "2h45m".
// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h". (anything
// time.ParseDuration can parse)
func SetInvalidationTime(time string) {
	cacheInvalidation = time
	return
}

// Fetch contents of url from web and write to cache
func getFromWeb(url string) (resp *http.Response, err error) {
	resp, err = http.Get(url)

	if err != nil {
		return
	}
	var resp2 http.Response
	resp2 = *resp

	defer resp.Body.Close()

	// Draining the body
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Print("Reading response body went bad. ", err)

	}

	writeToCache(url, body, &resp2)
	resp.Body = nopCloser{bytes.NewBufferString(string(body))}

	return
}

func writeToCache(url string, body []byte, resp *http.Response) {

	cacheEntry := generateCacheEntry(resp, body)

	b, err := json.Marshal(cacheEntry)

	if err != nil {
		fmt.Println("Could not marshal response ", err)
		return
	}

	filename := getFilePath(url)

	// Set up any directory needed to write the file.
	// e.g. vg.no/first/second/file will be stored as
	// cache/first/second/file

	err = createDirectories(filename)
	if err != nil {
		pe, _ := err.(*os.PathError)

		if !strings.Contains(pe.Error(), "file exists") {
			fmt.Println("Could not create directories ", err)
			return
		}
	}

	// If the file hasn't got an extension, set it to .json
	name := strings.Split(filename, "/")
	fn := name[len(name)-1]

	if len(strings.Split(fn, ".")) < 2 {
		filename = filename + ".json"
	}

	// Create file
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("could not create cache file ", err)
		return
	}

	// Close and check for error on exit
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println("Could not close file ", err)
		}
	}()

	// Write file
	_, err = file.Write(b)
	if err != nil {
		fmt.Println("Could not write to cache file ", err)
	}

}

func generateCacheEntry(resp *http.Response, body []byte) Entry {
	Response := resp
	Content := string(body)
	entry := Entry{Response, Content}

	// Cannot marshal the body from get go
	entry.Response.Body = nil

	// we don't care for the original request.
	// commented the line below (14:55 1.8.2016 work on luftkvalitet).
	// don't know why I set it to nil...

	entry.Response.Request = nil
	//entry.Response.Request.Cancel = nil

	return entry

}

// Try to fetch contents of addr from cache
func getFromCache(URL string) (resp *http.Response, err error) {

	filename := getFilePath(URL)

	// If the file hasn't got an extension, set it to .json
	name := strings.Split(filename, "/")
	fn := name[len(name)-1]

	if len(strings.Split(fn, ".")) < 2 {
		filename = filename + ".json"
	}

	file, err := os.Open(filename)
	if err != nil {
		err = errors.New("File '" + filename + "' not Found")
		return
	}

	defer file.Close()

	entry, err := readFromFile(file)

	if err != nil {
		return
	}

	resp = entry.GenerateHttpResponse()

	// FIRST CHECK IF OLD ENTRY
	respTime := resp.Header.Get("Date")

	loc, err := time.LoadLocation("Local")
	if err != nil {
		fmt.Println("Could not load location", err)
		return nil, err
	}

	t0, err := time.ParseInLocation("Mon, 2 Jan 2006 15:04:05 MST", respTime, loc)
	if err != nil {
		fmt.Println("Could not parse time", err, respTime)
		return nil, err
	}

	waitTime, err := time.ParseDuration(cacheInvalidation)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	t0 = t0.Add(waitTime)

	if time.Now().After(t0) {
		fmt.Println("Cache entry was too old, fetching new from", URL)
		return getFromWeb(URL)
	}

	// Uodate response.request with applicable fields.
	resp.Request = &http.Request{}
	u, _ := url.Parse("")
	resp.Request.URL = u

	return resp, nil

}

// Read cache entry from file
func readFromFile(file *os.File) (entry *Entry, err error) {
	fileinfo, err := file.Stat()
	var size int
	size = int(fileinfo.Size())

	buf := make([]byte, size)

	_, err = file.Read(buf)

	if err != nil {
		return entry, errors.New("Reading file went wrong")
	}

	// Must trim buffer before unmarshaling it. This is because of
	// the unmarshaling failing if entire buffer is returned
	buf = bytes.Trim(buf[0:], string(0))

	entry = new(Entry)
	err = json.Unmarshal(buf, entry)

	if err != nil {
		fmt.Print("Unmarshaling gone wrong: ", err)
		return entry, errors.New("Unmarshaling gone wrong")
	}

	return
}

func (entry *Entry) Print() {
	fmt.Print("Content: ", entry.Content)
}

func (entry *Entry) GenerateHttpResponse() (resp *http.Response) {

	resp = entry.Response
	resp.Body = nopCloser{bytes.NewBufferString(entry.Content)}

	return resp

}

func getFilePath(url string) (dir string) {
	urlTokens := strings.Split(url, "/")
	// 2 because we need to strip away 'http:', ' '
	strippedUrl := urlTokens[2:]
	dir = "cache/" + strings.Join(strippedUrl, "/")
	return
}

func createDirectories(filename string) error {

	filepath := path.Dir(filename)
	directories := strings.Split(filepath, "/")

	p := ""
	for i := range directories {
		p += directories[i] + "/"
		err := os.Mkdir(p, 0755)

		if err != nil {
			pe, _ := err.(*os.PathError)

			// if folder exists, continue to the next one
			if !strings.Contains(pe.Error(), "file exists") {
				fmt.Println("Mkdir failed miserably: ", err)
				return err
			}
		}
	}

	return nil
}

/* Below are from:
   https://groups.google.com/forum/#!topic/golang-nuts/J-Y4LtdGNSw
*/

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error {
	return nil
}
