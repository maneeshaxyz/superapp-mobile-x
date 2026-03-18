import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { PublicHoliday } from './types';
import { holidayService } from './service';

interface HolidayContextType {
  holidays: PublicHoliday[];
  isLoading: boolean;
  error: string | null;
  refreshHolidays: () => Promise<void>;
}

const HolidayContext = createContext<HolidayContextType | undefined>(undefined);

export const HolidayProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [holidays, setHolidays] = useState<PublicHoliday[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchHolidays = async () => {
    setIsLoading(true);
    setError(null);
    try {
      const currentYear = new Date().getFullYear();
      const currentYearHolidays = await holidayService.getHolidays(currentYear);
      const nextYearHolidays = await holidayService.getHolidays(currentYear + 1);
      
      setHolidays([...currentYearHolidays, ...nextYearHolidays]);
    } catch (err) {
      console.error("Failed to load holidays", err);
      setError("Failed to load holiday data");
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchHolidays();
  }, []);

  return (
    <HolidayContext.Provider value={{
      holidays,
      isLoading,
      error,
      refreshHolidays: fetchHolidays
    }}>
      {children}
    </HolidayContext.Provider>
  );
};

export const useHoliday = () => {
  const context = useContext(HolidayContext);
  if (context === undefined) {
    throw new Error('useHoliday must be used within a HolidayProvider');
  }
  return context;
};
