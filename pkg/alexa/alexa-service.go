package alexa

import (
	"context"
	"electricity-prices/pkg/date"
	"electricity-prices/pkg/price"
	"fmt"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"strings"
	"time"
)

type Service struct {
	PriceService price.Service
}

func (s *Service) GetTitle(lang language.Tag) string {
	p := message.NewPrinter(lang)
	return p.Sprintf("alexa_full_title")
}

func (s *Service) GetFullFeed(ctx context.Context, t time.Time, lang language.Tag) (string, error) {
	// Get the daily info for the given date
	dailyInfo, err := s.PriceService.GetDailyInfo(ctx, t)

	if err != nil {
		return "", err
	}
	if len(dailyInfo.Prices) == 0 {
		return s.getTodayNoDataMessage(lang), nil
	}
	var messages []string

	// Parse day rating message
	messages = append(messages, s.getTodayRatingMessage(dailyInfo.DayRating, dailyInfo.DayAverage, lang))

	// Get current price
	messages = append(messages, s.getPriceMessage(dailyInfo.Prices, t, lang))

	// Get next cheap period
	messages = append(messages, s.getNextCheapPeriodMessage(dailyInfo.CheapPeriods, t, lang))

	// Get next expensive period
	messages = append(messages, s.getNextExpensivePeriodMessage(dailyInfo.ExpensivePeriods, t, lang))

	// Get tomorrow's data
	tomorrowInfo, err := s.PriceService.GetDailyInfo(ctx, t.AddDate(0, 0, 1))
	if err == nil && len(tomorrowInfo.Prices) > 0 {
		messages = append(messages, s.getTomorrowRatingMessage(tomorrowInfo.DayRating, tomorrowInfo.DayAverage, lang))
	}

	return strings.Join(messages, " "), nil
}

func (s *Service) ProcessAlexaSkillRequest(ctx context.Context, intent AlexaIntent, t time.Time, lang language.Tag) AlexaSkillResponse {
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
	case "AMAZON.NavigateHomeIntent", "AMAZON.FallbackIntent":
		msg = p.Sprintf("alexa_welcome")
	case "FULL":
		feed, err := s.GetFullFeed(ctx, t, lang)
		if err != nil {
			msg = s.getUnknownError(lang)
		}
		if feed != "" {
			msg = feed
		}
	case "TODAY", "TODAY_AVERAGE":
		rating, err := s.PriceService.GetDayRating(ctx, t)
		avg, err2 := s.PriceService.GetDayAverage(ctx, t)
		if err != nil || err2 != nil {
			msg = s.getUnknownError(lang)
		} else {
			msg = s.getTodayRatingMessage(rating, avg, lang)
		}
	case "TOMORROW":
		tomorrow := t.AddDate(0, 0, 1)
		rating, err := s.PriceService.GetDayRating(ctx, tomorrow)
		avg, err2 := s.PriceService.GetDayAverage(ctx, tomorrow)
		if err != nil || err2 != nil {
			msg = s.getTomorrowNoDataMessage(lang)
		} else {
			msg = s.getTomorrowRatingMessage(rating, avg, lang)
		}
	case "NEXT_CHEAP":
		cheapPeriods, err := s.PriceService.GetCheapPeriods(ctx, t)
		if err != nil {
			msg = s.getUnknownError(lang)
		} else {
			msg = s.getNextCheapPeriodMessage(cheapPeriods, t, lang)
		}
	case "NEXT_EXPENSIVE":
		expensivePeriods, err := s.PriceService.GetExpensivePeriods(ctx, t)
		if err != nil {
			msg = s.getUnknownError(lang)
		} else {
			msg = s.getNextExpensivePeriodMessage(expensivePeriods, t, lang)
		}
	case "CURRENT_PRICE":
		pr, err := s.PriceService.GetPrice(ctx, t)
		if err != nil {
			msg = s.getUnknownError(lang)
		} else {
			msg = p.Sprintf("alexa_current_price", price.FormatPrice(pr.Price))
		}

	case "THIRTY_DAY_AVERAGE":
		avg, err := s.PriceService.GetThirtyDayAverage(ctx, t)
		if err != nil {
			msg = s.getUnknownError(lang)
		} else {
			msg = p.Sprintf("alexa_thirty_day_average", price.FormatPrice(avg))
		}
	default:
		msg = p.Sprintf("alexa_welcome")
	}

	return WrapAlexaSkillResponse(msg, endSess)
}

