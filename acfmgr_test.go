package acfmgr

import (
    "bytes"
    "io/ioutil"
    "os"
	"strings"
    "testing"
	"os/user"
	"time"
	"runtime"
	"github.com/aws/aws-sdk-go-v2/aws"
)

const baseCredFile string = `
[testing]
foo
bar

[newentry]
bar
foo

`

const expectedResult string = `
[testing]
foo
bar

[newentry]
bar
foo

[acfmgrtest]
# DO NOT EDIT
# ACFMGR MANAGED SECTION
# (Will be overwritten regularly)
####################################################
# ASSUMED ROLE: NA
# ASSUMED FROM INSTANCE ROLE: NA
# GENERATED: 2020-01-13 18:38:19.389637 +0000 UTC
# EXPIRES@   2020-01-08 14:03:02 +0000 UTC
aws_access_key_id = AHENVMSKIRUEQNFHGZTA
aws_secret_access_key = ZcqCQl34NF8PtXHSdbBk3mZze1plNNSWqnmsz523
aws_session_token = f8sNh8tocFpiabpbOGHfpqSYSgOQcNqvbzyNpAYW9gxWOlAcGpaPJMQoeDM/0AQjHnvA8qMA8Q2jdxFmPwLHA184JI9YXVXs3a6ig2GMKvtTYXYwe4HKbymJm4zWxcG7OWwPee8BlZbY+F/T+lmNguge42ePV3mA5uyK5oTgryTG9TNFBtmh518OCdRXBDwwPWwQbfLWM/95KaOnZRIr/TpkjdWk4iCFXmKTIs5RKwDrS9mmD66cj6KTNsAGDxw29wYLOXlcB3MXbuEZzgew6tn8vpzonBIRiFy74Oym6Ct1sFcNXVKrwmn2Ojnmec3KCAbFwynyTHPxE2PpHlVhQhvb2Azw2FeLGAw1btiItcvLDrS3cDI3TfQNaa8L2MX3Zfr2yBv9UUS4MfS2pZQ42Czze7PMRk6LrWh0HA+SdBUG6XeXDHcvXH3rH4GxJHuDhALCgNabFYwuXysXdGP=

`

const expectedResultDeletionEnd string = `
[testing]
foo
bar

`

const expectedResultDeletionMiddle string = `
[newentry]
bar
foo

`

func writeBaseFile(filename string) error {
    var b bytes.Buffer
    _, err := b.WriteString(baseCredFile)
    err = ioutil.WriteFile(filename, b.Bytes(), 0644)
    return err
}

func replaceLine(contents []byte, linecontains, replaceline string) (output []byte) {
	lines := strings.Split(string(contents), "\n")
	for i, line := range lines {
		if strings.Contains(line, linecontains) {
				lines[i] = replaceline
		}
	}
	newcontents := strings.Join(lines, "\n")
	output = []byte(newcontents)
	return output
}

func TestModifyEntry(t *testing.T) {
    filename := "./acfmgr_credfile_test.txt"
    err := writeBaseFile(filename)
    if err != nil {
        t.Errorf("Error making basefile: %s", err)
    }
    sess, err := NewCredFileSession(filename)
    if err != nil {
        t.Errorf("Error making credfile session: %s", err)
    }
	pfi := ProfileEntryInput{
		Credential: getFakeCreds(),
		ProfileEntryName: "acfmgrtest",
	}
    err = sess.NewEntry(&pfi)
	if err != nil {
		t.Errorf("Error adding entry: %s", err)
	}
    err = sess.AssertEntries()
    if err != nil {
        t.Errorf("Error asserting entries: %s", err)
    }
    fullContents, err := ioutil.ReadFile(filename)
    if err != nil {
        t.Errorf("Error reading file: %s", err)
    }
	// fix generated date so contents match
	replacement := "# GENERATED: 2020-01-13 18:38:19.389637 +0000 UTC"
	fullContents = replaceLine(fullContents, "GENERATED", replacement)
    got := string(fullContents)
    if got != expectedResult {
        t.Errorf("Result not expected. Got: %s", got)
    }
    defer os.Remove(filename)
}

func TestDeleteEntryAtEnd(t *testing.T) {
    filename := "./acfmgr_credfile_test1.txt"
    err := writeBaseFile(filename)
    if err != nil {
        t.Errorf("Error making basefile: %s", err)
    }
    sess, err := NewCredFileSession(filename)
    if err != nil {
        t.Errorf("Error making credfile session: %s", err)
    }
	pfi := ProfileEntryInput{
		Credential: getFakeCreds(),
		ProfileEntryName: "newentry",
	}
    err = sess.NewEntry(&pfi)
	if err != nil {
		t.Errorf("Error adding entry: %s", err)
	}
    err = sess.DeleteEntries()
    if err != nil {
        t.Errorf("Error deleting entries: %s", err)
    }
    fullContents, err := ioutil.ReadFile(filename)
    if err != nil {
        t.Errorf("Error reading file: %s", err)
    }
	// fix generated date so contents match
	replacement := "# GENERATED: 2020-01-13 18:38:19.389637 +0000 UTC"
	fullContents = replaceLine(fullContents, "GENERATED", replacement)
    got := string(fullContents)
    if got != expectedResultDeletionEnd {
        t.Errorf("Result not expected. Got: %s", got)
    }
    defer os.Remove(filename)
}

