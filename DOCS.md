# acfmgr godocs
```
package acfmgr // import "github.com/GESkunkworks/acfmgr"

Package acfmgr is a package to manage entries in an AWS credentials file

Sample resulting AWS creds file entries:

    [account2]
    # DO NOT EDIT
    # ACFMGR MANAGED SECTION
    # (Will be overwritten regularly)
    ####################################################
    # ASSUMED ROLE: arn:aws:iam::098765432123:role/aj/d-readonly
    # ASSUMED FROM INSTANCE ROLE: NA
    # GENERATED: 2020-01-09 22:30:07.527022 +0000 UTC
    # EXPIRES@   2020-01-09 23:30:04 +0000 UTC
    # DESCRIPTION: gossamer-legacy
    output = json
    region = us-east-1
    aws_access_key_id = ASIASDIVWOEIOBINAIE
    aws_secret_access_key = GyD6rud3Q06qk90pTlECLncbbKx7GPXjM2N5ocVe
    aws_session_token = FwoGZXIvYXdzED0aDNHw4GhQvSFSCn8vUCK6Af+KK2QGsRbN5F22xJvXyNyYoAzxTkPYrSgvvuL7/17tyBa5LMeHWSKV/9E3ON2vRSLIz0iFfeEE5cj4zmbqpw/5LAiDiptTvbQQKmzCE4Pt05khFcsTmwsju9ibR5Mx2oJKdHHQXCsqk0XjvugSuu+KbU0wigO2oSXvu1dguNg+j6RTdxGAS7Uoih2WZR4ZlJCdcFNOivhf/kWs18mMRQ43r47GWsV9Z3vlTaMimHLWuBMldPgBcJV2iCiWrpnwBTIt2Dfkgvi8Bs7OcInotWE751K48QJnzcwPMKjsNKBE0tf1kGI9JArO8x+aDQJX

    [account1]
    # DO NOT EDIT
    # ACFMGR MANAGED SECTION
    # (Will be overwritten regularly)
    ####################################################
    # ASSUMED ROLE: arn:aws:iam::123456789012:role/aj/d-admin
    # ASSUMED FROM INSTANCE ROLE: NA
    # GENERATED: 2020-01-09 22:30:08.518222 +0000 UTC
    # EXPIRES@   2020-01-09 23:30:05 +0000 UTC
    # DESCRIPTION: gossamer-legacy
    output = json
    region = us-east-2
    aws_access_key_id = ASIZIPVKAVLEIGH
    aws_secret_access_key = GFIzAhjrBu9h3kkUNimntlN0zTF0TLmMebkMM0AP
    aws_session_token = FwoGZXIvYXdzED0aDNnD7rODZ+iSKGSi6yK6ASyAn9hJz0ER81eYa842hUwQnca7QNeFq7ZOrYvKb3ZegoVSRFEOaMvgw5La/taN8udMAdFINmxxV7Fx7JGTWuMK9JbWQA8I/AHqCN/NmOC1PrbIRvUhAZ8FgTdjNbiyh8CoOEvFqI3n4uQ57oWG5EZGZh8DSfENoVANR1AIaod7sFU1yHnHKOlr5Zp/iIUcD5j8X8yY8m05Vj2JFNipwcUsIVTNCeaWMud5n/30F4g/sQJLrIcV3nNoTCiWrpnwBTItwq++PEN9OzQYkCIEFNcvJe2ZnkiPz/+4xDNSSfiBGsHMKCGdChINizxQQHzE
    ...

Adding and removing entries manually is a pain when doing large amounts of
assume role operations so this package was created to assist in
programattically adding them once you have sessions built from the Golang
AWS SDK.

To start building a credentials file all you have to do is start a
NewCredFileSession with the desired filename. If the file does not exist
acfmgr will create it fore you. Use the CredFile object that's returned to
manage the entries in the file.

From there you use ProfileEntryInput objects and the CredFile.NewEntry() to
put new credentials in the queue. After you're done creating objects you can
write them all to file with CredFile.AssertEntries or delete all of the
queued profile names with CredFile.DeleteEntries.

The package will not touch any existing profile entries unless they match
the name of a profile entry that you provide to the CredFile object.

Sample with minimum input params

    c, err := acfmgr.NewCredFileSession("~/.aws/credentials")
    check(err)
    profileInput := acfmgr.ProfileEntryInput{
        Credential: <some sts.Credentials object>,
        ProfileEntryName: "devaccount1",
    }
    err = c.NewEntry(profileInput)
    check(err)
    err = c.AssertEntries()
    check(err)

Another Sample using more of the input parameters:

    c, err := acfmgr.NewCredFileSession("~/.aws/credentials")
    check(err)
    profileInput := acfmgr.ProfileEntryInput{
        Credential: <some aws.Credentials object>,
        ProfileEntryName: "devaccount",
        OutputFormat: "text",
        Region: "us-east-2",
        AssumeRoleARN: "arn:aws:iam::123456789012:role/aj/d-admin"
    }
    err = c.NewEntry(profileInput)
    check(err)
    err = c.DeleteEntries()

TYPES

type CredFile struct {
	// Has unexported fields.
}
    CredFile should be built with the exported NewCredFileSession function.

func NewCredFileSession(filename string) (cf *CredFile, err error)
    NewCredFileSession creates a new interactive credentials file session. Needs
    target filename and returns CredFile obj and err.

func (c *CredFile) AssertEntries() (err error)
    AssertEntries loops through all of the credEntry objs attached to CredFile
    obj and makes sure there is an occurrence with the credEntry.name and
    contents. Existing entries of the same name with different contents will be
    clobbered.

func (c *CredFile) DeleteEntries() (err error)
    DeleteEntries loops through all of the credEntry objs attached to CredFile
    obj and makes sure entries with the same credEntry.name are removed. Will
    remove ALL entries with the same name.

func (c *CredFile) NewEntry(pfi *ProfileEntryInput) (err error)

type ProfileEntryInput struct {
	Credential       *aws.Credentials   // MANDATORY: credentials object from aws.STS
	ProfileEntryName string             // MANDATORY: name of the desired profile entry e.g., '[devaccount]'. Brackets will be removed and spaces converted to dashes.
	Region           string             // OPTIONAL: region to include in the profile entry
	OutputFormat     string             // OPTIONAL: format for output when this credential is used, e.g., ('json', 'text')
	ExpiresToken     string             // OPTIONAL: a token so that string parsers can find the expiry date later
	InstanceRoleARN  string             // OPTIONAL: the ARN of the Instance Profile Role used to get these credentials
	AssumeRoleARN    string             // OPTIONAL: the ARN of the role that was assumed to get these credentials
	Description      string             // OPTIONAL: a description to give this entry
	TemplateOverride *template.Template // OPTIONAL: a text/template.Template to override the package default for this entry
}
    ProfileEntryInput holds properties required for acfmgr CredFile to write a
    profile entry to the file.
```
