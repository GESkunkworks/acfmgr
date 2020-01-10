// Package acfmgr is a package to manage entries in an AWS credentials file
//
// Sample resulting AWS creds file entries:
//
//  [account2]
//  # DO NOT EDIT
//  # ACFMGR MANAGED SECTION
//  # (Will be overwritten regularly)
//  ####################################################
//  # ASSUMED ROLE: arn:aws:iam::098765432123:role/aj/d-readonly
//  # ASSUMED FROM INSTANCE ROLE: NA
//  # GENERATED: 2019-12-27 14:10:37.282148008 -0500 EST
//  # EXPIRES@   2019-12-27 20:10:14 +0000 UTC
//  output = json
//  region = us-east-1
//  aws_access_key_id = ASIASDIVWOEIOBINAIE
//  aws_secret_access_key = GyD6rud3Q06qk90pTlECLncbbKx7GPXjM2N5ocVe
//  aws_session_token = FwoGZXIvYXdzED0aDNHw4GhQvSFSCn8vUCK6Af+KK2QGsRbN5F22xJvXyNyYoAzxTkPYrSgvvuL7/17tyBa5LMeHWSKV/9E3ON2vRSLIz0iFfeEE5cj4zmbqpw/5LAiDiptTvbQQKmzCE4Pt05khFcsTmwsju9ibR5Mx2oJKdHHQXCsqk0XjvugSuu+KbU0wigO2oSXvu1dguNg+j6RTdxGAS7Uoih2WZR4ZlJCdcFNOivhf/kWs18mMRQ43r47GWsV9Z3vlTaMimHLWuBMldPgBcJV2iCiWrpnwBTIt2Dfkgvi8Bs7OcInotWE751K48QJnzcwPMKjsNKBE0tf1kGI9JArO8x+aDQJX
//
//  [account1]
//  # DO NOT EDIT
//  # ACFMGR MANAGED SECTION
//  # (Will be overwritten regularly)
//  ####################################################
//  # ASSUMED ROLE: arn:aws:iam::123456789012:role/aj/d-admin
//  # ASSUMED FROM INSTANCE ROLE: NA
//  # GENERATED: 2019-12-27 14:10:37.332225334 -0500 EST
//  # EXPIRES@   2019-12-27 20:10:14 +0000 UTC
//  output = json
//  region = us-east-2
//  aws_access_key_id = ASIZIPVKAVLEIGH
//  aws_secret_access_key = GFIzAhjrBu9h3kkUNimntlN0zTF0TLmMebkMM0AP
//  aws_session_token = FwoGZXIvYXdzED0aDNnD7rODZ+iSKGSi6yK6ASyAn9hJz0ER81eYa842hUwQnca7QNeFq7ZOrYvKb3ZegoVSRFEOaMvgw5La/taN8udMAdFINmxxV7Fx7JGTWuMK9JbWQA8I/AHqCN/NmOC1PrbIRvUhAZ8FgTdjNbiyh8CoOEvFqI3n4uQ57oWG5EZGZh8DSfENoVANR1AIaod7sFU1yHnHKOlr5Zp/iIUcD5j8X8yY8m05Vj2JFNipwcUsIVTNCeaWMud5n/30F4g/sQJLrIcV3nNoTCiWrpnwBTItwq++PEN9OzQYkCIEFNcvJe2ZnkiPz/+4xDNSSfiBGsHMKCGdChINizxQQHzE
//  ...
//
//
// Adding and removing entries manually is a pain when doing large amounts of assume
// role operations so this package was created to assist in programattically adding
// them once you have sessions built from the Golang AWS SDK.
//
// To start building a credentials file all you have to do is start a NewCredFileSession
// with the desired filename. If the file does not exist acfmgr will create it fore you.
// Use the CredFile object that's returned to manage the entries in the file.
//
// From there you use ProfileEntryInput objects and the CredFile.NewEntry() to put new
// credentials in the queue. After you're done creating objects you can write them all
// to file with CredFile.AssertEntries or delete all of the queued profile names with
// CredFile.DeleteEntries.
//
// The package will not touch any existing profile entries unless they match the name 
// of a profile entry that you provide to the CredFile object. 
//
//
// Sample with minimum input params
//
//  c, err := acfmgr.NewCredFileSession("~/.aws/credentials")
//  check(err)
//  profileInput := acfmgr.ProfileEntryInput{
//      Credential: <some sts.Credentials object>,
//      ProfileEntryName: "devaccount1",
//  }
//  err = c.NewEntry(profileInput)
//  check(err)
//  err = c.AssertEntries()
//  check(err)
//
// Another Sample using more of the input parameters:
//  c, err := acfmgr.NewCredFileSession("~/.aws/credentials")
//  check(err)
//  profileInput := acfmgr.ProfileEntryInput{
//      Credential: <some sts.Credentials object>,
//      ProfileEntryName: "devaccount",
//      OutputFormat: "text",
//      Region: "us-east-2",
//      AssumeRoleARN: "arn:aws:iam::123456789012:role/aj/d-admin"
//  }
//  err = c.NewEntry(profileInput)
//  check(err)
//  err = c.DeleteEntries()
//
// Yields:
//     [devaccount]
//     # DO NOT EDIT
//     # ACFMGR MANAGED SECTION
//     # (Will be overwritten regularly)
//     ####################################################
//     # ASSUMED ROLE: arn:aws:iam::123456789012:role/aj/d-admin
//     # ASSUMED FROM INSTANCE ROLE: NA
//     # GENERATED: 2019-12-27 14:10:37.332225334 -0500 EST
//     # EXPIRES@2019-12-27 20:10:14 +0000 UTC
//     output = text
//     region = us-east-2
//     aws_access_key_id = ASIZIPVKAVLEIGHAODIEA
//     aws_secret_access_key = GFIzAhjrBu9h3kkUNimntlN0zTF0TLmMebkMM0AP
//     aws_session_token = FwoGZXIvYXdzED0aDNnD7rODZ+iSKGSi6yK6ASyAn9hJz0ER81eYa842hUwQnca7QNeFq7ZOrYvKb3ZegoVSRFEOaMvgw5La/taN8udMAdFINmxxV7Fx7JGTWuMK9JbWQA8I/AHqCN/NmOC1PrbIRvUhAZ8FgTdjNbiyh8CoOEvFqI3n4uQ57oWG5EZGZh8DSfENoVANR1AIaod7sFU1yHnHKOlr5Zp/iIUcD5j8X8yY8m05Vj2JFNipwcUsIVTNCeaWMud5n/30F4g/sQJLrIcV3nNoTCiWrpnwBTItwq++PEN9OzQYkCIEFNcvJe2ZnkiPz/+4xDNSSfiBGsHMKCGdChINizxQQHzE
//   ...
//
package acfmgr

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/service/sts"
	"io/ioutil"
	"os"
	"os/user"
	"regexp"
	"strings"
	"text/template"
	"time"
    "encoding/json"
)