func (s *Service) getUnknownError(lang language.Tag) string {
	p := message.NewPrinter(lang)
	errMesg := p.Sprintf("alexa_unknown_error")

	return errMesg
}
func (s *Service) getTodayNoDataMessage(lang language.Tag) string {
	p := message.NewPrinter(lang)
	noData := p.Sprintf("alexa_today_nodata")

	return noData
}

func (s *Service) getTomorrowNoDataMessage(lang language.Tag) string {
	p := message.NewPrinter(lang)
	noData := p.Sprintf("alexa_tomorrow_nodata")

	return noData
}

func (s *Service) getTodayRatingMessage(dayRating price.DayRating, dayAverage float64, lang language.Tag) string {
	p := message.NewPrinter(lang)

	if dayRating == price.Nil {
		return s.getTodayNoDataMessage(lang)
	}

	rating := p.Sprintf(fmt.Sprintf("alexa_rating_%s", strings.ToLower(string(dayRating))))
	todayRating := p.Sprintf("alexa_today_rating", rating, price.FormatPrice(dayAverage))

	return todayRating
}

func (s *Service) getTomorrowRatingMessage(dayRating price.DayRating, dayAverage float64, lang language.Tag) string {
	p := message.NewPrinter(lang)

	if dayRating == price.Nil {
		return s.getTomorrowNoDataMessage(lang)
	}

	rating := p.Sprintf(fmt.Sprintf("alexa_rating_%s", strings.ToLower(string(dayRating))))
	tomorrowRating := p.Sprintf("alexa_tomorrow_rating", rating, price.FormatPrice(dayAverage))

	return tomorrowRating
}

func (s *Service) getPriceMessage(prices []price.Price, t time.Time, lang language.Tag) string {
	p := message.NewPrinter(lang)

	for _, pr := range prices {
		if date.SameHour(t, pr.DateTime) {
			return p.Sprintf("alexa_current_price", price.FormatPrice(pr.Price))
		}
	}
	return p.Sprintf("alexa_current_price_nodata")
}

// getNextCheapPeriodMessage
// Get the next cheap period message.
func (s *Service) getNextCheapPeriodMessage(periods [][]price.Price, t time.Time, lang language.Tag) string {
	p := message.NewPrinter(lang)

	if len(periods) == 0 {
		return p.Sprintf("alexa_next_cheap_period_nodata")
	}

	next, started := price.GetNextPeriod(periods, t)

	if next == nil {
		return p.Sprintf("alexa_next_cheap_period_none_left")
	}

	avg := price.FormatPrice(price.CalculateAverage(next))
	start := date.FormatTime(next[0].DateTime)
	end := date.FormatTime(next[len(next)-1].DateTime.Add(time.Hour))

	if started {
		return p.Sprintf("alexa_current_cheap_period", start, avg, end)
	} else {
		return p.Sprintf("alexa_next_cheap_period", start, avg, end)
	}

}

// getNextExpensivePeriodMessage
// Get the next expensive period message.
func (s *Service) getNextExpensivePeriodMessage(periods [][]price.Price, t time.Time, lang language.Tag) string {
	p := message.NewPrinter(lang)

	if len(periods) == 0 {
		return p.Sprintf("alexa_next_expensive_period_nodata")
	}

	next, started := price.GetNextPeriod(periods, t)

	if next == nil {
		return p.Sprintf("alexa_next_expensive_period_none_left")
	}

	avg := price.FormatPrice(price.CalculateAverage(next))
	start := date.FormatTime(next[0].DateTime)
	end := date.FormatTime(next[len(next)-1].DateTime.Add(time.Hour))

	if started {
		return p.Sprintf("alexa_current_expensive_period", start, avg, end)
	} else {
		return p.Sprintf("alexa_next_expensive_period", start, avg, end)
	}
}
