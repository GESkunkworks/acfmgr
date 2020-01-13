# acfmgr
Package acfmgr is a package to manage entries in an [AWS credentials file](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html).

Sample resulting AWS creds file entries:

```
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
 ```

# Justification
Adding and removing entries manually is a pain when doing large amounts of assume
role operations so this package was created to assist in programattically adding
them once you have sessions built from the Golang AWS SDK.

To start building a credentials file all you have to do is start a `NewCredFileSession`
with the desired filename. If the file does not exist acfmgr will create it fore you.
Use the CredFile object that's returned to manage the entries in the file. 

    NOTE: If the filename string has "phrases of intent" such as `~`, `$HOME`, `%USERPROFILE%` the package will try to expand them. 

From there you use `ProfileEntryInput` objects and the `CredFile.NewEntry()` to put new
credentials in the queue. After you're done creating objects you can write them all
to file with `CredFile.AssertEntries()` or delete all of the queued profile names with
`CredFile.DeleteEntries()`.

The package will not touch any existing profile entries unless they match the name 
of a profile entry that you provide to the CredFile object. 


Sample with minimum input params

```
import (
	"github.com/GESkunkworks/acfmgr"
)

func main() {
	c, _ := acfmgr.NewCredFileSession("~/.aws/credentials")
	profileInput := acfmgr.ProfileEntryInput{
		Credential: <some sts.Credentials object>,
		ProfileEntryName: "devaccount1",
	}
	_ = c.NewEntry(profileInput)
	_ = c.AssertEntries()
}
```

Another Sample using more of the input parameters:

```
import (
	"github.com/GESkunkworks/acfmgr"
)

func main() {
	c, _ := acfmgr.NewCredFileSession("~/.aws/credentials")
	profileInput := acfmgr.ProfileEntryInput{
		Credential: <some sts.Credentials object>,
		ProfileEntryName: "devaccount",
		OutputFormat: "text",
		Region: "us-east-2",
		AssumeRoleARN: "arn:aws:iam::123456789012:role/aj/d-admin"
	}
	_ = c.NewEntry(profileInput)
	_ = c.AssertEntries()
}
```

