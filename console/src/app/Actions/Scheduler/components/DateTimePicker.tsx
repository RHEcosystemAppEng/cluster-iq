import { InputGroup, InputGroupItem, DatePicker, TimePicker, Button, Popover } from '@patternfly/react-core';
import { GlobeAmericasIcon } from '@patternfly/react-icons';
import React, { useState } from 'react';
import { DateTime } from 'luxon';

interface DateTimePickerProps {
  onChange: (formattedDateTime: string) => void;
  width?: string;
}

// Partil list of timezone
const timeZones = [
  // Universal
  'UTC',

  // Americas
  'America/New_York',
  'America/Chicago',
  'America/Denver',
  'America/Los_Angeles',
  'America/Toronto',
  'America/Vancouver',
  'America/Mexico_City',
  'America/Bogota',
  'America/Sao_Paulo',
  'America/Santiago',
  'America/Buenos_Aires',

  // Europe
  'Europe/London',
  'Europe/Dublin',
  'Europe/Paris',
  'Europe/Berlin',
  'Europe/Madrid',
  'Europe/Rome',
  'Europe/Amsterdam',
  'Europe/Brussels',
  'Europe/Zurich',
  'Europe/Stockholm',
  'Europe/Vienna',
  'Europe/Warsaw',
  'Europe/Athens',
  'Europe/Istanbul',
  'Europe/Moscow',
  'Europe/Kiev',

  // Middle East
  'Asia/Jerusalem',
  'Asia/Amman',
  'Asia/Beirut',
  'Asia/Riyadh',
  'Asia/Dubai',
  'Asia/Baghdad',
  'Asia/Tehran',

  // Asia & Pacific
  'Asia/Kolkata',
  'Asia/Karachi',
  'Asia/Bangkok',
  'Asia/Jakarta',
  'Asia/Shanghai',
  'Asia/Hong_Kong',
  'Asia/Taipei',
  'Asia/Seoul',
  'Asia/Tokyo',
  'Asia/Singapore',
  'Asia/Manila',

  // Oceania
  'Australia/Sydney',
  'Australia/Melbourne',
  'Australia/Brisbane',
  'Australia/Perth',
  'Pacific/Auckland',

  // Africa
  'Africa/Cairo',
  'Africa/Lagos',
  'Africa/Nairobi',
  'Africa/Johannesburg',
];

const DateTimePicker: React.FunctionComponent<DateTimePickerProps> = ({ onChange, width = '300px' }) => {
  const [isTimeZoneOpen, setIsTimeZoneOpen] = useState(false);
  const [selectedDate, setSelectedDate] = useState<Date | undefined>(undefined);
  const [selectedTimeZone, setSelectedTimeZone] = useState('UTC');

  // Handle date selection
  const onDateChange = (_event: React.FormEvent<HTMLInputElement>, _value: string, date?: Date) => {
    if (date) {
      setSelectedDate(date);
      updateDateTime(date, selectedTimeZone);
    }
  };

  // Handle time selection
  const onTimeChange = (_event: React.FormEvent<HTMLInputElement>, _time: string, hour?: number, minute?: number) => {
    if (selectedDate && hour !== undefined && minute !== undefined) {
      const updatedDate = new Date(selectedDate);
      updatedDate.setHours(hour, minute);
      setSelectedDate(updatedDate);
      updateDateTime(updatedDate, selectedTimeZone);
    }
  };

  // Handle timezone selection
  const onSelectTimeZone = (event: React.MouseEvent<HTMLDivElement>) => {
    const newTimeZone = event.currentTarget.textContent as string;
    setSelectedTimeZone(newTimeZone);
    setIsTimeZoneOpen(false);

    if (selectedDate) {
      updateDateTime(selectedDate, newTimeZone);
    }
  };

  // Convert to ISO 8601 Format with offset
  const updateDateTime = (date: Date, timeZone: string) => {
    const zonedDateTime = DateTime.fromJSDate(date, { zone: timeZone }).toISO();
    if (zonedDateTime) onChange(zonedDateTime);
  };

  return (
    <div style={{ width }}>
      <InputGroup>
        <InputGroupItem>
          <DatePicker
            onChange={onDateChange}
            aria-label="Date picker"
            placeholder="YYYY-MM-DD"
            popoverProps={{
              position: 'bottom',
              zIndex: 9999,
            }}
            appendTo={document.body}
          />
        </InputGroupItem>
        <InputGroupItem>
          <TimePicker
            aria-label="Time picker"
            is24Hour
            onChange={onTimeChange}
            menuAppendTo={document.body}
            style={{ width: '140px', zIndex: 9999 }}
          />
        </InputGroupItem>
        <InputGroupItem>
          <Popover
            position="bottom"
            bodyContent={
              <div className="pf-v6-u-p-sm" style={{ maxHeight: '200px', overflowY: 'auto' }}>
                {timeZones.map(zone => (
                  <div key={zone} onClick={onSelectTimeZone} className="pf-v6-u-p-sm" style={{ cursor: 'pointer' }}>
                    {zone}
                  </div>
                ))}
              </div>
            }
            showClose={false}
            isVisible={isTimeZoneOpen}
            hasNoPadding
            hasAutoWidth
            appendTo={document.body}
            zIndex={9999}
          >
            <Button
              icon={<GlobeAmericasIcon />}
              variant="control"
              aria-label="Toggle timezone"
              onClick={() => setIsTimeZoneOpen(!isTimeZoneOpen)}
            >
              {selectedTimeZone}
            </Button>
          </Popover>
        </InputGroupItem>
      </InputGroup>
    </div>
  );
};

export default DateTimePicker;
