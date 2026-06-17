package cmd

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/zwinslett/strava-cli-go/calculator"
	"github.com/zwinslett/strava-cli-go/format"
	"github.com/zwinslett/strava-cli-go/model"
	"github.com/zwinslett/strava-cli-go/telegram"
)

type DateRange = string

const (
	Monthly DateRange = "Monthly"
	Weekly  DateRange = "Weekly"
)

func fetchLastActivity(ctx context.Context) (model.DetailedActivity, []model.Zones, error) {
	var detailedActivity model.DetailedActivity
	var zones []model.Zones
	var detailedActivityErr, zonesErr error
	var wg sync.WaitGroup
	wg.Add(2)

	activity, err := client.GetRecentActivities(ctx, 10)
	if err != nil {
		return model.DetailedActivity{}, nil, err
	}
	// Filter out non-running activities
	activity = calculator.FilterByType("Run", activity)

	go func() {
		defer wg.Done()
		detailedActivity, detailedActivityErr = client.GetActivityById(ctx, activity[0].ID)
	}()
	go func() {
		defer wg.Done()
		zones, zonesErr = client.GetActivityZones(ctx, activity[0].ID)
	}()
	wg.Wait()
	if detailedActivityErr != nil {
		return model.DetailedActivity{}, nil, detailedActivityErr
	}
	if zonesErr != nil {
		return model.DetailedActivity{}, nil, zonesErr
	}
	return detailedActivity, zones, nil
}

func fetchStats(ctx context.Context, before int64, after int64) ([]model.DetailedActivity, []model.Zones, error) {
	var detailedActivities []model.DetailedActivity
	var allZones []model.Zones
	var wg sync.WaitGroup

	activities, err := client.GetActivitiesByRange(ctx, after, before)
	if err != nil {
		return []model.DetailedActivity{}, nil, err
	}
	// Filter out non-running activities
	activities = calculator.FilterByType("Run", activities)
	activitiesCh := make(chan model.DetailedActivity, len(activities))
	zonesCh := make(chan []model.Zones, len(activities))
	errCh := make(chan error, len(activities)*2)
	for _, activity := range activities {
		wg.Add(2)
		go func(activity model.Activity) {
			defer wg.Done()
			detailedActivity, err := client.GetActivityById(ctx, activity.ID)
			if err != nil {
				errCh <- err
				return
			}
			activitiesCh <- detailedActivity
		}(activity)
		go func(activity model.Activity) {
			defer wg.Done()
			zones, err := client.GetActivityZones(ctx, activity.ID)
			if err != nil {
				errCh <- err
				return
			}
			zonesCh <- zones
		}(activity)
	}
	wg.Wait()
	close(activitiesCh)
	close(zonesCh)
	close(errCh)

	for err := range errCh {
		log.Println(err)
	}

	for detailedActivity := range activitiesCh {
		detailedActivities = append(detailedActivities, detailedActivity)
	}
	for zones := range zonesCh {
		allZones = append(allZones, zones...)
	}

	return detailedActivities, allZones, nil
}

func fetchComparativeStats(ctx context.Context, previousBefore int64, previousAfter int64, currentBefore int64, currentAfter int64) ([]model.DetailedActivity, []model.DetailedActivity, []model.Zones, []model.Zones, error) {
	var previousActivities []model.DetailedActivity
	var currentActivities []model.DetailedActivity
	var previousZones []model.Zones
	var currentZones []model.Zones
	var previousErr, currentErr error
	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		previousActivities, previousZones, previousErr = fetchStats(ctx, previousBefore, previousAfter)
		if previousErr != nil {
			return
		}
	}()

	go func() {
		defer wg.Done()
		currentActivities, currentZones, currentErr = fetchStats(ctx, currentBefore, currentAfter)
		if currentErr != nil {
			return
		}
	}()
	wg.Wait()

	if previousErr != nil {
		return nil, nil, nil, nil, previousErr
	}
	if currentErr != nil {
		return nil, nil, nil, nil, currentErr
	}
	return previousActivities, currentActivities, previousZones, currentZones, nil
}

func statsMessageBuilder(ctx context.Context, dateRange DateRange) error {
	var detailedActivities []model.DetailedActivity
	var allZones []model.Zones
	var fetchErr error
	before := time.Now().Unix()
	if dateRange == Weekly {
		detailedActivities, allZones, fetchErr = fetchStats(ctx, before, time.Now().AddDate(0, 0, -7).Unix())
	} else if dateRange == Monthly {
		detailedActivities, allZones, fetchErr = fetchStats(ctx, before, time.Now().AddDate(0, -1, 0).Unix())
	}

	if fetchErr != nil {
		return fetchErr
	}
	err := bot.SendMessage(ctx, format.ActivitiesMessage(detailedActivities, dateRange)+"\n\n"+format.ZonesMessage(allZones, calculator.Heartrate)+"\n\n"+format.GearMessage(detailedActivities))
	if err != nil {
		return err
	}
	return nil
}

func statsComparisonMessageBuilder(ctx context.Context, dateRange DateRange) error {
	var previousActivities, currentActivities []model.DetailedActivity
	var previousZones, currentZones []model.Zones
	var fetchErr error
	before := time.Now().Unix()
	if dateRange == Weekly {
		previousActivities, currentActivities, previousZones, currentZones, fetchErr = fetchComparativeStats(ctx, time.Now().AddDate(0, 0, -7).Unix(), time.Now().AddDate(0, 0, -14).Unix(), before, time.Now().AddDate(0, 0, -7).Unix())
	} else if dateRange == Monthly {
		previousActivities, currentActivities, previousZones, currentZones, fetchErr = fetchComparativeStats(ctx, time.Now().AddDate(0, -1, 0).Unix(), time.Now().AddDate(0, -2, 0).Unix(), before, time.Now().AddDate(0, -1, 0).Unix())
	}
	if fetchErr != nil {
		return fetchErr
	}
	err := bot.SendMessage(ctx, format.ActivitiesComparisonMessage(previousActivities, currentActivities, dateRange)+"\n\n"+format.ZonesComparisonMessage(previousZones, currentZones, calculator.Heartrate)+"\n\n"+format.GearMessage(currentActivities))
	if err != nil {
		return err
	}
	return nil
}

func activityMessageBuilder(ctx context.Context) error {
	detailedActivity, allZones, err := fetchLastActivity(ctx)
	if err != nil {
		return err
	}
	err = bot.SendMessage(ctx, format.ActivityMessage(detailedActivity)+"\n\n"+format.SplitMessage(detailedActivity)+"\n\n"+format.ZonesMessage(allZones, calculator.Heartrate))
	if err != nil {
		return err
	}
	return nil
}

func handleUpdates(ctx context.Context, result telegram.Result) error {
	switch result.Message.Text {
	case telegram.CmdLatest:
		return activityMessageBuilder(ctx)
	case telegram.CmdWeekly:
		return statsMessageBuilder(ctx, Weekly)
	case telegram.CmdMonthly:
		return statsMessageBuilder(ctx, Monthly)
	default:
		return bot.SendMessage(ctx, "Unsupported Command")
	}
}

func pollForUpdates(ctx context.Context) {
	offset := 0
	for {
		results, err := bot.GetUpdates(ctx, 10, 60, offset)
		if err != nil {
			log.Println(err)
			continue
		}
		for _, result := range results {
			err := handleUpdates(ctx, result)
			if err != nil {
				log.Println(err)
				continue
			}
			offset = int(result.UpdateID) + 1
		}
	}
}
