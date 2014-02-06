package main

import (
	"crypto/sha1"
	"fmt"
	"bytes"
	"log"
	"os"
	"io"
	"os/exec"
	"encoding/hex"
	"io/ioutil"
	"strconv"
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

func solve(counter *int, tree string, parent string, author string, committer string, difficulty string) {
	hasher := sha1.New()
	for 1 == 1 {
		(*counter)++
		body := fmt.Sprintf("%s%s%s\n%s\nGive me a Gitcoin\n\n%d", tree, parent, author, committer, *counter)
		store := fmt.Sprintf("commit %d\\0%s", len(body), body)
		hasher.Reset()
		io.WriteString(hasher, store)
		digest := hex.EncodeToString(hasher.Sum(nil))
		if digest < difficulty {
			 fmt.Println(digest, *counter)
		}
	}
}

func main() {
	public_username := "user-dwj9pqp4"
	a, b := exists("/home/ajcrites/projects/personal/gominer/level1")
	if a == true {
		fmt.Println("level 1 exists")
		fmt.Println(b)
	} else {
		fmt.Println("can't find level 1, time to clone")
		// outString := shellcmd("git", "clone", "lvl1-d6wr0qcx@stripe-ctf.com:level1", "/Users/bstange/StripeCTF/gominer/level1")
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
	counter := 0

	maxprocs, err := strconv.Atoi(os.Getenv("GOMAXPROCS"))
	for x := 0; x < maxprocs; x++ {
		go solve(&counter, tree, parent, author, committer, difficulty)
	}

	var input string
	fmt.Scanln(&input)
	fmt.Println("done")
}
