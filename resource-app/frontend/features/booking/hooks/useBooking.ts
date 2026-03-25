import { useState, useMemo } from 'react';
import { useBookingContext } from '../context';
import { useUser } from '../../../features/user';
import { Resource } from '../../../features/resource/types';
import { BookingStatus } from '../types';
import { addDays, addMinutes, addHours, isBefore, format, parse, set } from 'date-fns';
import { APP_CONFIG } from '../../../config';

export const useBooking = (resource: Resource, onSuccess: () => void) => {
  const { createBooking, bookings } = useBookingContext();
  const { currentUser } = useUser();
  const [step, setStep] = useState<'time' | 'form'>('time');
  const [isSubmitting, setIsSubmitting] = useState(false);

  // --- State ---
  const initialDate = useMemo(() => {
    const now = new Date();
    if (now.getHours() >= APP_CONFIG.VISUAL_END_HOUR) return addDays(now, 1).toISOString().split('T')[0];
    return now.toISOString().split('T')[0];
  }, []);

  const [date, setDate] = useState(initialDate);
  const [duration, setDuration] = useState(APP_CONFIG.DEFAULT_DURATION_MINUTES);
  const [startTime, setStartTime] = useState<string>(''); // "HH:mm"
  const [formData, setFormData] = useState<Record<string, unknown>>({});

  // --- Logic ---
  const minStartTime = useMemo(() => {
    return addHours(new Date(), resource.minLeadTimeHours);
  }, [resource.minLeadTimeHours]);

  // Get existing bookings for display and conflict check
  const existingBookings = useMemo(() => {
    return bookings.filter(b =>
      b.resourceId === resource.id &&
      b.status !== BookingStatus.CANCELLED &&
      b.status !== BookingStatus.REJECTED &&
      b.start.startsWith(date) // Simple string match for same day
    ).sort((a, b) => new Date(a.start).getTime() - new Date(b.start).getTime());
  }, [bookings, resource.id, date]);

  // Calculate status of selected time
  const timeStatus = useMemo(() => {
    if (!startTime || !date) return { valid: false, message: 'Select a time', conflict: null };

    const baseDate = parse(date, 'yyyy-MM-dd', new Date());
    const startDateTime = parse(startTime, 'HH:mm', baseDate);
    const endDateTime = addMinutes(startDateTime, duration);

    // 1. Past check
    if (isBefore(startDateTime, new Date())) return { valid: false, message: 'Cannot book in the past', conflict: null };

    // 2. Lead time check
    if (isBefore(startDateTime, minStartTime)) return { valid: false, message: `Book ${resource.minLeadTimeHours}h in advance`, conflict: null };

    // 3. Conflict check
    const conflict = existingBookings.find(b => {
      const bStart = new Date(b.start);
      const bEnd = new Date(b.end);
      return (startDateTime < bEnd && endDateTime > bStart);
    });

    if (conflict) {
      const range = `${format(new Date(conflict.start), APP_CONFIG.TIME_FORMAT)} - ${format(new Date(conflict.end), APP_CONFIG.TIME_FORMAT)}`;
      return { valid: false, message: `Overlaps: ${range}`, conflict };
    }

    return { valid: true, message: 'Available', conflict: null };
  }, [date, startTime, duration, existingBookings, minStartTime, resource.minLeadTimeHours]);

  // Visual Timeline Data
  const timelineSegments = useMemo(() => {
    const segments = [];
    const startVisual = APP_CONFIG.VISUAL_START_HOUR;
    const endVisual = APP_CONFIG.VISUAL_END_HOUR;

    const baseDate = parse(date, 'yyyy-MM-dd', new Date());
    const startOfDay = set(baseDate, { hours: startVisual, minutes: 0, seconds: 0, milliseconds: 0 });
    const endOfDay = set(baseDate, { hours: endVisual, minutes: 0, seconds: 0, milliseconds: 0 });
    const totalMinutes = (endOfDay.getTime() - startOfDay.getTime()) / 60000;

    // Existing Bookings
    existingBookings.forEach(b => {
      const bStart = new Date(b.start);
      const bEnd = new Date(b.end);

      // Clip to visual range
      const effectiveStart = bStart < startOfDay ? startOfDay : bStart;
      const effectiveEnd = bEnd > endOfDay ? endOfDay : bEnd;

      if (effectiveEnd > effectiveStart) {
        const startPct = Math.max(0, (effectiveStart.getTime() - startOfDay.getTime()) / 60000 / totalMinutes * 100);
        const widthPct = Math.min(100, (effectiveEnd.getTime() - effectiveStart.getTime()) / 60000 / totalMinutes * 100);
        segments.push({ left: startPct, width: widthPct, type: 'booked' });
      }
    });

    // Current Selection
    if (startTime) {
      const selfStart = parse(startTime, 'HH:mm', baseDate);
      const selfEnd = addMinutes(selfStart, duration);
      const effectiveStart = selfStart < startOfDay ? startOfDay : selfStart;
      const effectiveEnd = selfEnd > endOfDay ? endOfDay : selfEnd;

      if (effectiveEnd > effectiveStart) {
        const startPct = Math.max(0, (effectiveStart.getTime() - startOfDay.getTime()) / 60000 / totalMinutes * 100);
        const widthPct = Math.min(100, (effectiveEnd.getTime() - effectiveStart.getTime()) / 60000 / totalMinutes * 100);
        segments.push({ left: startPct, width: widthPct, type: timeStatus.valid ? 'selection' : 'conflict' });
      }
    }

    return segments;
  }, [existingBookings, date, startTime, duration, timeStatus]);

  const [validationError, setValidationError] = useState<string | null>(null);

  const handleSubmit = async () => {
    if (!timeStatus.valid || isSubmitting || !currentUser) return;
    setValidationError(null);

    const missingFields = resource.formFields
      .filter(field => field.required)
      .filter(field => {
        const val = formData[field.id];
        if (field.type === 'boolean') return val === undefined;
        if (field.type === 'number') return val === undefined || val === null || val === '';
        return !val || (typeof val === 'string' && val.trim() === '');
      });

    if (missingFields.length > 0) {
      setValidationError(`Please fill in all required fields: ${missingFields.map(f => f.label).join(', ')}`);
      return;
    }

    setIsSubmitting(true);

    try {
      const baseDate = parse(date, 'yyyy-MM-dd', new Date());
      const startDateTime = parse(startTime, 'HH:mm', baseDate);
      const startIso = startDateTime.toISOString();
      const endIso = addMinutes(startDateTime, duration).toISOString();

      const res = await createBooking({
        resourceId: resource.id,
        userId: currentUser?.id,
        start: startIso,
        end: endIso,
        details: formData
      });

      if (res.success) {
        onSuccess();
      } else {
        setValidationError(res.error || 'Failed to create booking');
      }
    } finally {
      setIsSubmitting(false);
    }
  };

  return {
    step,
    setStep,
    isSubmitting,
    date,
    setDate,
    duration,
    setDuration,
    startTime,
    setStartTime,
    formData,
    setFormData,
    timeStatus,
    existingBookings,
    timelineSegments,
    handleSubmit,
    validationError,
  };
};
