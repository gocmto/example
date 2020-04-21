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