func TestDeleteEntryInMiddle(t *testing.T) {
    filename := "./acfmgr_credfile_test2.txt"
    err := writeBaseFile(filename)
    if err != nil {
        t.Errorf("Error making basefile: %s", err)
    }
    sess, err := NewCredFileSession(filename)
    if err != nil {
        t.Errorf("Error making credfile session: %s", err)
    }
	pfi := ProfileEntryInput{
		Credential: getFakeCreds(),
		ProfileEntryName: "testing",
	}
    err = sess.NewEntry(&pfi)
	if err != nil {
		t.Errorf("Error adding entry: %s", err)
	}
    err = sess.DeleteEntries()
    if err != nil {
        t.Errorf("Error deleting entries: %s", err)
    }
    fullContents, err := ioutil.ReadFile(filename)
    if err != nil {
        t.Errorf("Error reading file: %s", err)
    }
    got := string(fullContents)
    if got != expectedResultDeletionMiddle {
        t.Errorf("Result not expected. Got: %s", got)
    }
    defer os.Remove(filename)
}

func gimmeUserWindows() (fakeUser *user.User) {
        u := user.User{
                Uid:      "S-1-5-21-693214013-1772980081-1954060963-500",
                Gid:      "S-1-5-21-693214013-1772980081-1954060963-513",
                Username: "EC2AMAZ-K9GK6S2\\Administrator",
                Name:     "",
                HomeDir:  "C:\\Users\\Administrator",
        }
        fakeUser = &u
        return fakeUser
}

func gimmeUserLinux() (fakeUser *user.User) {
        u := user.User{
                Uid:      "1001",
                Gid:      "1001",
                Username: "admin",
                Name:     "",
                HomeDir:  "/home/admin",
        }
        fakeUser = &u
        return fakeUser
}

func gimmeUserMac() (fakeUser *user.User) {
        u := user.User{
                Uid:      "501",
                Gid:      "20",
                Username: "dudedudem",
                Name:     "Dude Dudem",
                HomeDir:  "/Users/dudedudem",
        }
        fakeUser = &u
        return fakeUser
}

func TestExpandPath(t *testing.T) {
        cases := []struct {
                Path                  string
                Want                  string
                User                  *user.User
                ExpectedErr           error
                ExpectedRuntime       string
                PrefixWantCurrentUser bool
        }{
                {
                        User:            gimmeUserMac(),
                        Path:            "~/.aws/credentials",
                        Want:            "/Users/dudedudem/.aws/credentials",
                        ExpectedRuntime: "darwin",
                },
                {
                        User:            gimmeUserMac(),
                        Path:            "/Users/dudedudem/.aws/credentials",
                        Want:            "/Users/dudedudem/.aws/credentials",
                        ExpectedRuntime: "darwin",
                },
                // test will always fail because os.ExpandEnv always takes runtime user, ugh
                // {
                //      User: gimmeUserMac(),
                //  Path: "$HOME/.aws/credentials",
                //      ExpectedRuntime: "darwin",
                // },
                {
                        User:            gimmeUserWindows(),
                        Path:            "~\\.aws\\credentials",
                        Want:            "C:\\Users\\Administrator\\.aws\\credentials",
                        ExpectedRuntime: "windows",
                },
                // test will probably fail unless test system user is exactly "Administrator"
                // {
                //      User:            gimmeUserWindows(),
                //      Path:            "%userprofile%\\.aws\\credentials",
                //      Want:            "C:\\Users\\Administrator\\.aws\\credentials",
                //      ExpectedRuntime: "windows",
                // },
                {
                        User:            gimmeUserLinux(),
                        Path:            "~/.aws/credentials",
                        Want:            "/home/admin/.aws/credentials",
                        ExpectedRuntime: "linux",
                },
        }
        for _, c := range cases {
                if runtime.GOOS == c.ExpectedRuntime { // otherwise tests will bomb
                        have, err := expandPath(c.Path, c.User)
                        if err != nil {
                                if c.ExpectedErr == nil {
                                        t.Errorf("Unexpected error: %s\n", err.Error())
                                }
                        } else if have != c.Want {
                                        t.Errorf("Unexpected result. Have: '%s', Want: '%s'\n", have, c.Want)
                        } else {
                                t.Logf("PASS: Have '%s', Want '%s'\n", have, c.Want)
                        }

                }
        }
}

func getFakeCreds() *aws.Credentials {
    timeFormat := "2006-01-02T15:04:05Z" // time.RFC3339
    t, _ := time.Parse(timeFormat, "2020-01-08T14:03:02Z")
    st := "f8sNh8tocFpiabpbOGHfpqSYSgOQcNqvbzyNpAYW9gxWOlAcGpaPJMQoeD" +
        "M/0AQjHnvA8qMA8Q2jdxFmPwLHA184JI9YXVXs3a6ig2GMKvtTYXYwe4HK" +
        "bymJm4zWxcG7OWwPee8BlZbY+F/T+lmNguge42ePV3mA5uyK5oTgryTG9" +
        "TNFBtmh518OCdRXBDwwPWwQbfLWM/95KaOnZRIr/TpkjdWk4iCFXmKTIs5" +
        "RKwDrS9mmD66cj6KTNsAGDxw29wYLOXlcB3MXbuEZzgew6tn8vpzonBIRi" +
        "Fy74Oym6Ct1sFcNXVKrwmn2Ojnmec3KCAbFwynyTHPxE2PpHlVhQhvb2Az" +
        "w2FeLGAw1btiItcvLDrS3cDI3TfQNaa8L2MX3Zfr2yBv9UUS4MfS2pZQ42" +
        "Czze7PMRk6LrWh0HA+SdBUG6XeXDHcvXH3rH4GxJHuDhALCgNabFYwuXysXdGP="
    cred := aws.Credentials{
        AccessKeyID:     []string{"AHENVMSKIRUEQNFHGZTA"}[0],
        SecretAccessKey: []string{"ZcqCQl34NF8PtXHSdbBk3mZze1plNNSWqnmsz523"}[0],
        SessionToken:    st,
        Expires:      t,
    }
    return &cred
}

