package cmd

import (
	"context"
	"sync"
	"time"

	"github.com/zwinslett/strava-cli-go/calculator"
	"github.com/zwinslett/strava-cli-go/format"
	"github.com/zwinslett/strava-cli-go/model"
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
	for _, activity := range activities {
		wg.Add(2)
		go func(activity model.Activity) {
			defer wg.Done()
			detailedActivity, err := client.GetActivityById(ctx, activity.ID)
			if err != nil {
				return
			}
			activitiesCh <- detailedActivity
		}(activity)
		go func(activity model.Activity) {
			defer wg.Done()
			zones, err := client.GetActivityZones(ctx, activity.ID)
			if err != nil {
				return
			}
			zonesCh <- zones
		}(activity)
	}
	wg.Wait()
	close(activitiesCh)
	close(zonesCh)

	for detailedActivity := range activitiesCh {
		detailedActivities = append(detailedActivities, detailedActivity)
	}
	for zones := range zonesCh {
		allZones = append(allZones, zones...)
	}

	return detailedActivities, allZones, nil
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
