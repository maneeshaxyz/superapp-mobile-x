export interface PublicHoliday {
  date: string; // yyyy-MM-dd
  localName: string;
  name: string;
  description?: string;
  countryCode: string;
  fixed: boolean;
  global: boolean;
  types: string[];
}
