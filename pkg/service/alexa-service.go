package service

import (
	"electricity-prices/pkg/i18n"
	"electricity-prices/pkg/model"
	"electricity-prices/pkg/utils"
	"fmt"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"log"
	"strings"
	"time"
)

func init() {
	// Initialise the i18n translations
	i18n.InitialiseTranslations()
}

func GetTitle(lang language.Tag) string {
	p := message.NewPrinter(lang)
	return p.Sprintf("alexa_full_title")
}

func GetFullFeed(date time.Time, lang language.Tag) (string, error) {
	// Get the daily info for the given date
	dailyInfo, err := GetDailyInfo(date)

	if err != nil {
		return "", err
	}
	if len(dailyInfo.Prices) == 0 {
		return getTodayNoDataMessage(lang), nil
	}
	var messages []string

	// Parse day rating message
	messages = append(messages, getTodayRatingMessage(dailyInfo.DayRating, dailyInfo.DayAverage, lang))

	// Get current price
	cpMsg, err := getCurrentPriceMessage(dailyInfo.Prices, lang)
	if err == nil {
		messages = append(messages, cpMsg)
	}

	// Get next cheap period
	messages = append(messages, getNextCheapPeriodMessage(dailyInfo.CheapPeriods, lang))

	// Get next expensive period
	messages = append(messages, getNextExpensivePeriodMessage(dailyInfo.ExpensivePeriods, lang))

	// Get tomorrow's data
	tomorrowInfo, err := GetDailyInfo(date.AddDate(0, 0, 1))
	if err == nil && len(tomorrowInfo.Prices) > 0 {
		messages = append(messages, getTomorrowRatingMessage(tomorrowInfo.DayRating, tomorrowInfo.DayAverage, lang))
	}

	return strings.Join(messages, " "), nil
}

func getTodayNoDataMessage(lang language.Tag) string {
	p := message.NewPrinter(lang)
	noData := p.Sprintf("alexa_today_nodata")

	return noData
}

func getTodayRatingMessage(dayRating model.DayRating, dayAverage float64, lang language.Tag) string {
	p := message.NewPrinter(lang)

	if dayRating == model.Nil {
		return ""
	}

	rating := p.Sprintf(fmt.Sprintf("alexa_rating_%s", strings.ToLower(string(dayRating))))
	todayRating := p.Sprintf("alexa_today_rating", rating, utils.FormatPrice(dayAverage))

	return todayRating
}

func getTomorrowRatingMessage(dayRating model.DayRating, dayAverage float64, lang language.Tag) string {
	p := message.NewPrinter(lang)

	if dayRating == model.Nil {
		return ""
	}

	rating := p.Sprintf(fmt.Sprintf("alexa_rating_%s", strings.ToLower(string(dayRating))))
	log.Printf("Rating: %s", rating)
	tomorrowRating := p.Sprintf("alexa_tomorrow_rating", rating, utils.FormatPrice(dayAverage))

	return tomorrowRating
}

func getCurrentPriceMessage(prices []model.Price, lang language.Tag) (string, error) {
	p := message.NewPrinter(lang)
	now := time.Now()

	for _, price := range prices {
		if utils.SameHour(now, price.DateTime) {
			return p.Sprintf("alexa_current_price", utils.FormatPrice(price.Price)), nil
		}
	}
	return "", fmt.Errorf("no current price found")
}

// getNextCheapPeriodMessage
// Get the next cheap period message.
func getNextCheapPeriodMessage(periods [][]model.Price, lang language.Tag) string {
	p := message.NewPrinter(lang)

	next, started := utils.GetNextPeriod(periods, time.Now())

	if next == nil {
		return p.Sprintf("alexa_next_cheap_period_nodata")
	}

	avg := utils.FormatPrice(utils.CalculateAverage(next))
	start := utils.FormatTime(next[0].DateTime)
	end := utils.FormatTime(next[len(next)-1].DateTime)

	if started {
		return p.Sprintf("alexa_current_cheap_period", start, avg, end)
	} else {
		return p.Sprintf("alexa_next_cheap_period", start, avg, end)
	}

}

// getNextExpensivePeriodMessage
// Get the next expensive period message.
func getNextExpensivePeriodMessage(periods [][]model.Price, lang language.Tag) string {
	p := message.NewPrinter(lang)

	next, started := utils.GetNextPeriod(periods, time.Now())

	if next == nil {
		return p.Sprintf("alexa_next_expensive_period_nodata")
	}

	avg := utils.FormatPrice(utils.CalculateAverage(next))
	start := utils.FormatTime(next[0].DateTime)
	end := utils.FormatTime(next[len(next)-1].DateTime)

	if started {
		return p.Sprintf("alexa_current_expensive_period", start, avg, end)
	} else {
		return p.Sprintf("alexa_next_expensive_period", start, avg, end)
	}
}

func ParseAlexaSkillRequest(intent model.AlexaIntent, lang language.Tag) model.AlexaSkillResponse {
	p := message.NewPrinter(lang)
	var endSess bool
	var msg string

	switch intent.Name {
	case "AMAZON.CancelIntent":
		endSess = true
		msg = p.Sprintf("alexa_cancel")
	case "AMAZON.HelpIntent":
		msg = p.Sprintf("alexa_help")
	case "AMAZON.StopIntent":
		endSess = true
		msg = p.Sprintf("alexa_stop")
	case "AMAZON.NavigateHomeIntent":
		msg = ""
	case "AMAZON.FallbackIntent":
		msg = p.Sprintf("alexa_welcome")
	case "FULL":
		msg, _ = GetFullFeed(time.Now(), lang)
	case "TODAY_AVERAGE":
		fallthrough
	case "TODAY":
		today := time.Now()
		rating, err := GetDayRating(today)
		avg, err2 := GetDayAverage(today)
		if err != nil || err2 != nil {
			msg = getTodayNoDataMessage(lang)
		} else {
			msg = getTodayRatingMessage(rating, avg, lang)
		}
	case "TOMORROW":
		tomorrow := time.Now().AddDate(0, 0, 1)
		rating, err := GetDayRating(tomorrow)
		avg, err2 := GetDayAverage(tomorrow)
		if err != nil || err2 != nil {
			msg = p.Sprintf("alexa_tomorrow_nodata")
		} else {
			msg = getTomorrowRatingMessage(rating, avg, lang)
		}
	case "NEXT_CHEAP":
		cheapPeriods, err := GetCheapPeriods(time.Now())
		if err != nil {
			msg = p.Sprintf("alexa_next_cheap_period_nodata")
		} else {
			msg = getNextCheapPeriodMessage(cheapPeriods, lang)
		}
	case "NEXT_EXPENSIVE":
		expensivePeriods, err := GetExpensivePeriods(time.Now())
		if err != nil {
			msg = p.Sprintf("alexa_next_expensive_period_nodata")
		} else {
			msg = getNextExpensivePeriodMessage(expensivePeriods, lang)
		}
	case "CURRENT_PRICE":
		price, err := GetPrice(time.Now())
		if err != nil {
			msg = p.Sprintf("alexa_today_nodata")
		} else {
			msg = p.Sprintf("alexa_current_price", utils.FormatPrice(price.Price))
		}

	case "THIRTY_DAY_AVERAGE":
		avg, err := GetThirtyDayAverage(time.Now())
		if err != nil {
			msg = p.Sprintf("alexa_today_nodata")
		} else {
			msg = p.Sprintf("alexa_thirty_day_average", utils.FormatPrice(avg))
		}
	default:
		msg = p.Sprintf("alexa_welcome")
	}

	return utils.WrapAlexaSkillResponse(msg, endSess)
}
