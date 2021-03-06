package main

import (
	"fmt"
	"crypto/sha512"
	"golang.org/x/crypto/sha3"
	"golang.org/x/crypto/bcrypt"
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
	"hash"
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


func generateImage(size int, ft []float64, filename string, test bool, verbose bool ) {
	s := toS(ft[0]) 
	for _, f := range ft[1:] {
		s = s + "," + toS(f)
	}
	s = strconv.Itoa(size) + "," +s
	if(verbose) {
		fmt.Printf("parameters : %s\n", s)
	}
	runGmic(s,filename, test, verbose) 
}

//  return the folder of the executable
//  return "." if find inside /tmp/ like for go run execution which is a tiny work around
var executableFolder string = getExecutableFolder() 
func getExecutableFolder() string {
	path, _ := getExecutablePathOnLinux()
	goRun, _ := regexp.MatchString("/tmp/*", path) // if "go run" is used
	if (goRun) {
		return "."
	}
	
	return filepath.Dir(path)
}

func getExecutablePathOnLinux() (string, error) {
	const deletedTag = " (deleted)"
	execpath, err := os.Readlink("/proc/self/exe")
	if err != nil {
		return execpath, err
	}
	execpath = strings.TrimSuffix(execpath, deletedTag)
	execpath = strings.TrimPrefix(execpath, deletedTag)
	return execpath, nil
}

func findGmic() string {
	p := path.Join(executableFolder,"gmic") 
	cmd := exec.Command(p, "-version") // try gmic beside binary 
	_, err := cmd.CombinedOutput()
	if(err != nil) {
		p = "gmic" // try gmic of the system
	}
	return p
}

var gmicScriptPath string = path.Join(executableFolder,"ghash.gmic") 
var gmicPath string = findGmic() 
func runGmic(s string, filename string, test bool, verbose bool ) {
	cmd := exec.Command(gmicPath, gmicScriptPath, "-ghash", s, "-o[0]", filename)
	if(verbose) {
		cmd = exec.Command(gmicPath, gmicScriptPath, "-ghash", s, "-o[0]", filename, "-o[1]", "gradient_" + filename )
	}
	//

	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf(" %s\n", out)
		cmd := exec.Command(gmicPath, "-version")
		out, _ := cmd.CombinedOutput()
		fmt.Printf(" %s\n", out)
		fmt.Printf("[ERROR] G'MIC have been tested using version 1.7.7.\n")
		log.Fatal(err)
	}
	if(verbose) {
		fmt.Printf(" %s\n", out)
	}
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

// hash is  already full of data 
// 3 consecutive hash for 24 float64 number
func hashToParameters(hash hash.Hash) []float64 {
	md := hash.Sum(nil) //1st round
	
	hash.Reset()
	hash.Write(md)
	md2 := hash.Sum(nil) //2nd round

	hash.Reset()
	hash.Write(md2)
	md3 := hash.Sum(nil )//3rd round

	ft := make([]float64,32,32)
	toFloat64Slice2(md,ft[0:8])
	toFloat64Slice2(md2,ft[8:16])
	toFloat64Slice2(md3,ft[16:24])
	return ft
}

func ScanAndHash2(strong bool) ([]float64, []float64) {
	hash := sha512.New()
	hash2 := sha3.New512()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		hash.Write(scanner.Bytes())
		if strong { 
			hash2.Write(scanner.Bytes())
		}
	}

	ft := hashToParameters(hash)
	var ft2 []float64 
	if strong { 
		md2 := hash.Sum(nil) // first sha3-512 hash 
		md3, _ := bcrypt.GenerateFromPassword(md2,12) //second bcrypt hash (work only on 50-70 bytes input)
		//fmt.Printf("%#v\n", err, len(md3), md3)
		hash2.Reset()  // if removed (?) this will mix 60 bytes bcrypt hash with the sha3-512 bytes hash
		hash2.Write(md3)
		ft2 = hashToParameters(hash2)
	}
	return ft,ft2
}

func main() {

	var name = flag.String("o", "ghash.jpg", "name of the output file")
	var name2 = flag.String("o2", "ghash2.jpg", "name of the 2nd output file (strong option)")
	var size = flag.Int("size", 128, "size of the square image")
	var test = flag.Bool("test", false, "generate many image with smallest variation")
	var verbose = flag.Bool("v", false, "more information")
	var strong = flag.Bool("9", false, "strong: generate a 2 images (one with sha2 and the other with sha3) for kind of cryptographic protection")
	flag.Parse();

	ft,ft2 := ScanAndHash2(*strong)
	generateImage(*size, ft, *name, *test, *verbose) 

	if * strong {

		generateImage(*size, ft2, *name2, *test, *verbose) 
	}
	

	if (*test){ 
		ftBis  := make([]float64,32,32)
		ftBis2 := make([]float64,32,32)

		for i, _ := range ft[0:23] {
			copy(ftBis[:],ft[:]);
			ftBis[i] += 0.01
			generateImage(*size, ftBis, "test"+ strconv.Itoa(102+i)+ "_" + *name, *test,*verbose)
			if * strong {
				copy(ftBis2[:],ft2[:]);
				ftBis2[i] += 0.01
				generateImage(*size, ftBis2, "test2"+ strconv.Itoa(102+i)+ "_" + *name2  , *test,*verbose) 
			}
		}
	}
}
