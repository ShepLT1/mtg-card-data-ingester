package ingestion

import (
	"SWE1-project-data-ingester/internal/data"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"strings"
)

type ExternalCard struct {
	Name					string						`json:"name"`
	SetName				string						`json:"set_name"`
	SetCode				string						`json:"set"`
	CollectorNum	string						`json:"collector_number"`
	PromoTypes		[]string					`json:"promo_types"`
	Finishes			[]string					`json:"finishes"`
	Prices				map[string]string	`json:"prices"`
	ImageURIs			map[string]string	`json:"image_uris"`
	Digital				bool							`json:"digital"`
	FrameEffects	[]string					`json:"frame_effects"`
	Border				string						`json:"border_color"`
	Language			string						`json:"lang"`
}

type IngestCardFailure struct {
	Data			any			`json:"data"`
	Msg				error		`json:"msg"`	
}

type BulkDataItem struct {
	Name				string	`json:"name"`
	DownloadURI	string	`json:"download_uri"`
}

type BulkDataResponse struct {
	Data	[]BulkDataItem	`json:"data"`
}

func FetchDefaultCardEndpoint() (string, error) {
	bulkUrl := "https://api.scryfall.com/bulk-data"
	res, err := http.Get(bulkUrl)
	if err != nil {
		return "", fmt.Errorf("failed to GET bulk data endpoint: %w", err)
	}
	defer res.Body.Close()

	var bulkDataRes	BulkDataResponse
	if err := json.NewDecoder(res.Body).Decode(&bulkDataRes); err != nil {
		return "", fmt.Errorf("failed to decode JSON: %w", err)
	}

	if len(bulkDataRes.Data) < 3 {
		return "", nil
	}

	return bulkDataRes.Data[2].DownloadURI, nil
}

func FetchExternalCards(url string) ([]ExternalCard, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to GET bulk data: %w", err)
	}
	defer res.Body.Close()

	var cards []ExternalCard
	if err := json.NewDecoder(res.Body).Decode(&cards); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	return cards, nil
}

func IngestCardData(ctx context.Context, extCards []ExternalCard) ([]IngestCardFailure, error) {
	var failures []IngestCardFailure

	for i := 0; i < len(extCards); i++ {
		if extCards[i].Digital || extCards[i].Language != "en" || (extCards[i].Prices["usd"] == "" && extCards[i].Prices["usd_foil"] == "" && extCards[i].Prices["usd_etched"] == "") {
			continue
		}
		
		set := &data.Set{
			Name:	extCards[i].SetName,
			Code:	extCards[i].SetCode,
		}
		err := data.SaveSetUpsert(ctx, set)
		if err != nil {
			var failure IngestCardFailure
			failure.Data = extCards[i]
			failure.Msg = fmt.Errorf("failed to PUT SET: %w", err)
			failures = append(failures, failure)
			continue
		}

		card := &data.Card{
			Name: extCards[i].Name,
			SetID: set.ID,
			CollectorNum: extCards[i].CollectorNum,
			ImageURI: extCards[i].ImageURIs["normal"],
		}

		for j := 0; j < len(extCards[i].PromoTypes); j++ {
			switch extCards[i].PromoTypes[j] {
				case "prerelease":
					card.PromoType = "prerelease"
				case "promopack":
					card.PromoType = "promopack"
				case "serialized":
					card.AltStyle = "serialized"
				default:
				if strings.Contains(extCards[i].PromoTypes[j], "foil") {
					if card.AltStyle != "serialized" {
						card.AltStyle = extCards[i].PromoTypes[j]
					}
				}
			}
		}

		for k := 0; k < len(extCards[i].FrameEffects); k++ {
			switch extCards[i].FrameEffects[k] {
				case "showcase", "extendedart", "fullart":
					card.AltStyle = extCards[i].FrameEffects[k]
			}
		}

		if card.AltStyle == "" && extCards[i].Border == "borderless" {
			card.AltStyle = "borderless"
		}

		var priceTypes []string

		if extCards[i].Prices["usd"] != "" {
			priceTypes = append(priceTypes, "usd")
		}
		if extCards[i].Prices["usd_foil"] != "" {
			priceTypes = append(priceTypes, "usd_foil")
		}
		if extCards[i].Prices["usd_etched"] != "" {
			priceTypes = append(priceTypes, "usd_etched")
		}

		now := time.Now().UTC()

		listing := &data.Listing{
			CreatedAt: now,
			CreatedDate: now.Truncate(24 * time.Hour),
		}

		for j := 0; j < len(priceTypes); j++ {
			switch priceTypes[j] {
				case "usd":
					card.Finish = "nonfoil"
					floatPrice, err := strconv.ParseFloat(extCards[i].Prices["usd"], 64)
					if err != nil {
						var failure IngestCardFailure
						failure.Data = *card
						failure.Msg = fmt.Errorf("error converting string to float64: %w", err)
						failures = append(failures, failure)
						continue
					}
					listing.Price = floatPrice
				case "usd_foil":
					card.Finish = "foil"
					floatPrice, err := strconv.ParseFloat(extCards[i].Prices["usd_foil"], 64)
					if err != nil {
						var failure IngestCardFailure
						failure.Data = *card
						failure.Msg = fmt.Errorf("error converting string to float64: %w", err)
						failures = append(failures, failure)
						continue
					}
					listing.Price = floatPrice
				case "usd_etched":
					card.Finish = "foil_etched"
					floatPrice, err := strconv.ParseFloat(extCards[i].Prices["usd_etched"], 64)
					if err != nil {
						var failure IngestCardFailure
						failure.Data = *card
						failure.Msg = fmt.Errorf("error converting string to float64: %w", err)
						failures = append(failures, failure)
						continue
					}
					listing.Price = floatPrice
			}
			card.ID = 0
			cardErr := data.SaveCardUpsert(ctx, card)
			if cardErr != nil {
				var failure IngestCardFailure
				failure.Data = *card
				failure.Msg = fmt.Errorf("failed to PUT CARD: %w", cardErr)
				failures = append(failures, failure)
				continue
			}
			listing.CardID = card.ID
			listing.ID = 0
			listingErr := data.SaveListing(ctx, listing)
			if listingErr != nil {
				var failure IngestCardFailure
				failure.Data = card
				failure.Msg = fmt.Errorf("failed to PUT LISTING: %w", listingErr)
				failures = append(failures, failure)
				continue
			}
		}

	}

	return failures, nil
}
