package main

import (
	"github.com/thejini3/go-helper"
	"flag"
	"strconv"
	"regexp"
	"strings"
	"sort"
	"os"
)

type BlogSpotURL struct {
	UrlStr string
	Year   int
	Month  int
	Serial int
}

func main() {

	flag.Parse()

	if len(flag.Args()) != 1 {
		hel.PP("Only file-path [ex: /urls.txt] should be passed!")
	}

	if flag.Arg(0) == "-h" || flag.Arg(0) == "help" ||
		flag.Arg(0) == "--h" || flag.Arg(0) == "--help" {
		hel.PP("Example: $ sort-url-by-path-date urls.txt\nYou should set urls protocol to `https` :)")
	}

	// given file :D
	gf := flag.Arg(0)
	if !hel.FileExists(gf) {
		hel.PP("The file `" + gf + "` doesn't exists!")
	}

	urlStr := hel.GetFileStr(gf)
	urlArr := hel.StrToArr(urlStr, "\n")

	hel.P("given urls count: " + strconv.Itoa(len(urlArr)))

	/*
		https://x.blogspot.com/2014/11/blog-post_2.html
		https://x.blogspot.com/2014/11/blog-post.html
		https://x.blogspot.com/2014/11/blog-post_1.html
	*/

	regexStr := `(https:\/\/.+?\/)(\d+?\/)(\d+?\/)(.+)`
	r := regexp.MustCompile(regexStr)
	var finalUrls []BlogSpotURL

	for _, urlStr := range urlArr {

		/* subStr ex-1

		0: https://idream4life.blogspot.com/2010/12/blog-post.html
		1: https://idream4life.blogspot.com/
		2: 2010/
		3: 12/
		4: blog-post.html

		subStr ex-2

		0: https://idream4life.blogspot.com/2015/12/blog-post_12.html
		1: https://idream4life.blogspot.com/
		2: 2015/
		3: 12/
		4: blog-post_12.html

		 */

		var url = BlogSpotURL{UrlStr: urlStr}

		for i, subStr := range r.FindStringSubmatch(urlStr) {

			// fmt.Printf("%d: %v\n", i, subStr)

			subStrF := subStr[:len(subStr)-1]
			subStrInt, err := strconv.Atoi(subStrF)

			if i == 2 {

				hel.PErrExit("Converting "+subStrF+"  to integer", err)
				url.Year = subStrInt

			} else if i == 3 {

				hel.PErrExit("Converting "+subStrF+"  to integer", err)
				url.Month = subStrInt

			} else if i == 4 {

				subStrArrOfLastItem := strings.Split(subStr, "_")
				/*
				ex-1
				4-0: blog-post.html

				ex-2
				4-0: blog-post
				4-1: 23.html
				 */

				if len(subStrArrOfLastItem) == 1 {
					url.Serial = 0
				} else if len(subStrArrOfLastItem) == 2 {

					// s = 23
					s := strings.ReplaceAll(subStrArrOfLastItem[1], ".html", "")
					serial, err := strconv.Atoi(s)
					hel.PErrExit("Converting "+s+"  to integer", err)
					url.Serial = serial

				} else {
					hel.PP("len(subStrArrOfLastItem) exceeded! for: " + urlStr)
				}

			}

		}

		finalUrls = append(finalUrls, url)

		// hel.P("")
	}

	count := len(finalUrls)
	hel.P("Final urls count: " + strconv.Itoa(count))

	var sorted []BlogSpotURL

	for _, year := range hel.SortIntAsc(getUniqueYears(finalUrls)) {

		// fmt.Printf("%v = ", year)

		for _, month := range hel.SortIntAsc(getMonthsOfYear(year, finalUrls)) {
			// fmt.Printf("%v-%v\n", year, month)

			for _, finalSortedUrl := range sortSerial(getAllOfYearAndMonth(year, month, finalUrls)) {
				sorted = append(sorted, finalSortedUrl)
			}

		}

		// hel.P("")
	}

	hel.P("[SORTED] Final urls count: " + strconv.Itoa(len(sorted)))

	var newFileStr string
	for i, url := range sorted {
		// fmt.Printf("%+v\n", url)
		newFileStr += url.UrlStr
		if len(sorted)-1 != i {
			newFileStr += "\n"
		}
	}

	fn := hel.GetNonCreatedFileName("sorted_urls", ".txt", 1)
	f, _ := os.Create(fn)
	f.WriteString(newFileStr)
	f.Close()

	hel.P("Generated file: " + fn)
}

func sortSerial(urls []BlogSpotURL) []BlogSpotURL {
	sort.Slice(urls, func(i, j int) bool {
		return urls[i].Serial < urls[j].Serial
	})
	return urls
}

func getUniqueYears(urls []BlogSpotURL) []int {
	var uniques []int
	for _, url := range urls {
		if !hel.ContainsInt(uniques, url.Year) {
			uniques = append(uniques, url.Year)
		}
	}
	return uniques
}

func getMonthsOfYear(year int, urls []BlogSpotURL) []int {
	var uniques []int
	for _, url := range urls {
		if url.Year == year && !hel.ContainsInt(uniques, url.Month) {
			uniques = append(uniques, url.Month)
		}
	}
	return uniques
}

func getAllOfYearAndMonth(year int, month int, urls []BlogSpotURL) []BlogSpotURL {
	var uniques []BlogSpotURL
	for _, url := range urls {
		if url.Year == year && url.Month == month {
			uniques = append(uniques, url)
		}
	}
	return uniques
}
