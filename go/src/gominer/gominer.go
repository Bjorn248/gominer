package main

import (
	"crypto/sha1"
	"fmt"
	"bytes"
	"log"
	"os"
	// "strconv"
	"io"
	"os/exec"
	"encoding/hex"
	"io/ioutil"
	// "strings"
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

func solve(start int, interval int, tree string, parent string, author string, committer string, difficulty string) {
	hasher := sha1.New()
	iterator := start
	for 1 == 1 {
		iterator += interval
		body := fmt.Sprintf("%s%s%s\n%s\nGive me a Gitcoin\n\n%d", tree, parent, author, committer, iterator)
		store := fmt.Sprintf("commit %d\\0%s", len(body), body)
		hasher.Reset()
		io.WriteString(hasher, store)
		digest := hex.EncodeToString(hasher.Sum(nil))
		if digest < difficulty {
			 fmt.Println(digest, iterator)
		}
	}
}

func main() {
	public_username := "user-dwj9pqp4"
	a, b := exists("/Users/bstange/StripeCTF/gominer/level1")
	if a == true {
		fmt.Println("level 1 exists")
		fmt.Println(b)
	} else {
		fmt.Println("can't find level 1, time to clone")
		outString := shellcmd("git", "clone", "lvl1-d6wr0qcx@stripe-ctf.com:level1", "/Users/bstange/StripeCTF/gominer/level1")
		fmt.Printf("Clone Output: %q\n", outString)
	}
	os.Chdir("/Users/bstange/StripeCTF/gominer/level1")
	fmt.Println(shellcmd("git", "reset", "--hard", "HEAD"))
	ledger, err := ioutil.ReadFile("LEDGER.txt")
	if err != nil { panic(err) }
	updatedledger := []byte(fmt.Sprintf("%s%s: 1", ledger, public_username))
	err = ioutil.WriteFile("LEDGER.txt", updatedledger, 0644)
	if err != nil { panic(err) }
	fmt.Println(shellcmd("git", "add", "LEDGER.txt"))
	tree := fmt.Sprintf("tree %s", shellcmd("git", "write-tree"))
	difficulty := shellcmd("cat", "difficulty.txt")
	parent := fmt.Sprintf("parent %s", shellcmd("git", "rev-parse", "HEAD"))
	timestamp := shellcmd("date", "+%s")
	author := fmt.Sprintf("author CTF user <me@example.com> %s +0000", timestamp)
	committer := fmt.Sprintf("committer CTF user <me@example.com> %s +0000", timestamp)

	go solve(0, 2, tree, parent, author, committer, difficulty)
	go solve(1, 2, tree, parent, author, committer, difficulty)
	go solve(-1, -2, tree, parent, author, committer, difficulty)
	go solve(0, -2, tree, parent, author, committer, difficulty)

	var input string
	fmt.Scanln(&input)
	fmt.Println("done")
}
