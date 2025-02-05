package utils

import (
	"math"
	"receipt-processor/models"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func CalculatePoints(receipt models.Receipt) int {
	points := 0

	re := regexp.MustCompile(`[a-zA-Z0-9]`)
	points += len(re.FindAllString(receipt.Retailer, -1))

	total, _ := strconv.ParseFloat(receipt.Total, 64)

	if total == math.Floor(total) {
		points += 50
	}

	if math.Mod(total, 0.25) == 0 {
		points += 25
	}

	points += (len(receipt.Items) / 2) * 5

	for _, item := range receipt.Items {
		trimmedDesc := strings.TrimSpace(item.ShortDescription)
		if len(trimmedDesc)%3 == 0 {
			price, _ := strconv.ParseFloat(item.Price, 64)
			points += int(math.Ceil(price * 0.2))
		}
	}

	if total > 10.00 {
		points += 5
	}

	date, _ := time.Parse("2006-01-02", receipt.PurchaseDate)
	if date.Day()%2 != 0 {
		points += 6
	}

	timeParsed, _ := time.Parse("15:04", receipt.PurchaseTime)
	if timeParsed.Hour() == 14 || (timeParsed.Hour() == 15 && timeParsed.Minute() < 60) {
		points += 10
	}

	return points
}
