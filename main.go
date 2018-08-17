// main
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/tealeg/xlsx"
)

var (
	session_id = os.Getenv("YD_SESSION_ID")
)

func main() {
	tm := time.Now()
	bg := time.Date(tm.Year(), tm.Month(), 1, 0, 0, 0, 0, time.Local) //.AddDate(0, -1, 0)
	fn := time.Date(tm.Year(), tm.Month(), tm.Day(), 0, 0, 0, 0, time.Local)
	fmt.Println(bg, fn)
	s := fmt.Sprintf("%s_%s_%s", bg.Format("2006-01-02-15-04-05"), tm.Format("2006-01-02-15-04-05"), tm.Format("2006-01-02-15-04-05"))
	o, err := DownloadReport(bg, fn, fmt.Sprintf("./data_%s.xlsx", s))
	if err != nil {
		fmt.Println(err)
	}
	if err := xls(o, fmt.Sprintf("./data_details_%s.xlsx", s)); err != nil {
		fmt.Println(err)
	}
}

func DownloadReport(bg, fn time.Time, fname string) ([]OrderDetailsFlat, error) {

	// Get the data
	buf, err := download(fmt.Sprintf("http://youdrive.today/partner/report?from=%v.000Z&to=%v.000Z", bg.Format("2006-01-02T15:04:05"), fn.Format("2006-01-02T15:04:05")))
	if err != nil {
		return nil, err
	}
	err = ioutil.WriteFile(fname, buf, 0644)
	if err != nil {
		return nil, err
	}

	xlFile, err := xlsx.OpenBinary(buf)
	if err != nil {
		return nil, err
	}
	flats := make([]OrderDetailsFlat, 0)
	for _, sheet := range xlFile.Sheets {
		if sheet.Name == "Otchet_po_zakazam" {
			for rnum, row := range sheet.Rows {
				if rnum > 0 {
					for _, cell := range row.Cells {
						text := cell.String()
						parts := strings.Split(text, "=")
						if len(parts) == 2 {
							time.Sleep(1 * time.Second / 2.0)
							buf, err := download(fmt.Sprintf("https://youdrive.today/api/partner/order/%s/details", parts[1]))
							if err != nil {
								return nil, err
							}
							var od OrderDetails
							if err := json.Unmarshal(buf, &od); err != nil {
								return nil, err
							}
							user := ""
							for _, v := range od.Events {
								if v.Details.UserID != "" {
									user = v.Details.UserID
									break
								}
								//								flats = append(flats, OrderDetailsFlat{
								//									CarID:                       od.CarID,
								//									CarNumber:                   od.CarNumber,
								//									CarModel:                    od.CarModel,
								//									CarImg:                      od.CarImg,
								//									CarImgSide:                  od.CarImgSide,
								//									EventId:                     k,
								//									EventState:                  v.State,
								//									EventName:                   v.Name,
								//									EventStatus:                 v.Status,
								//									EventTime:                   v.Time,
								//									EventLat:                    v.Lat,
								//									EventLon:                    v.Lon,
								//									EventCost:                   v.Cost,
								//									EventDuration:               v.Duration,
								//									EventAdminID:                v.AdminID,
								//									EventUserID:                 v.Details.UserID,
								//									CheckUsageTime:              od.Check.UsageTime,
								//									CheckUsageCost:              od.Check.UsageCost,
								//									CheckUsagePrice:             od.Check.UsagePrice,
								//									CheckUsagePriceType:         od.Check.UsagePriceType,
								//									CheckUsageWorkdayTime:       od.Check.UsageWorkdayTime,
								//									CheckUsageWorkdayCost:       od.Check.UsageWorkdayCost,
								//									CheckUsageWorkdayPrice:      od.Check.UsageWorkdayPrice,
								//									CheckUsageWorkdayPriceType:  od.Check.UsageWorkdayPriceType,
								//									CheckUsageWeekendTime:       od.Check.UsageWeekendTime,
								//									CheckUsageWeekendCost:       od.Check.UsageWeekendCost,
								//									CheckUsageWeekendPrice:      od.Check.UsageWeekendPrice,
								//									CheckUsageWeekendPriceType:  od.Check.UsageWeekendPriceType,
								//									CheckChargingTime:           od.Check.ChargingTime,
								//									CheckChargingCost:           od.Check.ChargingCost,
								//									CheckChargingPrice:          od.Check.ChargingPrice,
								//									CheckChargingPriceType:      od.Check.ChargingPriceType,
								//									CheckParkingTime:            od.Check.ParkingTime,
								//									CheckParkingCost:            od.Check.ParkingCost,
								//									CheckParkingPrice:           od.Check.ParkingPrice,
								//									CheckParkingPriceType:       od.Check.ParkingPriceType,
								//									CheckParkingNightTime:       od.Check.ParkingNightTime,
								//									CheckParkingNightCost:       od.Check.ParkingNightCost,
								//									CheckParkingNightPrice:      od.Check.ParkingNightPrice,
								//									CheckParkingNightPriceType:  od.Check.ParkingNightPriceType,
								//									CheckTransferTime:           od.Check.TransferTime,
								//									CheckTransferCost:           od.Check.TransferCost,
								//									CheckTransferPrice:          od.Check.TransferPrice,
								//									CheckTransferPriceType:      od.Check.TransferPriceType,
								//									CheckTransferNightTime:      od.Check.TransferNightTime,
								//									CheckTransferNightCost:      od.Check.TransferNightCost,
								//									CheckTransferNightPrice:     od.Check.TransferNightPrice,
								//									CheckTransferNightPriceType: od.Check.TransferNightPriceType,
								//									CheckWaitingTime:            od.Check.WaitingTime,
								//									CheckWaitingCost:            od.Check.WaitingCost,
								//									CheckWaitingPrice:           od.Check.WaitingPrice,
								//									CheckWaitingPriceType:       od.Check.WaitingPriceType,
								//									CheckBookingTime:            od.Check.BookingTime,
								//									CheckBookingTimeLeft:        od.Check.BookingTimeLeft,
								//									CheckWaitingTimeLeft:        od.Check.WaitingTimeLeft,
								//									CheckFinishCost:             od.Check.FinishCost,
								//									CheckInsuranceIncluded:      od.Check.InsuranceIncluded,
								//									CheckDailyPrice:             od.Check.DailyPrice,
								//									CheckDailyPriceType:         od.Check.DailyPriceType,
								//									CheckDailyCost:              od.Check.DailyCost,
								//									CheckDailyTime:              od.Check.DailyTime,
								//									CheckDailyStatus:            od.Check.DailyStatus,
								//									CheckDiscountPercent:        od.Check.DiscountPercent,
								//									CheckDiscountPrice:          od.Check.DiscountPrice,
								//									CheckTotalCost:              od.Check.TotalCost,
								//									PeriodStart:                 od.Period.Start,
								//									PeriodFinish:                od.Period.Finish,
								//									PathStartLat:                od.Path.Start.Lat,
								//									PathStartLon:                od.Path.Start.Lon,
								//									PathFinishLat:               od.Path.Finish.Lat,
								//									PathFinishLon:               od.Path.Finish.Lon,
								//									CarOwnerID:                  od.CarOwnerID,
								//									Success:                     od.Success,
								//								})
								//fmt.Println(od)
							}
							flats = append(flats, OrderDetailsFlat{
								OrderID:    parts[1],
								CarID:      od.CarID,
								CarNumber:  od.CarNumber,
								CarModel:   od.CarModel,
								CarImg:     od.CarImg,
								CarImgSide: od.CarImgSide,
								//								EventId:                     k,
								//								EventState:                  v.State,
								//								EventName:                   v.Name,
								//								EventStatus:                 v.Status,
								//								EventTime:                   v.Time,
								//								EventLat:                    v.Lat,
								//								EventLon:                    v.Lon,
								//								EventCost:                   v.Cost,
								//								EventDuration:               v.Duration,
								//								EventAdminID:                v.AdminID,
								UserID:                      user,
								CheckUsageTime:              od.Check.UsageTime,
								CheckUsageCost:              od.Check.UsageCost,
								CheckUsagePrice:             od.Check.UsagePrice,
								CheckUsagePriceType:         od.Check.UsagePriceType,
								CheckUsageWorkdayTime:       od.Check.UsageWorkdayTime,
								CheckUsageWorkdayCost:       od.Check.UsageWorkdayCost,
								CheckUsageWorkdayPrice:      od.Check.UsageWorkdayPrice,
								CheckUsageWorkdayPriceType:  od.Check.UsageWorkdayPriceType,
								CheckUsageWeekendTime:       od.Check.UsageWeekendTime,
								CheckUsageWeekendCost:       od.Check.UsageWeekendCost,
								CheckUsageWeekendPrice:      od.Check.UsageWeekendPrice,
								CheckUsageWeekendPriceType:  od.Check.UsageWeekendPriceType,
								CheckChargingTime:           od.Check.ChargingTime,
								CheckChargingCost:           od.Check.ChargingCost,
								CheckChargingPrice:          od.Check.ChargingPrice,
								CheckChargingPriceType:      od.Check.ChargingPriceType,
								CheckParkingTime:            od.Check.ParkingTime,
								CheckParkingCost:            od.Check.ParkingCost,
								CheckParkingPrice:           od.Check.ParkingPrice,
								CheckParkingPriceType:       od.Check.ParkingPriceType,
								CheckParkingNightTime:       od.Check.ParkingNightTime,
								CheckParkingNightCost:       od.Check.ParkingNightCost,
								CheckParkingNightPrice:      od.Check.ParkingNightPrice,
								CheckParkingNightPriceType:  od.Check.ParkingNightPriceType,
								CheckTransferTime:           od.Check.TransferTime,
								CheckTransferCost:           od.Check.TransferCost,
								CheckTransferPrice:          od.Check.TransferPrice,
								CheckTransferPriceType:      od.Check.TransferPriceType,
								CheckTransferNightTime:      od.Check.TransferNightTime,
								CheckTransferNightCost:      od.Check.TransferNightCost,
								CheckTransferNightPrice:     od.Check.TransferNightPrice,
								CheckTransferNightPriceType: od.Check.TransferNightPriceType,
								CheckWaitingTime:            od.Check.WaitingTime,
								CheckWaitingCost:            od.Check.WaitingCost,
								CheckWaitingPrice:           od.Check.WaitingPrice,
								CheckWaitingPriceType:       od.Check.WaitingPriceType,
								CheckBookingTime:            od.Check.BookingTime,
								CheckBookingTimeLeft:        od.Check.BookingTimeLeft,
								CheckWaitingTimeLeft:        od.Check.WaitingTimeLeft,
								CheckFinishCost:             od.Check.FinishCost,
								CheckInsuranceIncluded:      od.Check.InsuranceIncluded,
								CheckDailyPrice:             od.Check.DailyPrice,
								CheckDailyPriceType:         od.Check.DailyPriceType,
								CheckDailyCost:              od.Check.DailyCost,
								CheckDailyTime:              od.Check.DailyTime,
								CheckDailyStatus:            od.Check.DailyStatus,
								CheckDiscountPercent:        od.Check.DiscountPercent,
								CheckDiscountPrice:          od.Check.DiscountPrice,
								CheckTotalCost:              od.Check.TotalCost,
								PeriodStart:                 od.Period.Start,
								PeriodFinish:                od.Period.Finish,
								PathStartLat:                od.Path.Start.Lat,
								PathStartLon:                od.Path.Start.Lon,
								PathFinishLat:               od.Path.Finish.Lat,
								PathFinishLon:               od.Path.Finish.Lon,
								CarOwnerID:                  od.CarOwnerID,
								Success:                     od.Success,
							})

							fmt.Println(od)
						}
						//fmt.Printf("%d - %s\n", cnum, text)
					}
				}
			}
		}
	}
	return flats, nil
}

