package Utils

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	crand "math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const garblePackage string = "mvdan.cc/garble@latest"

func Version() {
	Version := runtime.Version()
	Version = strings.Replace(Version, "go1.", "", -1)
	VerNumb, _ := strconv.ParseFloat(Version, 64)
	if VerNumb >= 19.1 {
	} else {
		log.Fatal("Error: The version of Go is to old, please update to version 1.19.1 or later")
	}
}

func CheckGarble() {
	bin, _ := exec.LookPath("env")
	var cmd *exec.Cmd
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	garble := "garble"
	if runtime.GOOS == "windows" {
		garble = garble + ".exe"
	}

	if _, err := os.Stat(filepath.Join(cwd, ".lib", garble)); err == nil {
		fmt.Println("[+] Garble is present")
	} else {
		fmt.Println("[!] Missing Garble... Downloading it now")

		switch runtime.GOOS {
		case "windows":
			pre_code := `
$env:GOBINB=$GOBIN;
$env:GOBIN="%s";

%s

$env:GOBIN=$GOBINB;
$env:GOBINB=$null
			`
			cmd_code := fmt.Sprintf("go install %s", garblePackage)
			code := fmt.Sprintf(pre_code, filepath.Join(cwd, ".lib"), cmd_code)
			fmt.Printf("[+] Executed code:\n%s\n", code)

			opt := strings.Join([]string{"-NonInteractive"}, " ")
			cmd = exec.Command("powershell.exe", opt, code)
		default:
			cmd = exec.Command(bin, "GOBIN="+filepath.Join(cwd, ".lib"), "go", "install", garblePackage)
		}

		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("%s: %s\n", err, stderr.String())
		}
		fmt.Println(out.String(), stderr.String())
	}
}

func Sha256(input string) {
	f, err := os.Open(input)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("[!] Sha256 hash of "+input+": %x\n", h.Sum(nil))
}

func Writefile(outFile string, result string) {
	cf, err := os.OpenFile(outFile, os.O_CREATE|os.O_WRONLY, 0644)
	check(err)
	defer cf.Close()
	_, err = cf.Write([]byte(result))
	check(err)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

const capletters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const hexchar = "abcef12345678890"

var (
	ErrInvalidBlockSize = errors.New("[-] Invalid Blocksize")

	ErrInvalidPKCS7Data = errors.New("[-] Invalid PKCS7 Data (Empty or Not Padded)")

	ErrInvalidPKCS7Padding = errors.New("[-] Invalid Padding on Input")
)

func Pkcs7Pad(b []byte, blocksize int) ([]byte, error) {
	if blocksize <= 0 {
		return nil, ErrInvalidBlockSize
	}
	if b == nil || len(b) == 0 {
		return nil, ErrInvalidPKCS7Data
	}
	n := blocksize - (len(b) % blocksize)
	pb := make([]byte, len(b)+n)
	copy(pb, b)
	copy(pb[len(b):], bytes.Repeat([]byte{byte(n)}, n))
	return pb, nil
}

func RandomBuffer(size int) []byte {
	buffer := make([]byte, size)
	_, err := rand.Read(buffer)
	if err != nil {
		log.Fatal(err)
	}
	return buffer
}

var randomGen *crand.Rand = crand.New(crand.NewSource(time.Now().UnixNano()))

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[randomGen.Intn(len(letters))]

	}
	return string(b)
}

func VarNumberLength(min, max int) string {
	var r string
	num := randomGen.Intn(max-min) + min
	n := num
	r = RandStringBytes(n)
	return r
}

func printHexOutput(input ...[]byte) {
	for _, i := range input {
		fmt.Println(hex.EncodeToString(i))
	}
}

func GenerateNumer(min, max int) int {
	num := randomGen.Intn(max-min) + min
	n := num
	return n

}

func CapLetter() string {
	n := 1
	b := make([]byte, n)
	for i := range b {
		b[i] = capletters[randomGen.Intn(len(capletters))]

	}
	return string(b)
}
