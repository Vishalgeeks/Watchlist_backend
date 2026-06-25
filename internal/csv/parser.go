package csv

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"watchlist-backend/pkg/models"
)

type metaResponse struct {
	SmURL     string `json:"sm_url"`
	Timestamp int64  `json:"timestamp"`
}

func ParseCSV(
	ctx context.Context,
	metaURL string,
) ([]models.Stock, error) {
	//just to check context
	go func() {
		<-ctx.Done()
		log.Println("Context cancelled:", ctx.Err())
	}()

	// Step 1 — Meta URL se actual CSV URL nikalo
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		metaURL,
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var meta metaResponse
	if err := json.Unmarshal(body, &meta); err != nil {
		return nil, err
	}

	log.Println("Actual CSV URL:", meta.SmURL)

	// Step 2 — Actual CSV URL se data fetch karo
	csvReq, err := http.NewRequestWithContext(
		ctx,
		"GET",
		meta.SmURL,
		nil,
	)
	if err != nil {
		return nil, err
	}

	csvResp, err := http.DefaultClient.Do(csvReq)
	if err != nil {
		return nil, err
	}
	defer csvResp.Body.Close()

	rawBody, err := io.ReadAll(csvResp.Body)
	if err != nil {
		return nil, err
	}

	log.Println("First 50 chars:", string(rawBody[:50]))

	// Step 3 — Format detect karo
	var csvBody []byte
	// yaha se hm error k liye detection code likhe th
	// Gzip try karo
	gzReader, err := gzip.NewReader(bytes.NewReader(rawBody))
	if err == nil {
		defer gzReader.Close()

		// Tar try karo andar
		tarReader := tar.NewReader(gzReader)
		for {
			header, err := tarReader.Next()
			if err != nil {
				break
			}
			log.Println("Tar file found:", header.Name)
			if strings.HasSuffix(header.Name, ".csv") {
				csvBody, err = io.ReadAll(tarReader)
				if err != nil {
					return nil, err
				}
				break
			}
		}

		// Tar nahi tha — seedha gzip content
		if csvBody == nil {
			gzReader2, _ := gzip.NewReader(bytes.NewReader(rawBody))
			csvBody, _ = io.ReadAll(gzReader2)
		}
	} else {
		// Gzip bhi nahi — seedha use karo
		csvBody = rawBody
	}

	log.Println("CSV First 300 chars:", string(csvBody[:300]))

	// Step 4 — CSV parse karo
	reader := csv.NewReader(strings.NewReader(string(csvBody)))
	reader.Comma = ','
	reader.LazyQuotes = true

	// Header skip karo
	header, err := reader.Read()
	if err != nil {
		return nil, err
	}
	log.Println("Header columns:", len(header))
	log.Println("Headers:", header[:5]) // pehle 5 columns

	var stocks []models.Stock
	for {
		row, err := reader.Read()
		if err != nil {
			break
		}

		if len(row) < 32 {
			continue
		}

		stock := models.Stock{
			ExchangeInstrumentID:  strings.TrimSpace(row[0]),
			Segment:               strings.TrimSpace(row[1]),
			InstrumentType:        strings.TrimSpace(row[2]),
			Symbol:                strings.TrimSpace(row[3]),
			DisplayName:           strings.TrimSpace(row[4]),
			CompanyName:           strings.TrimSpace(row[5]),
			ISIN:                  strings.TrimSpace(row[6]),
			Series:                strings.TrimSpace(row[7]),
			Exchange:              strings.TrimSpace(row[8]),
			ContractExpiration:    strings.TrimSpace(row[9]),
			OptionType:            strings.TrimSpace(row[11]),
			UnderlyingSymbolID:    strings.TrimSpace(row[12]),
			UnderlyingSymbol:      strings.TrimSpace(row[13]),
			Description:           strings.TrimSpace(row[19]),
			CautionaryMessageInfo: strings.TrimSpace(row[31]),
		}

		stock.Strike = parseFloat(row[10])
		stock.TickSize = parseFloat(row[15])
		stock.UpperCircuit = parseFloat(row[16])
		stock.LowerCircuit = parseFloat(row[17])
		stock.LTP = parseFloat(row[20])
		stock.Open = parseFloat(row[21])
		stock.High = parseFloat(row[22])
		stock.Low = parseFloat(row[23])
		stock.Close = parseFloat(row[24])
		stock.Bid = parseFloat(row[27])
		stock.Ask = parseFloat(row[28])
		stock.LotSize = parseInt(row[14])
		stock.FreezeQty = parseInt(row[18])
		stock.BidQty = parseInt(row[29])
		stock.AskQty = parseInt(row[30])
		stock.Vol = parseInt64(row[25])
		stock.OI = parseInt64(row[26])

		stocks = append(stocks, stock)
	}

	log.Println("Total stocks parsed:", len(stocks))
	return stocks, nil
}

func parseFloat(s string) float64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	val, _ := strconv.ParseFloat(s, 64)
	return val
}

func parseInt(s string) int {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	val, _ := strconv.Atoi(s)
	return val
}

func parseInt64(s string) int64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	val, _ := strconv.ParseInt(s, 10, 64)
	return val
}

//justy to check
