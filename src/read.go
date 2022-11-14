package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type email struct {
	total     float64
	avgDebit  float64
	avgCredit float64
	months    map[int64]*month
}

type month struct {
	numTransactions int
	debit           float64
	credit          float64
}

func ReadFile() (email, error) {

	bucket := os.Getenv("BUCKET")
	key := os.Getenv("KEY")

	// key = "txns.csv"
	// bucket = "transactions-us-east-1-383193739157"

	fmt.Println("bucket" + bucket)
	fmt.Println("key" + key)

	mySession := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")}))
	svc := s3.New(mySession)

	req := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	r, err := svc.GetObject(req)
	if err != nil {
		print(err.Error())
		return email{}, fmt.Errorf("failed to get result: %w", err)
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		print("her2 dd")
		return email{}, fmt.Errorf("error in reading file %s: %s", key, err)
	}

	defer r.Body.Close()

	reader := csv.NewReader(bytes.NewBuffer(body))

	// file, err := os.Open("../data/txns.csv")
	// if err != nil {
	// 	log.Println(err)
	// 	return email{}, err
	// }

	// defer file.Close()

	//eader := csv.NewReader(file)
	if _, err := reader.Read(); err != nil {
		log.Println("here")
	}

	rows, err := reader.ReadAll()

	if err != nil {
		log.Println(err)
		return email{}, err
	}

	fmt.Println("rows :")
	fmt.Println(rows)

	email := email{total: 0, months: map[int64]*month{}}
	for _, row := range rows {
		//id, _ := strconv.ParseInt(row[0], 0, 0)

		monthNum, _ := strconv.ParseInt(strings.Split(row[1], "/")[0], 0, 0)

		if value, ok := email.months[monthNum]; ok {
			value.numTransactions++

		} else {
			email.months[monthNum] = &month{numTransactions: 1, debit: 0, credit: 0}
		}

		temp := email.months[monthNum]

		isDebit := false
		tx := row[2]
		if tx[:1] == "+" {
			isDebit = true
			tx = tx[1:]
		}

		amount, _ := strconv.ParseFloat(tx, 64)
		email.total += amount

		if isDebit {
			temp.debit += amount
		} else {
			temp.credit += amount
		}

	}

	fmt.Println(email.total)

	totalDebit := 0.0
	totalCredit := 0.0

	for num, mon := range email.months {
		mo := *mon
		fmt.Println(num, mo)
		totalDebit += mon.debit
		totalCredit += mon.credit
	}

	email.avgCredit = totalCredit / float64(len(email.months))
	email.avgDebit = totalDebit / float64(len(email.months))

	fmt.Printf("Avg credit: % v , Avg debit %v", totalCredit/float64(len(email.months)), totalDebit/float64(len(email.months)))

	return email, nil

}
