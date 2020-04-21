package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const oneKilo = 1024
const oneMega = 1048576
const oneGiga = 1073741824
const oneTera = 1099511627776

func checkFile(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	// получить файл
	dat, err := ioutil.ReadFile("subdir/some-file-1.txt")
	checkFile(err)
	fmt.Print(string(dat))

	fmt.Println("\n..................................")

	// открыть файл
	f, err := os.Open("subdir/some-file-1.txt")
	checkFile(err)

	// создали slice, map, or chan
	// прочитаем 5 байт из файла subdir/some-file-1.txt
	b1 := make([]byte, 5)
	n1, err := f.Read(b1)
	checkFile(err)

	// выводим подсчитанное количество байтов
	fmt.Printf("%d bytes: %s\n", n1, string(b1[:n1]))

	o2, err := f.Seek(6, 0)
	checkFile(err)
	b2 := make([]byte, 2)
	n2, err := f.Read(b2)
	checkFile(err)
	fmt.Printf("%d bytes @ %d: ", n2, o2)
	fmt.Printf("%v\n", string(b2[:n2]))

	o3, err := f.Seek(6, 0)
	checkFile(err)
	b3 := make([]byte, 2)
	n3, err := io.ReadAtLeast(f, b3, 2)
	checkFile(err)
	fmt.Printf("%d bytes @ %d: %s\n", n3, o3, string(b3))

	_, err = f.Seek(0, 0)
	checkFile(err)

	r4 := bufio.NewReader(f)
	b4, err := r4.Peek(5)
	checkFile(err)
	fmt.Printf("5 bytes: %s\n", string(b4))

	f.Close()

	// Write to File
	writeFile()

	fmt.Println("\n..................................\n")

	fmt.Println("--- Read file 'some-file-1.txt' by row ---")
	readFileByRow("subdir/some-file-1.txt")
	fmt.Println("--- Read file 'some-file-2.txt' by row ---")
	readFileByRow("subdir/some-file-2.txt")

	fmt.Println("\n..................................\n")

	FilePaths()

	fmt.Println("\n..................................\n")

	TemporaryFilesDirectories()

}

func readFileByRow(path string) {
	// открыли файл
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// создали сканер буферизированого ввода/вывода
	scanner := bufio.NewScanner(file)

	// читаем файл в цикле построчно
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func TemporaryFilesDirectories() {
	// The easiest way to create a temporary file is by
	// calling `ioutil.TempFile`. It creates a file *and*
	// opens it for reading and writing. We provide `""`
	// as the first argument, so `ioutil.TempFile` will
	// create the file in the default location for our OS.
	f, err := ioutil.TempFile("", "sample")
	checkFile(err)

	// Display the name of the temporary file. On
	// Unix-based OSes the directory will likely be `/tmp`.
	// The file name starts with the prefix given as the
	// second argument to `ioutil.TempFile` and the rest
	// is chosen automatically to ensure that concurrent
	// calls will always create different file names.
	fmt.Println("Temp file name:", f.Name())

	// Clean up the file after we're done. The OS is
	// likely to clean up temporary files by itself after
	// some time, but it's good practice to do this
	// explicitly.
	defer os.Remove(f.Name())

	// We can write some data to the file.
	_, err = f.Write([]byte{1, 2, 3, 4})
	checkFile(err)

	// If we intend to write many temporary files, we may
	// prefer to create a temporary *directory*.
	// `ioutil.TempDir`'s arguments are the same as
	// `TempFile`'s, but it returns a directory *name*
	// rather than an open file.
	dname, err := ioutil.TempDir("", "sampledir")
	checkFile(err)
	fmt.Println("Temp dir name:", dname)

	defer os.RemoveAll(dname)

	// Now we can synthesize temporary file names by
	// prefixing them with our temporary directory.
	fname := filepath.Join(dname, "file1")
	err = ioutil.WriteFile(fname, []byte{1, 2}, 0666)
	checkFile(err)
}

func FilePaths() {
	// `Join` should be used to construct paths in a
	// portable way. It takes any number of arguments
	// and constructs a hierarchical path from them.
	p := filepath.Join("dir1", "dir2", "filename")
	fmt.Println("p:", p)

	// You should always use `Join` instead of
	// concatenating `/`s or `\`s manually. In addition
	// to providing portability, `Join` will also
	// normalize paths by removing superfluous separators
	// and directory changes.
	fmt.Println(filepath.Join("dir1//", "filename"))
	fmt.Println(filepath.Join("dir1/../dir1", "filename"))

	// `Dir` and `Base` can be used to split a path to the
	// directory and the file. Alternatively, `Split` will
	// return both in the same call.
	fmt.Println("Dir(p):", filepath.Dir(p))
	fmt.Println("Base(p):", filepath.Base(p))

	// We can check whether a path is absolute.
	fmt.Println(filepath.IsAbs("dir/file"))
	fmt.Println(filepath.IsAbs("/dir/file"))

	filename := "config.json"

	// Some file names have extensions following a dot. We
	// can split the extension out of such names with `Ext`.
	ext := filepath.Ext(filename)
	fmt.Println(ext)

	// To find the file's name with the extension removed,
	// use `strings.TrimSuffix`.
	fmt.Println(strings.TrimSuffix(filename, ext))

	// `Rel` finds a relative path between a *base* and a
	// *target*. It returns an error if the target cannot
	// be made relative to base.
	rel, err := filepath.Rel("a/b", "a/b/t/file")
	if err != nil {
		panic(err)
	}
	fmt.Println(rel)

	rel, err = filepath.Rel("a/b", "a/c/t/file")
	if err != nil {
		panic(err)
	}
	fmt.Println(rel)
}

func writeFile() {
	// To start, here's how to dump a string (or just
	// bytes) into a file.
	d1 := []byte("hello\ngo\n")
	err := ioutil.WriteFile("subdir/some-file-1.txt", d1, 0644)
	checkFile(err)

	// For more granular writes, open a file for writing.
	f, err := os.Create("subdir/some-file-2.txt")
	checkFile(err)

	// It's idiomatic to defer a `Close` immediately
	// after opening a file.
	defer f.Close()

	// You can `Write` byte slices as you'd expect.
	d2 := []byte{115, 111, 109, 101, 10}
	n2, err := f.Write(d2)
	checkFile(err)
	fmt.Printf("wrote %d bytes\n", n2)

	// A `WriteString` is also available.
	n3, err := f.WriteString("writes\n")
	checkFile(err)
	fmt.Printf("wrote %d bytes\n", n3)

	// Issue a `Sync` to flush writes to stable storage.
	f.Sync()

	// `bufio` provides buffered writers in addition
	// to the buffered readers we saw earlier.
	w := bufio.NewWriter(f)
	n4, err := w.WriteString("buffered\n")
	checkFile(err)
	fmt.Printf("wrote %d bytes\n", n4)

	// Use `Flush` to ensure all buffered operations have
	// been applied to the underlying writer.
	w.Flush()
}
