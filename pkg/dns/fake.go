package dns

import (
	"fmt"
)

type FakeRecords map[string]Record

type FakeDNSService struct {
	ZoneMap   map[string]Zone
	RecordMap map[string]FakeRecords
}

var _ = Service(&FakeDNSService{})

func (f *FakeDNSService) Zones() ([]Zone, error) {
	result := []Zone{}
	for _, value := range f.ZoneMap {
		result = append(result, value)
	}
	return result, nil
}

func (f *FakeDNSService) WriteZone(zone Zone, create bool) error {
	if f.ZoneMap == nil {
		f.ZoneMap = map[string]Zone{}
		f.RecordMap = map[string]FakeRecords{}
	}
	_, exists := f.ZoneMap[zone.Name]
	if exists && create {
		return fmt.Errorf("zone already exists!")
	}
	if !exists && !create {
		return fmt.Errorf("zone doesn't exist!")
	}
	f.ZoneMap[zone.Name] = zone
	f.RecordMap[zone.Name] = FakeRecords{}
	return nil
}

func (f *FakeDNSService) DeleteZone(zone Zone) error {
	if f.ZoneMap == nil {
		return nil
	}
	delete(f.ZoneMap, zone.Name)
	return nil
}

func (f *FakeDNSService) Records(zone Zone) ([]Record, error) {
	result := []Record{}
	for _, value := range f.RecordMap[zone.Name] {
		result = append(result, value)
	}
	return result, nil
}

func (f *FakeDNSService) WriteRecord(zone Zone, oldRecord, record Record) error {
	if _, exists := f.RecordMap[zone.Name]; !exists {
		f.RecordMap[zone.Name] = map[string]Record{}
	}
	_, exists := f.RecordMap[zone.Name][record.RecordName()]
	if oldRecord != nil && !exists {
		return fmt.Errorf("record doesn't exist!")
	}
	if oldRecord == nil && exists {
		return fmt.Errorf("conflict, record exists")
	}
	f.RecordMap[zone.Name][record.RecordName()] = record
	return nil
}

func (f *FakeDNSService) DeleteRecord(zone Zone, record Record) error {
	if _, exists := f.RecordMap[zone.Name]; !exists {
		return fmt.Errorf("zone doesn't exist!")
	}
	delete(f.RecordMap[zone.Name], record.RecordName())
	return nil
}
