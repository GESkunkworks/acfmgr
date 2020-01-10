package acfmgr

import (
    "bytes"
    "io/ioutil"
    "os"
    "testing"
	"os/user"
	"runtime"
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
my
test
here
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
    entryName := "[acfmgrtest]"
    entryContents := []string{"my", "test", "here"}
    sess.NewEntry(entryName, entryContents)
    err = sess.AssertEntries()
    if err != nil {
        t.Errorf("Error asserting entries: %s", err)
    }
    fullContents, err := ioutil.ReadFile(filename)
    if err != nil {
        t.Errorf("Error reading file: %s", err)
    }
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
    entryName := "[newentry]"
    entryContents := []string{"whocares"}
    sess.NewEntry(entryName, entryContents)
    err = sess.DeleteEntries()
    if err != nil {
        t.Errorf("Error deleting entries: %s", err)
    }
    fullContents, err := ioutil.ReadFile(filename)
    if err != nil {
        t.Errorf("Error reading file: %s", err)
    }
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
    entryName := "[testing]"
    entryContents := []string{"whocares"}
    sess.NewEntry(entryName, entryContents)
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

