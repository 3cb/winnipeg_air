package main

import (
	"github.com/3cb/winnipeg_air/wair"
	flatbuffers "github.com/google/flatbuffers/go"
	// "github.com/3cb/winnipeg_air/wair"
)

func serialize(data []Reading, date string) []byte {
	builder := flatbuffers.NewBuilder(1024)

	readings := []flatbuffers.UOffsetT{}
	for _, v := range data {
		wair.ReadingStartCoordinatesVector(builder, 2)
		for i := len(v.Coordinates) - 1; i >= 0; i-- {
			builder.PrependFloat64(v.Coordinates[i])
		}
		coordinates := builder.EndVector(2)
		observationid := builder.CreateString(v.ObservationID)
		observationtime := builder.CreateString(v.ObservationTime)
		thingid := builder.CreateString(v.ThingID)
		locationname := builder.CreateString(v.LocationName)
		measurementtype := builder.CreateString(v.MeasurementType)
		measurementvalue := builder.CreateString(v.MeasurementValue)
		measurementunit := builder.CreateString(v.MeasurementUnit)

		wair.ReadingStart(builder)
		wair.ReadingAddCoordinates(builder, coordinates)
		wair.ReadingAddObservationid(builder, observationid)
		wair.ReadingAddObservationtime(builder, observationtime)
		wair.ReadingAddThingid(builder, thingid)
		wair.ReadingAddLocationname(builder, locationname)
		wair.ReadingAddMeasurementtype(builder, measurementtype)
		wair.ReadingAddMeasurementvalue(builder, measurementvalue)
		wair.ReadingAddMeasurementunit(builder, measurementunit)
		readings = append(readings, wair.ReadingEnd(builder))
	}
	wair.MessageStartReadingsVector(builder, len(data))
	for i := len(readings) - 1; i >= 0; i-- {
		builder.PrependUOffsetT(readings[i])
	}
	r := builder.EndVector(len(data))

	d := builder.CreateString(date)

	wair.MessageStart(builder)
	wair.MessageAddDate(builder, d)
	wair.MessageAddReadings(builder, r)
	msg := wair.MessageEnd(builder)

	builder.Finish(msg)
	buf := builder.FinishedBytes()

	return buf, nil
}
