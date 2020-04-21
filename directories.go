package main

import (
	"fmt"
	"io/ioutil"
	_ "log"
	"os"
	"path/filepath"
)

func checkSubdir() {

	// первый возвращенный результат нам не нужен
	// поэтому вместо него подчеркивание - пустота
	_, er := os.Stat("subdir")

	// если ошибки нет - значит есть директория
	// значит её надо удалить
	if er == nil {
		os.RemoveAll("subdir")
	}
}

func getFullPath(path string) {
	// выводим полный путь к месту где была создана временна директория
	dir, error := filepath.Abs(path)
	check(error)
	fmt.Println("Created at: " + dir)
}

func main() {

	// проверяем сабдиректорию
	checkSubdir()

	// удаляем временную директорию в самом конце
	//defer os.RemoveAll("subdir")

	// анонимная функция: объявили и тут же вызвали
	func(message string) {
		dir, error := filepath.Abs(filepath.Dir(os.Args[0]))
		check(error)
		fmt.Println(message + dir + "\n")
	}("Actual folder: ")
	/**
	актуальной директорией будет та, в которой будет находиться скомпилированный файл
	Actual folder: C:\Users\79161\go\src\example\_output
	*/

	// директория будет создана там, где находится этот файл
	err := os.Mkdir("subdir", 0755)
	check(err)

	// выводим полный путь к месту где была создана временна директория
	getFullPath("subdir")

	// замыкание: записали анонимную функцию в переменную
	createEmptyFile := func(name string) {
		d := []byte("")
		check(ioutil.WriteFile(name, d, 0644))
	}

	createEmptyFile("subdir/file1")

	err = os.MkdirAll("subdir/parent/child", 0755)
	check(err)

	createEmptyFile("subdir/parent/file2")
	createEmptyFile("subdir/parent/file3")
	createEmptyFile("subdir/parent/child/file4")

	info, err := os.Stat("subdir/parent/child/file4")
	fmt.Println("info.Mode()") // -rw-rw-rw-
	fmt.Println(info.Mode())
	fmt.Println("info.ModTime()") // 2020-04-21 09:06:08.104909 +0300 MSK
	fmt.Println(info.ModTime())
	fmt.Println("info.Size()")
	fmt.Println(info.Size())

	/**
	// Exactly one of O_RDONLY, O_WRONLY, or O_RDWR must be specified.
	O_RDONLY int = syscall.O_RDONLY // open the file read-only.
	O_WRONLY int = syscall.O_WRONLY // open the file write-only.
	O_RDWR   int = syscall.O_RDWR   // open the file read-write.
	// The remaining values may be or'ed in to control behavior.
	O_APPEND int = syscall.O_APPEND // append data to the file when writing.
	O_CREATE int = syscall.O_CREAT  // create a new file if none exists.
	O_EXCL   int = syscall.O_EXCL   // used with O_CREATE, file must not exist.
	O_SYNC   int = syscall.O_SYNC   // open for synchronous I/O.
	O_TRUNC  int = syscall.O_TRUNC  // truncate regular writable file when opened.
	*/
	_, e := os.OpenFile("subdir/parent/child/file4", os.O_WRONLY, 0666)
	if os.IsPermission(e) {
		fmt.Println("Unable to write to ", "subdir/parent/child/file4")
		fmt.Println(err)
		//os.Exit(1)
	}

	/**
	644 (-rw-r--r--)
	Все пользователи имеют право чтения; владелец может редактировать
	660 (-rw-rw----)
	Владелец и группа могут читать и редактировать; остальные не имеют права выполнять никаких действий
	664 (-rw-rw-r--)
	Все пользователи имеют право чтения; владелец и группа могут редактировать
	666 (-rw-rw-rw-)
	Все пользователи могут читать и редактировать
	700 (-rwx------)
	Владелец может читать, записывать и запускать на выполнение; никто другой не имеет права выполнять никакие действия
	744 (-rwxr--r--)
	Каждый пользователь может читать, владелец имеет право редактировать и запускать на выполнение
	755 (-rwxr-xr-x)
	Каждый пользователь имеет право читать и запускать на выполнение; владелец может редактировать
	777 (-rwxrwxrwx)
	*/

	// получили все файлы из директории
	c, err := ioutil.ReadDir("subdir/parent")
	check(err)

	fmt.Println("\nListing subdir/parent")

	/**
	range - это форма цикла for
	её используют для итерации среза или карты
	*/
	for _, entry := range c {
		fmt.Println(" ", entry.Name(), entry.IsDir())
	}

	// перешли из subdir в директорию subdir/parent/child
	err = os.Chdir("subdir/parent/child")
	check(err)

	c, err = ioutil.ReadDir(".")
	check(err)

	fmt.Println("Listing subdir/parent/child")
	for _, entry := range c {
		fmt.Println(" ", entry.Name(), entry.IsDir())
	}

	// перешли на несколько уровней вверх
	getFullPath("../../../../..")

	err = os.Chdir("../../..")
	check(err)

	fmt.Println("Visiting subdir")
	err = filepath.Walk("subdir", visit)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func visit(p string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	fmt.Println(" ", p, info.IsDir())
	return nil
}
