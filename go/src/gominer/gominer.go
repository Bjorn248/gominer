package main

import (
	"crypto/sha1"
	// "runtime"
	"fmt"
	"bytes"
	"log"
	"os"
	"io"
	"os/exec"
	"encoding/hex"
	"io/ioutil"
	"time"
)

func shellcmd(name string, arg ...string) string {
	var cmdout bytes.Buffer
	cmd := exec.Command(name, arg...)
	cmd.Stdout = &cmdout
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	return cmdout.String()
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil { return true, nil }
	if os.IsNotExist(err) { return false, nil }
	return false, err
}

func solve(digestCount chan int, tree string, parent string, author string, committer string, difficulty string, ident int) {
	hasher := sha1.New()
	current := 0
	for {
		current++
		body := fmt.Sprintf("%s%s%s\n%s\nGive me a Gitcoin\n\n%d-%d", tree, parent, author, committer, ident, current)
		store := fmt.Sprintf("commit %d\\0%s", len(body), body)
		hasher.Reset()
		io.WriteString(hasher, store)
		digest := hex.EncodeToString(hasher.Sum(nil))
		if digest < difficulty {
			fmt.Println(digest, current)
		}
        digestCount <- 1
	}
}

func main() {
	public_username := "user-dwj9pqp4"
	leve1path = "/Users/bstange/StripeCTF/gominer/level1"
	a, b := exists(leve1path)
	if a == true {
		fmt.Println("level 1 exists")
		fmt.Println(b)
	} else {
		fmt.Println("can't find level 1, time to clone")
		// outString := shellcmd("git", "clone", "lvl1-d6wr0qcx@stripe-ctf.com:level1", leve1path)
		// fmt.Printf("Clone Output: %q\n", outString)
	}
	os.Chdir("/home/ajcrites/projects/personal/gominer/level1")
	// fmt.Println(shellcmd("git", "reset", "--hard", "HEAD"))
	ledger, err := ioutil.ReadFile("LEDGER.txt")
	if err != nil { panic(err) }
	updatedledger := []byte(fmt.Sprintf("%s%s: 1", ledger, public_username))
	err = ioutil.WriteFile("LEDGER.txt", updatedledger, 0644)
	if err != nil { panic(err) }
	// fmt.Println(shellcmd("git", "add", "LEDGER.txt"))
	tree := fmt.Sprintf("tree %s", shellcmd("git", "write-tree"))
	difficulty := shellcmd("cat", "difficulty.txt")
	parent := fmt.Sprintf("parent %s", shellcmd("git", "rev-parse", "HEAD"))
	timestamp := shellcmd("date", "+%s")
	author := fmt.Sprintf("author CTF user <me@example.com> %s +0000", timestamp)
	committer := fmt.Sprintf("committer CTF user <me@example.com> %s +0000", timestamp)
	var firstTime = time.Now()
    count := 0

	for x := 0; x < 4; x++ {
        digestCount := make(chan int)
		go solve(digestCount, tree, parent, author, committer, difficulty, x)
        go func (digestCount chan int) {
            for {
                <-digestCount
                count++
                //fmt.Println("Hash rate: ", int64(count) / int64(time.Since(firstTime)))
                if count == 1000000 {
                    fmt.Println(time.Since(firstTime))
                }
            }
        }(digestCount)
	}
    go func() {
        for {
            time.Sleep(1000 * 1000 * 1000)
            fmt.Println("Hash rate: ", int(count) / (int(time.Since(firstTime)) / 1000 / 1000 / 1000))
        }
    }()

	var input string
	fmt.Scanln(&input)
	fmt.Println("done")
}
