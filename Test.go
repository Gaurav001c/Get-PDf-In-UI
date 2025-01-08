package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo"
)

func main() {

	e := echo.New()

	e.Static("/", "static")

	e.Static("/uploads", "uploads")
	e.GET("/", serverIndex)
	e.GET("/home", serveHome)
	e.GET("/about", serveAbout)
	e.GET("/contact", serveContact)
	e.POST("/upload", handleUpload)
	e.GET("/readData", readFileData)

	// e.GET("/get/:id/:name", getValue)

	e.Logger.Fatal(e.Start(":8080"))

}

// func getValue(c echo.Context) error {

// 	// id := c.QueryParam("id")
// 	// name := c.QueryParam("name")

// 	id := c.Param("id")
// 	name := c.Param("name")
// 	fmt.Println(id, name)
// 	return nil

// }

func serverIndex(c echo.Context) error {

	fmt.Println("IP Address" + c.RealIP())
	fmt.Println("index Page")
	// return c.Echo().Static("")
	return c.File("static/index.html")

}

func serveHome(c echo.Context) error {
	fmt.Println("Home Page")
	return c.File("static/home.html")
}

func serveAbout(c echo.Context) error {
	return c.File("static/about.html")
}

func serveContact(c echo.Context) error {
	return c.File("static/contact.html")
}
func handleUpload(c echo.Context) error {

	fmt.Println("In Upload")

	file, err := c.FormFile("file")
	if err != nil {
		panic(err)
	}

	src, err := file.Open()

	if err != nil {
		panic(err)
	}

	fmt.Println("Src Pass : ", src)
	defer src.Close()

	dstPath := filepath.Join("./uploads", file.Filename)

	dst, _ := os.Create(dstPath)

	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return c.String(http.StatusOK, fmt.Sprintf("File %s uploaded successfully", file.Filename))
}

// func readFileData(c echo.Context) error {
// 	fileName := c.FormValue("fileName")

// 	data, err := os.ReadFile(filepath.Join("./uploads", fileName))
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 		return c.String(http.StatusBadRequest, fmt.Sprintf("Error reading file: %s", err))
// 	}

// 	return c.File(string(data))
// }

func readFileData(c echo.Context) error {
	fileName := c.QueryParam("fileName")

	// fmt.Println(fileName)

	filePath := filepath.Join("uploads", fileName)

	fmt.Println(filePath)
	// fmt.Println(filePath)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return c.String(http.StatusOK, fmt.Sprintf("Error reading file: %v", err))
	}

	ext := filepath.Ext(fileName) //get file extension

	fmt.Println("extension : ", ext)

	var response string
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif":
		response = fmt.Sprintf(`<img src="/uploads/%s" alt="Image">`, fileName)
	case ".pdf":
		response = fmt.Sprintf(`<embed src="/uploads/%s" type="application/pdf" width="600" height="500">`, fileName)

	default:
		fileContent := string(data)
		response = fmt.Sprintf(`<pre>%s</pre>`, fileContent)
	}

	return c.HTML(http.StatusOK, response)
}
