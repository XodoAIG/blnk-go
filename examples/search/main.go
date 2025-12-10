package main

import (
	"fmt"
	"net/url"

	blnkgo "github.com/blnkfinance/blnk-go"
)

func main() {
	baseURL, _ := url.Parse("http://localhost:5001/")
	client := blnkgo.NewClient(baseURL, nil)

	normalSearchParams := blnkgo.SearchParams{
		Q:        "*",
		QueryBy:  "transaction_id,reference,description",
		FilterBy: "status:APPLIED",
		SortBy:   "created_at:desc",
		Page:     1,
		PerPage:  10,
	}

	normalResult, resp, err := client.Search.SearchDocument(normalSearchParams, blnkgo.Transactions)
	if err != nil {
		fmt.Printf("Error in normal search: %s\n", err.Error())
		return
	}

	fmt.Printf("Status Code: %d\n", resp.StatusCode)
	fmt.Printf("Found: %d transactions\n", normalResult.Found)
	fmt.Printf("Page: %d\n", normalResult.Page)
	fmt.Printf("Search Time: %dms\n", normalResult.SearchTimeMs)

	for i, hit := range normalResult.Hits {
		doc := hit.Document
		fmt.Printf("\nTransaction %d:\n", i+1)
		fmt.Printf("  ID: %s\n", doc.TransactionID)
		fmt.Printf("  Amount: %.2f\n", doc.Amount)
		fmt.Printf("  Reference: %s\n", doc.Reference)
		fmt.Printf("  Status: %s\n", doc.Status)
		fmt.Printf("  Created At: %s\n", doc.CreatedAt.Time.Format("2006-01-02 15:04:05"))
	}

	groupedSearchParams := blnkgo.SearchParams{
		Q:          "*",
		QueryBy:    "balance_id,currency",
		SortBy:     "created_at:desc",
		Page:       1,
		PerPage:    50,
		GroupBy:    "currency",
		GroupLimit: 5,
	}

	groupedResult, resp, err := client.Search.SearchDocument(groupedSearchParams, blnkgo.Transactions)
	if err != nil {
		fmt.Printf("Error in grouped search: %s\n", err.Error())
		return
	}

	fmt.Printf("Status Code: %d\n", resp.StatusCode)
	fmt.Printf("Found: %d balances\n", groupedResult.Found)
	fmt.Printf("Total Groups: %d\n", len(groupedResult.GroupedHits))
	fmt.Printf("Search Time: %dms\n", groupedResult.SearchTimeMs)

	for i, group := range groupedResult.GroupedHits {
		if len(group.GroupKey) > 0 {
			fmt.Printf("\nGroup %d - Currency: %s\n", i+1, group.GroupKey[0])
		} else {
			fmt.Printf("\nGroup %d - No Group Key\n", i+1)
		}
		fmt.Printf("  Items in group: %d\n", len(group.Hits))

		for j, hit := range group.Hits {
			doc := hit.Document
			fmt.Printf("\n  Balance %d:\n", j+1)
			fmt.Printf("    ID: %s\n", doc.BalanceID)
			fmt.Printf("    Currency: %s\n", doc.Currency)
			fmt.Printf("    Balance: %s\n", doc.Balance)
			fmt.Printf("    Credit: %s\n", doc.CreditBalance)
			fmt.Printf("    Debit: %s\n", doc.DebitBalance)
			fmt.Printf("    Ledger ID: %s\n", doc.LedgerID)
		}
	}
}
