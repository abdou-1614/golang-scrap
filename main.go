package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

var Link_list = []string{
	"https://in.investing.com/equities/axis-bank-technical",
	"https://in.investing.com/equities/tata-motors-ltd-technical",
	"https://in.investing.com/equities/icici-bank-ltd-technical",
	"https://in.investing.com/equities/housing-development-finance-technical",
	"https://in.investing.com/equities/maruti-suzuki-india-technical",
	"https://in.investing.com/equities/infosys-technical",
	"https://in.investing.com/equities/kotak-mahindra-bank-technical",
	"https://in.investing.com/equities/reliance-industries-technical",
	"https://in.investing.com/equities/hdfc-bank-ltd-technical",
	"https://in.investing.com/equities/adani-enterprises-technical",
	"https://in.investing.com/equities/mundra-port-special-eco.-zone-technical",
	"https://in.investing.com/equities/apollo-hospitals-technical",
	"https://in.investing.com/equities/asian-paints-technical",
	"https://in.investing.com/equities/bajaj-auto-technical",
	"https://in.investing.com/equities/bajaj-finance-technical",
	"https://in.investing.com/equities/bajaj-finserv-limited-technical",
	"https://in.investing.com/equities/bharat-petroleum-technical",
	"https://in.investing.com/equities/bharti-airtel-technical",
	"https://in.investing.com/equities/cipla-technical",
	"https://in.investing.com/equities/coal-india-technical",
	"https://in.investing.com/equities/divis-laboratories-technical",
	"https://in.investing.com/equities/dr-reddys-laboratories-technical",
	"https://in.investing.com/equities/grasim-industries-technical",
	"https://in.investing.com/equities/hcl-technologies-technical",
	"https://in.investing.com/equities/hdfc-bank-ltd-technical",
	"https://in.investing.com/equities/itc-technical",
	"https://in.investing.com/equities/indusind-bank-technical",
	"https://in.investing.com/equities/infosys-technical",
	"https://in.investing.com/equities/jsw-steel-technical",
	"https://in.investing.com/equities/kotak-mahindra-bank-technical",
	"https://in.investing.com/equities/larsen---toubro-technical",
	"https://in.investing.com/equities/mahindra---mahindra-technical",
	"https://in.investing.com/equities/ntpc-technical",
	"https://in.investing.com/equities/oil---natural-gas-corporation-technical",
	"https://in.investing.com/equities/power-grid-corp.-of-india-technical",
	"https://in.investing.com/equities/apollo-hospitals-technical",
	"https://in.investing.com/equities/wipro-ltd-technical",
	"https://in.investing.com/equities/united-phosphorus-technical",
	"https://in.investing.com/equities/state-bank-of-india-technical",
	"https://in.investing.com/equities/sbi-life-insurance-technical",
	"https://in.investing.com/equities/sun-pharma-advanced-research-technical",
	"https://in.investing.com/equities/tata-consultancy-services-technical",
	"https://in.investing.com/equities/tata-global-beverages-technical",
	"https://in.investing.com/equities/tata-steel-technical",
	"https://in.investing.com/equities/titan-industries-technical",
}

func main() {

	sendRequest()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("GET REQUEST")
	})

	go func() {
		err := http.ListenAndServe(":4004", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	scrape()
}

func sendRequest() {
	res, err := http.Get("https://buysell-test.onrender.com/")

	if err != nil {
		fmt.Printf("ERROR TO SEND REQESUT : %v\n", err)
	} else {
		defer res.Body.Close()
		fmt.Println("Request sent successfully.")
	}

}

func scrape() {
	for {
		for _, link := range Link_list {
			buy := true
			sell := true
			focusCount := 0

			results, err := getResult(link)
			if err != nil {
				log.Printf("Error getting results: %v\n", err)
				continue
			}

			for _, result := range results {
				fmt.Println(result.Status)
				if result.Status == "Strong Sell" {
					focusCount++
					if sell && focusCount > 2 {
						fmt.Println("Telegram")
						sendAlertToTG(fmt.Sprintf("Alert for Bank %s - \"STRONG SELL\"", result.BankName))
						sell = false
						focusCount = 0
					}
				}
				if result.Status == "Strong Buy" {
					focusCount++
					if buy && focusCount > 2 {
						sendAlertToTG(fmt.Sprintf("Alert for Bank %s - \"STRONG BUY\"", result.BankName))
						buy = false
						focusCount = 0
					}
				}
			}

			fmt.Println("\n\nWaiting for 10 seconds")
			waitBeforeNextIteration(10 * time.Second)
		}
	}
}

func getResult(link string) ([]Result, error) {
	result := make([]Result, 0)
	timeFrame := []int{1, 5, 15}

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Headless,
		chromedp.NoSandbox,
		chromedp.Flag("disable-dev-shm-usage", true),
	)
	allowCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)

	defer cancel()

	ctx, cancel := chromedp.NewContext(allowCtx)

	defer cancel()

	for _, i := range timeFrame {
		url := fmt.Sprintf("%s?timeFrame=%d", link, i*60)
		fmt.Printf("GETTING, %s\n", url)

		err := chromedp.Run(ctx,
			network.Enable(),
			chromedp.Navigate(url),
			chromedp.WaitVisible("section.forecast-box-graph", chromedp.ByQuery),
		)

		if err != nil {
			log.Printf("Error navigating to URL: %v\n", err)
		}

		fmt.Println("GOTO SUCCESS")

		var status, bankName string

		err = chromedp.Run(ctx,
			chromedp.TextContent("section.forecast-box-graph .title", &status, chromedp.ByQuery),
			chromedp.TextContent("h1.main-title.js-main-title", &bankName, chromedp.ByQuery),
		)
		if err != nil {
			log.Printf("Error getting text content: %v\n", err)
			continue
		}
		result = append(result, Result{
			BankName: bankName,
			Status:   status,
			Link:     link,
			Url:      url,
		})
		log.Printf("BANK NAME: %v\n STATUS: %v\n LINK: %v\n URL: %v\n", bankName, status, link, url)
		fmt.Println("SEND RESULT SUCCESS.")
	}
	return result, nil
}

func sendAlertToTG(alertMsg string) {
	fmt.Println("Sending")

	alertBot := "https://api.telegram.org/bot5762212585:5a58d66d-8a32-45cc-8b16-7d3e6491d60c"
	chatID := "-855310893"
	alertText := alertMsg

	resp, err := http.Get(fmt.Sprintf("%s/sendMessage?chat_id=%s&text=%s", alertBot, chatID, alertText))
	if err != nil {
		log.Printf("Error sending alert to Telegram: %v\n", err)
	} else {
		defer resp.Body.Close()
		fmt.Printf("[Response] - %d\n", resp.StatusCode)
	}
}

func waitBeforeNextIteration(duration time.Duration) {
	time.Sleep(duration)
}

type Result struct {
	BankName string `json:"bankName"`
	Status   string `json:"status"`
	Link     string `json:"link"`
	Url      string `json:"url"`
}
