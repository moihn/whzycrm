package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"

	"github.com/integrii/flaggy"
)

type ConnectionConfig struct {
	Host     string `json:"host"`
	Port     *int   `json:"port"`
	Database string `json:"db"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (config ConnectionConfig) getConnectionString() string {
	portString := ":3306"
	if config.Port != nil {
		portString = fmt.Sprintf(":%v", *config.Port)
	}
	if len(config.Password) == 0 {
		logrus.Fatal("Password is not found in connection configuration.")
	}
	return fmt.Sprintf("%v:%v@tcp(%v%v)/%v", config.Username, config.Password, config.Host, portString, config.Database)
}

func getProductRef(fileName string) string {
	return strings.Split(fileName, " ")[0]
}

func main() {
	var imageFolder string
	var vendorID int = 0
	imageFileNameRegex := ".*(jpg|JPG|png|PNG)$"
	flaggy.String(&imageFolder, "", "directory", "Image folder")
	flaggy.String(&imageFileNameRegex, "", "pattern", "Image filename regex pattern")
	flaggy.Int(&vendorID, "", "vendor", "Vendor ID in the system.")
	flaggy.Parse()

	logrus.SetLevel(logrus.DebugLevel)

	imageFileNameRegexp, err := regexp.Compile(imageFileNameRegex)
	if err != nil {
		logrus.Fatalf("Failed to parse regex for file name %v: %v", imageFileNameRegex, err)
	}

	// Open our jsonFile
	connFile, err := os.Open("connection.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	logrus.Println("Successfully Opened connection.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer connFile.Close()
	// read our opened jsonFile as a byte array.
	byteValue, _ := ioutil.ReadAll(connFile)

	connConfig := ConnectionConfig{}
	json.Unmarshal(byteValue, &connConfig)

	// Open up our database connection.
	// I've set up a database on my local machine using phpmyadmin.
	// The database is called testDb
	connStr := connConfig.getConnectionString()
	logrus.Debug("Connecting to ", connStr)
	db, err := sql.Open("mysql", connStr)

	// if there is an error opening the connection, handle it
	if err != nil {
		panic(err.Error())
	}
	logrus.Info("Successfully Opened database.")
	// defer the close till after the main function has finished
	// executing
	defer db.Close()

	dbTransact, err := db.Begin()
	if err != nil {
		logrus.Fatalf("Failed to open database transaction: %v", err)
	}

	// scan a directory and get file name
	dirEntries, err := ioutil.ReadDir(imageFolder)
	if err != nil {
		logrus.Fatalf("Failed to open image directoriy %v", imageFolder)
	}
	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			continue
		}

		if imageFileNameRegexp.MatchString(dirEntry.Name()) {
			logrus.Infof("Found %v", dirEntry.Name())
			productRef := getProductRef(filepath.Base(dirEntry.Name()))
			logrus.Infof("Found product reference %v", productRef)
			productID := ensureProductDefined(dbTransact, productRef)
			logrus.Infof("Product %v has ID %v", productRef, productID)
			prodInVendorID := ensureProductInVendorDefined(dbTransact, productID, vendorID, productRef)
			logrus.Infof("Product %v has ID %v in vendor %v", productRef, prodInVendorID, vendorID)
			// load the image file content
			filePath := path.Join(imageFolder, dirEntry.Name())
			// check if the file has been uploaded using product_in_vendor_id and md5 checksum
			picId := ensureImageUploaded(dbTransact, prodInVendorID, filePath)
			logrus.Infof("Product image has ID %v in product_in_vendor_picture", picId)
		}
	}

	dbTransact.Commit()
}

func ensureImageUploaded(transact *sql.Tx, prodInVendorID int, filePath string) int {
query:
	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		logrus.Fatalf("Failed to read data from file %v: %v", filePath, err)
	}
	fileMd5 := md5.Sum(fileData)
	fileMd5Str := fmt.Sprintf("%x", fileMd5)
	rows, err := transact.Query(`
		select PICTURE_ID
		from PRODUCT_IN_VENDOR_PICTURE
		where PRODUCT_IN_VENDOR_ID = ?
		  and CHECK_SUM = ?
	`, prodInVendorID, fileMd5Str)
	if err != nil {
		logrus.Fatalf("Failed to check existence of image file %v in database: %v", filePath, err)
	}
	defer rows.Close()
	if rows.Next() {
		var picID int
		rows.Scan(&picID)
		return picID
	}
	// picture is not yet uploaded, let's insert
	_, err = transact.Exec(`
		insert into PRODUCT_IN_VENDOR_PICTURE
		(PRODUCT_IN_VENDOR_ID, FORMAT_TYPE, BODY, CHECK_SUM)
		values
		(?, (select FORMAT_ID from FILE_FORMAT where FILE_NAME_SUFFIX = ?), ?, ?)
	`, prodInVendorID, strings.ToLower(path.Ext(filePath)), fileData, fileMd5Str)
	if err != nil {
		logrus.Fatalf("Failed to add new product in vendor image %v: %v", filePath, err)
	}
	goto query
}

func ensureProductDefined(transact *sql.Tx, productRef string) int {
query:
	rows, err := transact.Query(`
		select PRODUCT_ID
		from PRODUCT
		where REFERENCE = ?
	`, productRef)
	if err != nil {
		logrus.Fatalf("Error occurred when finding existing product using reference %v: %v", productRef, err)
	}
	defer rows.Close()
	if rows.Next() {
		// we have the product
		var productID int
		rows.Scan(&productID)
		return productID
	}

	// product ref is not defined, let's add it
	_, err = transact.Exec(`
		insert into PRODUCT
		(REFERENCE, NAME, TYPE)
		values
		(?, ?, ?)
	`, productRef, productRef, 0)
	if err != nil {
		logrus.Fatalf("Failed to add new product with reference %v: %v", productRef, err)
	}
	goto query
}

func ensureProductInVendorDefined(transact *sql.Tx, productID int, vendorID int, productInVendorRef string) int {
query:
	rows, err := transact.Query(`
		select PRODUCT_IN_VENDOR_ID
		from PRODUCT_IN_VENDOR
		where PRODUCT_ID = ?
	`, productID)
	if err != nil {
		logrus.Fatalf("Error occurred when finding existing product(ID=%v,Ref=%v) in vendor(ID=%v): %v", productID, productInVendorRef, vendorID, err)
	}
	defer rows.Close()
	if rows.Next() {
		// we have the product in vendor
		var productInVendorID int
		rows.Scan(&productInVendorID)
		return productInVendorID
	}

	// product is not defined for this vendor yet, let's defined it now
	_, err = transact.Exec(`
		insert into PRODUCT_IN_VENDOR
		(PRODUCT_ID, VENDOR_ID, REFERENCE_IN_VENDOR)
		values
		(?, ?, ?)
	`, productID, vendorID, productInVendorRef)
	if err != nil {
		logrus.Fatalf("Failed to add new product(ID=%v,Ref=%v) to vendor(ID=%v): %v", productID, productInVendorRef, vendorID, err)
	}
	goto query
}
