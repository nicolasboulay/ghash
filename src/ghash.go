package main

import (
	"fmt"
	"crypto/sha512"
	"bufio"
	"os"
	"math"
	"os/exec"
	"log"
	"strconv"
	"flag"
	"strings"
	"path/filepath"
	"path"
	"regexp"
)

func process(a byte, b byte, c byte, d byte) float64 {
	return	(float64(a)+256.0*float64(b)+65536.0*float64(c)+16777216.0*float64(d) )/(math.Pow(2,32));
}

func processS(slice []byte) float64 {
	//fmt.Printf("len: %d, cap: %d\n", len(slice), cap(slice))
	//fmt.Println(slice);
	return	process(slice[0],slice[1],slice[2],slice[3])
}

func toFloat64Slice(self []byte) []float64 {
	r := make([]float64,8,8);

	r[0] = processS(self[0:4]);
	r[1] = processS(self[4:8]);
	r[2] = processS(self[8:12]);
	r[3] = processS(self[12:16]);
	r[4] = processS(self[16:20]);
	r[5] = processS(self[20:24]);
	r[6] = processS(self[24:28]);
	r[7] = processS(self[28:32]);
	return r
}

func toFloat64Slice2(self []byte, r []float64) []float64 {
	r[0] = processS(self[0:4]);
	r[1] = processS(self[4:8]);
	r[2] = processS(self[8:12]);
	r[3] = processS(self[12:16]);
	r[4] = processS(self[16:20]);
	r[5] = processS(self[20:24]);
	r[6] = processS(self[24:28]);
	r[7] = processS(self[28:32]);
	return r
}

func  toS(f float64) string {
	return strconv.FormatFloat(f,'f',5,32)
}


func generateImage(size int, ft []float64, filename string, test bool ) {
	s := toS(ft[0]) 
	for _, f := range ft[1:] {
		s = s + "," + toS(f)
	}
	s = strconv.Itoa(size) + "," +s
	fmt.Printf(" %s\n", s)
	runGmic(s,filename, test) 
}
var executableFolder string = getExecutableFolder() 
func getExecutableFolder() string {
	path, _ := getExecutablePathOnLinux()
	goRun, _ := regexp.MatchString("/tmp/*", path) // if "go run" is used
	if (goRun) {
		return "."
	}
	
	return filepath.Dir(path)
}

var executablePath string = "" // caching
func getExecutablePathOnLinux() (string, error) {

	if(executablePath == "") {
		const deletedTag = " (deleted)"
		execpath, err := os.Readlink("/proc/self/exe")
		if err != nil {
			return execpath, err
		}
		execpath = strings.TrimSuffix(execpath, deletedTag)
		execpath = strings.TrimPrefix(execpath, deletedTag)
		executablePath = execpath
		return execpath, nil
	} else {
		return executablePath, nil
	}
}

var gmicScriptPath string = path.Join(executableFolder,"ghash.gmic") 
var gmicPath string = path.Join(executableFolder,"gmic") 
func runGmic(s string, filename string, test bool ) {

	cmd := exec.Command(gmicPath, gmicScriptPath, "-ghash", s, "-o[0]", filename)
	if(test) {
		cmd = exec.Command(gmicPath, gmicScriptPath, "-ghash", s, "-o[0]", filename, "-o[1]", "gradient_" + filename )
	}
	//

	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf(" %s\n", out)
		fmt.Printf("[ERROR] gmic have been tested using version 1.7.7 : \n\n")
		cmd := exec.Command("gmic", "-version") // ok 1.7.7
		out, err := cmd.CombinedOutput()
		fmt.Printf(" %s\n", out)
		log.Fatal(err)
	}
	fmt.Printf(" %s\n", out)
}


//func runPSNR(filenameA string, filenameB string ) {
//	cmd := exec.Command("gmic",filenameA, filenameB, "-psnr" )
//
//	out, err := cmd.CombinedOutput()
//	if err != nil {
//		fmt.Printf(" %s\n", out)
//	log.Fatal(err)
//	}
//	fmt.Printf(" %s\n", out)
//}

func ScanAndHash() []float64 {
	hash := sha512.New()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
	    hash.Write([]byte(scanner.Text()))
	}
	md := hash.Sum(nil) //1st round

	hashBis := sha512.New()
	hashBis.Write(md)
	md2 := hashBis.Sum(nil) //2nd round

	hashTer := sha512.New()
	hashTer.Write(md2)
	md3 := hashTer.Sum(nil )//3rd round

	ft := make([]float64,32,32)
	toFloat64Slice2(md,ft[0:8])
	toFloat64Slice2(md2,ft[8:16])
	toFloat64Slice2(md3,ft[16:24])
	return ft
}

func main() {

	var name = flag.String("o", "ghash.jpg", "name of the output file")
	var size = flag.Int("size", 128, "size of the square image")
	var test = flag.Bool("test", false, "generate many image with smallest variation")

	flag.Parse();

	ft := ScanAndHash()
	generateImage(*size, ft, *name, *test) 

	ftBis := make([]float64,32,32)

	if (*test){ 
		for i, _ := range ft[0:23] {
			copy(ftBis[:],ft[:]);
			ftBis[i] += 0.01
			generateImage(*size, ftBis, "test"+ strconv.Itoa(102+i)+ "_"  + *name, *test) 
		}
	}
}
