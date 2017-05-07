package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"net/http"

	"database/sql"
	"fmt"
	"regexp"
	"time"
	"errors"
	"strings"
	"strconv"

	"log"

	"bufio"
)

type Tax struct {
	Id           int
	Municipality string
	PeriodStart  time.Time
	PeriodEnd    time.Time
	Rate         float64
}

func main() {

	db, err := sql.Open("mysql", "tax:tax@tcp(127.0.0.1:3306)/tax")
	if err != nil {
		fmt.Print(err.Error())
	}
	defer db.Close()
	// make sure connection is available
	err = db.Ping()
	if err != nil {
		fmt.Print(err.Error())
	}

	router := gin.Default()

	router.GET("/settax/:municipality/:date/:tax", func(c *gin.Context) {

		tax := new(Tax)
		tax.Municipality = strings.ToLower(c.Param("municipality"))
		taxDate := c.Param("date")

		tax.Rate, err = strconv.ParseFloat(c.Param("tax"), 64)

		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, "Internal server error")
			return
		}

		err = CalculateTaxDates(taxDate, tax)

		fmt.Println(tax.PeriodStart, tax.PeriodEnd);

		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		err = InsertTax(db, tax)

		if err != nil {
			c.JSON(http.StatusBadRequest, "Internal server error")
			return
		}

		c.JSON(http.StatusOK, "Record saved")

	})

	router.POST("/upload", func(c *gin.Context) {

		tax := new(Tax)

		file, header , err := c.Request.FormFile("upload")

		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(header.Filename)

		defer file.Close()

		scanner := bufio.NewScanner(file)
		bad := 0
		good := 0
		for scanner.Scan() {
			line := scanner.Text()
			cols := strings.Split(line,",")

			if (len(cols) != 4) {
				fmt.Println("Not correct line: " + line)
				bad++
				continue
			}

			if (cols[2] != "" && cols[2] != "week") {
				fmt.Println("Not correct column 2 in line : " + cols[2])
				bad++
				continue
			}

			err := CalculateTaxDates(cols[1] + cols[2], tax)

			if (err != nil) {
				bad++
				fmt.Println("Not correct date in line : " + line)
				continue
			}

			tax.Municipality = strings.ToLower(cols[0])
			tax.Rate, err = strconv.ParseFloat(cols[3], 64)

			if (err != nil) {
				bad++
				fmt.Println("Not correct rate in line : " + line)
				continue
			}

			err = InsertTax(db, tax)

			if err != nil {
				bad++
				fmt.Println("Internal error processing line : " + line)
				continue
			}


			good ++

		}

		c.JSON(http.StatusOK, "Processed: " + strconv.Itoa(good) + ". Not correct lines: " + strconv.Itoa(bad))
	})

	// GET tax detail
	router.GET("/tax/:municipality/:date", func(c *gin.Context) {
		var tax  Tax

		municipality := strings.ToLower(c.Param("municipality"))
		taxDate := c.Param("date")

		matched, matchedRes := regexp.MatchString("^[0-9]{4}-[0-9]{2}-[0-9]{2}$", taxDate)

		if matched && matchedRes == nil {
			requestDate, err := time.Parse("2006-01-02", taxDate)

			fmt.Println(requestDate)
			if err != nil {
				c.JSON(http.StatusBadRequest, "Not correct date")
				return
			}
		} else {
			c.JSON(http.StatusBadRequest, "Not correct date")
			return
		}

		row := db.QueryRow("select Rate "+
			"from Tax where ? BETWEEN PeriodStart and PeriodEnd and Municipality = ? "+
			"order by DATEDIFF(PeriodEnd, PeriodStart), PeriodStart limit 1;", taxDate, municipality)
		err = row.Scan(&tax.Rate)

		if err != nil {
			c.JSON(http.StatusOK, "No data")
			return
		}

		c.JSON(http.StatusOK, tax.Rate)
	})

	router.Run(":3000")
}

func InsertTax(db *sql.DB, tax *Tax) (err error) {
	// Could use update on duplicate, but too many indexes
	stmtDel, err := db.Prepare("DELETE FROM TAX WHERE Municipality = ? AND PeriodStart = ? AND PeriodEnd = ?")

	if err != nil {
		fmt.Println(err)
		return err
	}

	res, err := stmtDel.Exec(tax.Municipality, tax.PeriodStart, tax.PeriodEnd)
	fmt.Println(res)

	if err != nil {
		fmt.Println(err)
		return err
	}

	stmtIns, err := db.Prepare("INSERT INTO TAX (Municipality, PeriodStart, PeriodEnd, Rate) VALUES (?, ?, ?, ?)")

	if err != nil {
		fmt.Println(err)
		return err
	}

	res, err = stmtIns.Exec(tax.Municipality, tax.PeriodStart, tax.PeriodEnd, tax.Rate)
	fmt.Println(res)

	if err != nil {
		fmt.Println(err)
		return err
	}

	id, err := res.LastInsertId()

	fmt.Println(id)

	return
}

func CalculateTaxDates(taxDate string, tax *Tax) (errMsg error) {
	var err error
	matched, matchedRes := regexp.MatchString("^[0-9]{4}$", taxDate)

	if matched && matchedRes == nil {
		tax.PeriodStart, err = time.Parse("2006-01-02", taxDate+"-01-01")

		if err != nil {
			fmt.Println(err)
			return errors.New("Not correct year")
		}

		tax.PeriodEnd = tax.PeriodStart.AddDate(1, 0, -1)
	} else {

		matched, matchedRes := regexp.MatchString("^[0-9]{4}-[0-9]{2}$", taxDate)

		if matched && matchedRes == nil {
			tax.PeriodStart, err = time.Parse("2006-01-02", taxDate+"-01")

			if err != nil {
				fmt.Println(err)
				return errors.New("Not correct year and month")
			}

			tax.PeriodEnd = tax.PeriodStart.AddDate(0, 1, -1)
		} else {
			matched, matchedRes := regexp.MatchString("^[0-9]{4}-[0-9]{2}-[0-9]{2}$", taxDate)

			if matched && matchedRes == nil {
				tax.PeriodStart, err = time.Parse("2006-01-02", taxDate)

				if err != nil {
					fmt.Println(err)
					return errors.New("Not correct date")
				}

				tax.PeriodEnd = tax.PeriodStart
				fmt.Println(tax.PeriodStart)
			} else {
				matched, matchedRes := regexp.MatchString("^[0-9]{4}-[0-9]{2}-[0-9]{2}week$", taxDate)

				if matched && matchedRes == nil {
					tax.PeriodStart, err = time.Parse("2006-01-02", taxDate[:	10])

					if err != nil {
						fmt.Println(err)
						return errors.New("Not correct week date")
					}

					tax.PeriodEnd = tax.PeriodStart.AddDate(0, 0 , 6)
				} else {
					return errors.New("It must be provided year, year-month, year-month-day or year-month-day-week")
				}

			}

		}

	}
	return
}