func download(url string) ([]byte, error) {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.AddCookie(&http.Cookie{Name: "session_id", Value: session_id})

	var client = &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Write the body to file
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func csv(orders []OrderDetailsFlat) error {
	file, err := os.OpenFile("details.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	return gocsv.MarshalFile(&orders, file)
}

func xls(orders []OrderDetailsFlat, fname string) error {
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("details")
	if err != nil {
		return err
	}

	row := sheet.AddRow()
	cell := row.AddCell()
	cell.Value = "order_id"
	cell = row.AddCell()
	cell.Value = "car_id"
	cell = row.AddCell()
	cell.Value = "car_number"
	cell = row.AddCell()
	cell.Value = "car_model"
	cell = row.AddCell()
	cell.Value = "car_img"
	cell = row.AddCell()
	cell.Value = "car_img_side"
	cell = row.AddCell()
	cell.Value = "user_id"
	cell = row.AddCell()
	cell.Value = "check_usage_time, secs"
	cell = row.AddCell()
	cell.Value = "check_usage_cost"
	cell = row.AddCell()
	cell.Value = "check_usage_price"
	cell = row.AddCell()
	cell.Value = "check_usage_price_type"
	cell = row.AddCell()
	cell.Value = "check_usage_workday_time, secs"
	cell = row.AddCell()
	cell.Value = "check_usage_workday_cost"
	cell = row.AddCell()
	cell.Value = "check_usage_workday_price"
	cell = row.AddCell()
	cell.Value = "check_usage_workday_price_type"
	cell = row.AddCell()
	cell.Value = "check_usage_weekend_time, secs"
	cell = row.AddCell()
	cell.Value = "check_usage_weekend_cost"
	cell = row.AddCell()
	cell.Value = "check_usage_weekend_price"
	cell = row.AddCell()
	cell.Value = "check_usage_weekend_price_type"
	cell = row.AddCell()
	cell.Value = "check_charging_time, secs"
	cell = row.AddCell()
	cell.Value = "check_charging_cost"
	cell = row.AddCell()
	cell.Value = "check_charging_price"
	cell = row.AddCell()
	cell.Value = "check_charging_price_type"
	cell = row.AddCell()
	cell.Value = "check_parking_time, secs"
	cell = row.AddCell()
	cell.Value = "check_parking_cost"
	cell = row.AddCell()
	cell.Value = "check_parking_price"
	cell = row.AddCell()
	cell.Value = "check_parking_price_type"
	cell = row.AddCell()
	cell.Value = "check_parking_night_time, secs"
	cell = row.AddCell()
	cell.Value = "check_parking_night_cost"
	cell = row.AddCell()
	cell.Value = "check_parking_night_price"
	cell = row.AddCell()
	cell.Value = "check_parking_night_price_type"
	cell = row.AddCell()
	cell.Value = "check_transfer_time, secs"
	cell = row.AddCell()
	cell.Value = "check_transfer_cost"
	cell = row.AddCell()
	cell.Value = "check_transfer_price"
	cell = row.AddCell()
	cell.Value = "check_transfer_price_type"
	cell = row.AddCell()
	cell.Value = "check_transfer_night_time, secs"
	cell = row.AddCell()
	cell.Value = "check_transfer_night_cost"
	cell = row.AddCell()
	cell.Value = "check_transfer_night_price"
	cell = row.AddCell()
	cell.Value = "check_transfer_night_price_type"
	cell = row.AddCell()
	cell.Value = "check_waiting_time, secs"
	cell = row.AddCell()
	cell.Value = "check_waiting_cost"
	cell = row.AddCell()
	cell.Value = "check_waiting_price"
	cell = row.AddCell()
	cell.Value = "check_waiting_price_type"
	cell = row.AddCell()
	cell.Value = "check_booking_time, secs"
	cell = row.AddCell()
	cell.Value = "check_booking_time_left, secs"
	cell = row.AddCell()
	cell.Value = "check_waiting_time_left, secs"
	cell = row.AddCell()
	cell.Value = "check_finish_cost"
	cell = row.AddCell()
	cell.Value = "check_insurance_included"
	cell = row.AddCell()
	cell.Value = "check_daily_price"
	cell = row.AddCell()
	cell.Value = "check_daily_price_type"
	cell = row.AddCell()
	cell.Value = "check_daily_cost"
	cell = row.AddCell()
	cell.Value = "check_daily_time, secs"
	cell = row.AddCell()
	cell.Value = "check_daily_status"
	cell = row.AddCell()
	cell.Value = "check_discount_percent"
	cell = row.AddCell()
	cell.Value = "check_discount_price"
	cell = row.AddCell()
	cell.Value = "check_total_cost"
	cell = row.AddCell()
	cell.Value = "period_start"
	cell = row.AddCell()
	cell.Value = "period_finish"
	cell = row.AddCell()
	cell.Value = "path_start_lat"
	cell = row.AddCell()
	cell.Value = "path_start_lon"
	cell = row.AddCell()
	cell.Value = "path_finish_lat"
	cell = row.AddCell()
	cell.Value = "path_finish_lon"
	cell = row.AddCell()
	cell.Value = "car_owner_id"
	cell = row.AddCell()
	cell.Value = "success"
	//    cell = row.AddCell()
	//    cell.Value = "I am a cell!"
	for _, v := range orders {
		row = sheet.AddRow()
		cell = row.AddCell()
		cell.Value = v.OrderID
		cell = row.AddCell()
		cell.Value = v.CarID
		cell = row.AddCell()
		cell.Value = v.CarNumber
		cell = row.AddCell()
		cell.Value = v.CarModel
		cell = row.AddCell()
		cell.Value = v.CarImg
		cell = row.AddCell()
		cell.Value = v.CarImgSide
		cell = row.AddCell()
		cell.Value = v.UserID
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%v", v.CheckUsageTime)
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%v", c2r(v.CheckUsageCost))
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%v", c2r(v.CheckUsagePrice))
		cell = row.AddCell()
		cell.Value = v.CheckUsagePriceType
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%v", v.CheckUsageWorkdayTime)
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%v", v.CheckUsageWorkdayCost)
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%v", c2r(v.CheckUsageWorkdayPrice))
		cell = row.AddCell()
		cell.Value = v.CheckUsageWorkdayPriceType
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%v", v.CheckUsageWeekendTime)
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%v", c2r(v.CheckUsageWeekendCost))
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%v", c2r(v.CheckUsageWeekendPrice))
		cell = row.AddCell()
		cell.Value = v.CheckUsageWeekendPriceType
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%v", v.CheckChargingTime)
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%v", c2r(v.CheckChargingCost))
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%v", c2r(v.CheckChargingPrice))
		cell = row.AddCell()
		cell.Value = v.CheckChargingPriceType
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%v", v.CheckParkingTime)
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%v", c2r(v.CheckParkingCost))
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%v", c2r(v.CheckParkingPrice))
		cell = row.AddCell()
		cell.Value = v.CheckParkingPriceType
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%v", v.CheckParkingNightTime)
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%v", c2r(v.CheckParkingNightCost))
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%v", c2r(v.CheckParkingNightPrice))
		cell = row.AddCell()
		cell.Value = v.CheckParkingNightPriceType
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%v", v.CheckTransferTime)
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%v", c2r(v.CheckTransferCost))
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%v", c2r(v.CheckTransferPrice))
		cell = row.AddCell()
		cell.Value = v.CheckTransferPriceType
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%v", v.CheckTransferNightTime)
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%v", c2r(v.CheckTransferNightCost))
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%v", c2r(v.CheckTransferNightPrice))
		cell = row.AddCell()
		cell.Value = v.CheckTransferNightPriceType
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%v", v.CheckWaitingTime)
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%v", c2r(v.CheckWaitingCost))
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%v", c2r(v.CheckWaitingPrice))
		cell = row.AddCell()
		cell.Value = v.CheckWaitingPriceType
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%v", v.CheckBookingTime)
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%v", v.CheckBookingTimeLeft)
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%v", v.CheckWaitingTimeLeft)
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%v", v.CheckFinishCost)
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%v", v.CheckInsuranceIncluded)
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%v", c2r(v.CheckDailyPrice))
		cell = row.AddCell()
		cell.Value = v.CheckDailyPriceType
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%v", c2r(v.CheckDailyCost))
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%v", v.CheckDailyTime)
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%v", v.CheckDailyStatus)
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%v", v.CheckDiscountPercent)
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%v", c2r(v.CheckDiscountPrice))
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%v", c2r(v.CheckTotalCost))
		cell = row.AddCell()
		cell.Value = v.PeriodStart.Format("01/02/2006 15:04:05")
		cell = row.AddCell()
		cell.Value = v.PeriodFinish.Format("01/02/2006 15:04:05")
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%v", v.PathStartLat)
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%v", v.PathStartLon)
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%v", v.PathFinishLat)
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%v", v.PathFinishLon)
		cell = row.AddCell()
		cell.Value = v.CarOwnerID
		cell = row.AddCell()
		cell.Value = fmt.Sprintf("%v", v.Success)
	}

	return file.Save(fname)
}

func c2r(value int) float64 {
	return float64(value) / 100.0
}

type OrderDetails struct {
	CarID      string `json:"car_id"`
	CarNumber  string `json:"car_number"`
	CarModel   string `json:"car_model"`
	CarImg     string `json:"car_img"`
	CarImgSide string `json:"car_img_side"`
	Events     []struct {
		State  string    `json:"state"`
		Name   string    `json:"name"`
		Status string    `json:"status"`
		Time   time.Time `json:"time"`
		Lat    float64   `json:"lat,omitempty"`
		Lon    float64   `json:"lon,omitempty"`
		//IsPassive bool      `json:"is_passive,omitempty"`
		Cost     int     `json:"cost"`
		Duration float64 `json:"duration,omitempty"`
		AdminID  bool    `json:"admin_id"`
		Details  struct {
			UserID string `json:"user_id"`
		} `json:"details,omitempty"`
	} `json:"events"`
	Check struct {
		UsageTime              int    `json:"usage_time"`
		UsageCost              int    `json:"usage_cost"`
		UsagePrice             int    `json:"usage_price"`
		UsagePriceType         string `json:"usage_price_type"`
		UsageWorkdayTime       int    `json:"usage_workday_time"`
		UsageWorkdayCost       int    `json:"usage_workday_cost"`
		UsageWorkdayPrice      int    `json:"usage_workday_price"`
		UsageWorkdayPriceType  string `json:"usage_workday_price_type"`
		UsageWeekendTime       int    `json:"usage_weekend_time"`
		UsageWeekendCost       int    `json:"usage_weekend_cost"`
		UsageWeekendPrice      int    `json:"usage_weekend_price"`
		UsageWeekendPriceType  string `json:"usage_weekend_price_type"`
		ChargingTime           int    `json:"charging_time"`
		ChargingCost           int    `json:"charging_cost"`
		ChargingPrice          int    `json:"charging_price"`
		ChargingPriceType      string `json:"charging_price_type"`
		ParkingTime            int    `json:"parking_time"`
		ParkingCost            int    `json:"parking_cost"`
		ParkingPrice           int    `json:"parking_price"`
		ParkingPriceType       string `json:"parking_price_type"`
		ParkingNightTime       int    `json:"parking_night_time"`
		ParkingNightCost       int    `json:"parking_night_cost"`
		ParkingNightPrice      int    `json:"parking_night_price"`
		ParkingNightPriceType  string `json:"parking_night_price_type"`
		TransferTime           int    `json:"transfer_time"`
		TransferCost           int    `json:"transfer_cost"`
		TransferPrice          int    `json:"transfer_price"`
		TransferPriceType      string `json:"transfer_price_type"`
		TransferNightTime      int    `json:"transfer_night_time"`
		TransferNightCost      int    `json:"transfer_night_cost"`
		TransferNightPrice     int    `json:"transfer_night_price"`
		TransferNightPriceType string `json:"transfer_night_price_type"`
		WaitingTime            int    `json:"waiting_time"`
		WaitingCost            int    `json:"waiting_cost"`
		WaitingPrice           int    `json:"waiting_price"`
		WaitingPriceType       string `json:"waiting_price_type"`
		BookingTime            int    `json:"booking_time"`
		BookingTimeLeft        int    `json:"booking_time_left"`
		WaitingTimeLeft        int    `json:"waiting_time_left"`
		FinishCost             int    `json:"finish_cost"`
		InsuranceIncluded      bool   `json:"insurance_included"`
		DailyPrice             int    `json:"daily_price"`
		DailyPriceType         string `json:"daily_price_type"`
		DailyCost              int    `json:"daily_cost"`
		DailyTime              int    `json:"daily_time"`
		DailyStatus            bool   `json:"daily_status"`
		DiscountPercent        int    `json:"discount_percent"`
		DiscountPrice          int    `json:"discount_price"`
		TotalCost              int    `json:"total_cost"`
	} `json:"check"`
	Period struct {
		Start  time.Time `json:"start"`
		Finish time.Time `json:"finish"`
	} `json:"period"`
	Path struct {
		Start struct {
			Lat float64 `json:"lat"`
			Lon float64 `json:"lon"`
		} `json:"start"`
		Finish struct {
			Lat float64 `json:"lat"`
			Lon float64 `json:"lon"`
		} `json:"finish"`
	} `json:"path"`
	CarOwnerID string `json:"car_owner_id"`
	Success    bool   `json:"success"`
}

type OrderDetailsFlat struct {
	OrderID    string `csv:"order_id"`
	CarID      string `csv:"car_id"`
	CarNumber  string `csv:"car_number"`
	CarModel   string `csv:"car_model"`
	CarImg     string `csv:"car_img"`
	CarImgSide string `csv:"car_img_side"`
	//	EventId                     int       `csv:"event_id"`
	//	EventState                  string    `csv:"event_state"`
	//	EventName                   string    `csv:"event_name"`
	//	EventStatus                 string    `csv:"event_status"`
	//	EventTime                   time.Time `csv:"event_time"`
	//	EventLat                    float64   `csv:"event_at,omitempty"`
	//	EventLon                    float64   `csv:"event_lon,omitempty"`
	//	EventCost                   int       `csv:"event_cost"`
	//	EventDuration               float64   `csv:"event_duration,omitempty"`
	//	EventAdminID                bool      `csv:"event_admin_id"`
	//	EventUserID                 string    `csv:"event_user_id"`
	UserID                      string    `csv:"user_id"`
	CheckUsageTime              int       `csv:"check_usage_time"`
	CheckUsageCost              int       `csv:"check_usage_cost"`
	CheckUsagePrice             int       `csv:"check_usage_price"`
	CheckUsagePriceType         string    `csv:"check_usage_price_type"`
	CheckUsageWorkdayTime       int       `csv:"check_usage_workday_time"`
	CheckUsageWorkdayCost       int       `csv:"check_usage_workday_cost"`
	CheckUsageWorkdayPrice      int       `csv:"check_usage_workday_price"`
	CheckUsageWorkdayPriceType  string    `csv:"check_usage_workday_price_type"`
	CheckUsageWeekendTime       int       `csv:"check_usage_weekend_time"`
	CheckUsageWeekendCost       int       `csv:"check_usage_weekend_cost"`
	CheckUsageWeekendPrice      int       `csv:"check_usage_weekend_price"`
	CheckUsageWeekendPriceType  string    `csv:"check_usage_weekend_price_type"`
	CheckChargingTime           int       `csv:"check_charging_time"`
	CheckChargingCost           int       `csv:"check_charging_cost"`
	CheckChargingPrice          int       `csv:"check_charging_price"`
	CheckChargingPriceType      string    `csv:"check_charging_price_type"`
	CheckParkingTime            int       `csv:"check_parking_time"`
	CheckParkingCost            int       `csv:"check_parking_cost"`
	CheckParkingPrice           int       `csv:"check_parking_price"`
	CheckParkingPriceType       string    `csv:"check_parking_price_type"`
	CheckParkingNightTime       int       `csv:"check_parking_night_time"`
	CheckParkingNightCost       int       `csv:"check_parking_night_cost"`
	CheckParkingNightPrice      int       `csv:"check_parking_night_price"`
	CheckParkingNightPriceType  string    `csv:"check_parking_night_price_type"`
	CheckTransferTime           int       `csv:"check_transfer_time"`
	CheckTransferCost           int       `csv:"check_transfer_cost"`
	CheckTransferPrice          int       `csv:"check_transfer_price"`
	CheckTransferPriceType      string    `csv:"check_transfer_price_type"`
	CheckTransferNightTime      int       `csv:"check_transfer_night_time"`
	CheckTransferNightCost      int       `csv:"check_transfer_night_cost"`
	CheckTransferNightPrice     int       `csv:"check_transfer_night_price"`
	CheckTransferNightPriceType string    `csv:"check_transfer_night_price_type"`
	CheckWaitingTime            int       `csv:"check_waiting_time"`
	CheckWaitingCost            int       `csv:"check_waiting_cost"`
	CheckWaitingPrice           int       `csv:"check_waiting_price"`
	CheckWaitingPriceType       string    `csv:"check_waiting_price_type"`
	CheckBookingTime            int       `csv:"check_booking_time"`
	CheckBookingTimeLeft        int       `csv:"check_booking_time_left"`
	CheckWaitingTimeLeft        int       `csv:"check_waiting_time_left"`
	CheckFinishCost             int       `csv:"check_finish_cost"`
	CheckInsuranceIncluded      bool      `csv:"check_insurance_included"`
	CheckDailyPrice             int       `csv:"check_daily_price"`
	CheckDailyPriceType         string    `csv:"check_daily_price_type"`
	CheckDailyCost              int       `csv:"check_daily_cost"`
	CheckDailyTime              int       `csv:"check_daily_time"`
	CheckDailyStatus            bool      `csv:"check_daily_status"`
	CheckDiscountPercent        int       `csv:"check_discount_percent"`
	CheckDiscountPrice          int       `csv:"check_discount_price"`
	CheckTotalCost              int       `csv:"check_total_cost"`
	PeriodStart                 time.Time `csv:"period_start"`
	PeriodFinish                time.Time `csv:"period_finish"`
	PathStartLat                float64   `csv:"path_start_lat"`
	PathStartLon                float64   `csv:"path_start_lon"`
	PathFinishLat               float64   `csv:"path_finish_lat"`
	PathFinishLon               float64   `csv:"path_finish_lon"`
	CarOwnerID                  string    `csv:"car_owner_id"`
	Success                     bool      `csv:"success"`
}