var defaultTemplate *template.Template

const credFileTemplate = `# DO NOT EDIT
# ACFMGR MANAGED SECTION
# (Will be overwritten regularly)
####################################################
# ASSUMED ROLE: {{.AssumeRoleARN}}
# ASSUMED FROM INSTANCE ROLE: {{.InstanceRoleARN}}
# GENERATED: {{.Generated}}
{{ .ExpiresToken }}   {{.Expiration}}
{{- if .HasDescription}}
# DESCRIPTION: {{.Description}}{{end}}
{{- if .HasRegion}}
region = {{.Region}}{{end}}
{{- if .HasOutput}}
output = {{.Output}}{{end}}
aws_access_key_id = {{.AccessKeyID}}
aws_secret_access_key = {{.SecretAccessKey}}
aws_session_token = {{.SessionToken}}
`

func init() {
	// setup the default template hoping we don't introduce errors ourself
	defaultTemplate, _ = template.New("test").Parse(credFileTemplate)
}

// NewCredFileSession creates a new interactive credentials file
// session. Needs target filename and returns CredFile obj and err.
func NewCredFileSession(filename string) (cf *CredFile, err error) {
	usr, err := user.Current()
	if err != nil {
		return cf, err
	}
	// try to get absolute path of file
	filenameExpanded, err := expandPath(filename, usr)
	if err != nil {
		return cf, err
	}
	credfile := CredFile{filename: filenameExpanded,
		currBuff: new(bytes.Buffer),
		reSep:    regexp.MustCompile(`\[.*\]`),
	}
	err = credfile.loadFile()
	if err != nil {
		return cf, err
	}
	cf = &credfile
	return cf, err
}

