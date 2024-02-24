package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	helpers "github.com/alfaa19/go-scraper/helper"
	"github.com/alfaa19/go-scraper/model"
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
)

var url = "https://logammulia.com/id"

func GetItemHandler(c *gin.Context) {
	//Get Query Param
	name := c.Query("name")
	items, err := scrapeData("div.hero-price div[class=content]", url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	//Check If Query Param not equal empty string
	if name != "" {
		item, err := findByName(items, name)
		fmt.Println(name)
		if err != nil {
			c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}
		helpers.ResponseSuccessJson(c, "", item)
		return
	}
	helpers.ResponseSuccessJson(c, "", items)
}

// Find Item struct by Item Name
func findByName(items []model.Item, name string) (model.Item, error) {
	n := strings.ToLower(name)
	for _, v := range items {
		if strings.ToLower(v.Name) == n {
			return v, nil
		}
	}
	return model.Item{}, errors.New("not found")
}

// Scrape data from web, return slice of item & Error
func scrapeData(selector string, url string) ([]model.Item, error) {
	c := colly.NewCollector()
	var items []model.Item

	c.OnHTML(selector, extractData(&items))

	if err := c.Visit(url); err != nil {
		return nil, err
	}

	return items, nil
}

// Extract data from scraper to Struct
func extractData(items *[]model.Item) func(*colly.HTMLElement) {
	return func(h *colly.HTMLElement) {
		var item model.Item

		name := h.ChildText("div.ngc-title")
		price := h.ChildText("p.price span[class=current]")
		lastPrice := h.ChildText("p.last-price")
		change := h.ChildText("p.price span[class=change]")

		if !isEmptyString(name) && !isEmptyString(price) && !isEmptyString(lastPrice) && !isEmptyString(change) {
			//Convert Price to Float
			splitPrice := strings.Split(price, " ")
			splitPrice = strings.Split(splitPrice[1], "Rp")
			splitPrice = strings.Split(splitPrice[1], ",")
			price = strings.ReplaceAll(splitPrice[0], ".", "")
			p, _ := strconv.ParseFloat(price, 64)

			//Convert Last Price to Float
			splitLastPrice := strings.Split(lastPrice, " ")
			splitLastPrice = strings.Split(splitLastPrice[2], "Rp")
			splitLastPrice = strings.Split(splitLastPrice[1], ",")
			lastPrice = strings.ReplaceAll(splitLastPrice[0], ".", "")
			lp, _ := strconv.ParseFloat(lastPrice, 64)

			//convert Change to Float
			trimChange := strings.Trim(change, "Rp")
			trimChange = strings.Split(trimChange, ".")[0]
			trimChange = strings.ReplaceAll(trimChange, ",", "")
			c, _ := strconv.ParseFloat(trimChange, 64)

			//Insert data to struct
			item = model.Item{
				Name:      name,
				Price:     p,
				LastPrice: lp,
				Change:    c,
			}
			*items = append(*items, item)
		}

	}
}

// check if string is empty or not
func isEmptyString(s string) bool {
	return strings.TrimSpace(s) == ""
}