// CredFile should be built with the exported
// NewCredFileSession function.
type CredFile struct {
	filename string
	ents     []*credEntry
	currBuff *bytes.Buffer
	reSep    *regexp.Regexp // regex cred anchor separator e.g. "[\w*]"
}

type credEntry struct {
	name     string
	contents []string
}

// addEntry adds a new credentials entry to the queue
// to be written or deleted with the AssertEntries or
// DeleteEntries method.
func (c *CredFile) addEntry(entryName string, entryContents []string) {
	e := credEntry{name: entryName, contents: entryContents}
	c.ents = append(c.ents, &e)
}

// AssertEntries loops through all of the credEntry objs
// attached to CredFile obj and makes sure there is an
// occurrence with the credEntry.name and contents.
// Existing entries of the same name with different
// contents will be clobbered.
func (c *CredFile) AssertEntries() (err error) {
	for _, e := range c.ents {
		err = c.modifyEntry(true, e)
		if err != nil {
			return err
		}
	}
	return err
}

// DeleteEntries loops through all of the credEntry
// objs attached to CredFile obj and makes sure entries
// with the same credEntry.name are removed. Will remove
// ALL entries with the same name.
func (c *CredFile) DeleteEntries() (err error) {
	for _, e := range c.ents {
		err = c.modifyEntry(false, e)
		if err != nil {
			return err
		}
	}
	return err
}

func (e *credEntry) appendToList(lister []string) []string {
	lister = append(lister, e.name)
	for _, line := range e.contents {
		lister = append(lister, line)
	}
	return lister
}

func (c *CredFile) loadFile() error {
	if !c.fileExists() {
		_, err := c.createFile()
		if err != nil {
			panic(err)
		}
	}
	f, err := os.OpenFile(c.filename, os.O_RDONLY, os.ModeAppend)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		c.currBuff.WriteString(scanner.Text() + "\n")
	}
	return err
}

func (c *CredFile) writeBufferToFile() error {
	err := ioutil.WriteFile(c.filename, c.currBuff.Bytes(), 0644)
	return err
}

// indexOf find the index of a value in an []int
func indexOf(s []int, e int) (index int) {
	for index, a := range s {
		if a == e {
			return index
		}
	}
	return -1
}

func (c *CredFile) removeEntry(data []string, anchors []int, entry *credEntry) []string {
	ignoring := false
	ignoreUntil := 0
	var newLines []string
	for i, line := range data {
		if line == entry.name {
			currIndex := indexOf(anchors, i)
			if (currIndex + 1) >= len(anchors) {
				// this means it's at EOF
				ignoreUntil = len(data)
			} else {
				ignoreUntil = anchors[currIndex+1]
			}

			ignoring = true
		}
		if !(ignoring && i < ignoreUntil) {
			newLines = append(newLines, line)
		}
	}
	return newLines
}

// EnsureEntryExists makes sure that the attached Ent
// entry exists.
func (c *CredFile) modifyEntry(replace bool, entry *credEntry) (err error) {
	found := false
	// read buffer into []string
	var lines []string
	scanner := bufio.NewScanner(c.currBuff)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	// search for entry
	var anchors []int
	for i, line := range lines {
		reMatch := c.reSep.FindAllString(line, -1)
		if reMatch != nil {
			anchors = append(anchors, i)
		}
		if line == entry.name {
			found = true
		}
	}
	switch {
	case found && replace:
		lines = c.removeEntry(lines, anchors, entry)
		// make the credEntry append itself to the results
		lines = entry.appendToList(lines)
	case found && !replace:
		lines = c.removeEntry(lines, anchors, entry)
	case !found && !replace:
		// do nothing
	case !found && replace:
		lines = entry.appendToList(lines)
	}
	// now write []string to buffer adding newlines
	for _, line := range lines {
		c.currBuff.WriteString(fmt.Sprintf("%s\n", line))
	}
	err = c.writeBufferToFile()
	return err
}

func (c *CredFile) fileExists() bool {
	_, err := os.Stat(c.filename)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func (c *CredFile) createFile() (bool, error) {
	_, err := os.Create(c.filename)
	if err != nil {
		return false, err
	}
	return true, err
}

// ProfileEntryInput holds properties required
// for acfmgr CredFile to write a profile entry
// to the file.
type ProfileEntryInput struct {
	Credential       *sts.Credentials   // MANDATORY: credentials object from aws.STS
	ProfileEntryName string             // MANDATORY: name of the desired profile entry e.g., '[devaccount]'. Brackets will be removed and spaces converted to dashes.
	Region           string             // OPTIONAL: region to include in the profile entry
	OutputFormat     string             // OPTIONAL: format for output when this credential is used, e.g., ('json', 'text')
	ExpiresToken     string             // OPTIONAL: a token so that string parsers can find the expiry date later
	InstanceRoleARN  string             // OPTIONAL: the ARN of the Instance Profile Role used to get these credentials
	AssumeRoleARN    string             // OPTIONAL: the ARN of the role that was assumed to get these credentials
    Description      string             // OPTIONAL: a description to give this entry
	TemplateOverride *template.Template // OPTIONAL: a text/template.Template to override the package default for this entry
}

type basicCredential struct {
	AccessKeyID     string // comes from sts.Credentials
	SecretAccessKey string // comes from sts.Credentials
	SessionToken    string // comes from sts.Credentials
	Expiration      string // comes from sts.Credentials
	Generated       string // comes from time.Now()
	OutputFormat    string
	Region          string
	ExpiresToken    string
	InstanceRoleARN string
	AssumeRoleARN   string
    Description     string
    HasRegion       bool
    HasOutput       bool
    HasDescription  bool
}

// dump returns a formatted json version of the
// basic credential for debugging purposes
func (bc *basicCredential) dump() (string, error){
    b, err := json.MarshalIndent(bc, "", "    ")
    return string(b), err
}

func (c *CredFile) NewEntry(pfi *ProfileEntryInput) (err error) {
	// clean up profile name input
	credName := strings.Replace(pfi.ProfileEntryName, " ", "-", -1)
	credName = strings.Replace(credName, "[", "", -1)
	credName = strings.Replace(credName, "]", "", -1)
	credName = fmt.Sprintf("[%s]", credName)
	if len(credName) < 1 {
		err = errors.New("ProfileEntryName cannot be blank")
		return err
	}
	// build basicCredential with defaults unless user specifies
	var bc basicCredential
	loc, _ := time.LoadLocation("UTC")
	bc.Generated = time.Now().In(loc).String()
	if pfi.Region != "" {
        bc.HasRegion = true
		bc.Region = pfi.Region
	}

	if pfi.ExpiresToken == "" {
		bc.ExpiresToken = "# EXPIRES@"
	} else {
		bc.ExpiresToken = pfi.ExpiresToken
	}

	if pfi.AssumeRoleARN == "" {
		bc.AssumeRoleARN = "NA"
	} else {
		bc.AssumeRoleARN = pfi.AssumeRoleARN
	}

	if pfi.InstanceRoleARN == "" {
		bc.InstanceRoleARN = "NA"
	} else {
		bc.InstanceRoleARN = pfi.InstanceRoleARN
	}

	if pfi.OutputFormat != "" {
        bc.HasOutput = true
		bc.OutputFormat = pfi.OutputFormat
	}
	if pfi.Description != "" {
        bc.HasDescription = true
		bc.Description = pfi.Description
	}
	// certain things always come from the sts.Credentials object
	bc.AccessKeyID = *pfi.Credential.AccessKeyId
	bc.SecretAccessKey = *pfi.Credential.SecretAccessKey
	bc.SessionToken = *pfi.Credential.SessionToken
	bc.Expiration = pfi.Credential.Expiration.String()
	buf := new(bytes.Buffer)
	if pfi.TemplateOverride != nil {
		err = pfi.TemplateOverride.Execute(buf, bc)
		if err != nil {
			return err
		}
	} else {
		// use package default
		err = defaultTemplate.Execute(buf, bc)
		if err != nil {
			return err
		}
	}
	credContents := strings.Split(buf.String(), "\n")
	c.addEntry(credName, credContents)
	return err
}

// expandPath takes a file path as a string and attempts to
// expand things like tildes, %userprofile%, $HOME etc. to form a full
// absolute path.
func expandPath(path string, usr *user.User) (expandedPath string, err error) {
        // relay to the OS specific function built at build time as determined by +build flags
        return expandPathO(path, usr)
}

