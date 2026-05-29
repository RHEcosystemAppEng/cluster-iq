import { useState, useEffect } from 'react';
import { api, SystemEventResponseApi } from '@api';

export const useEventsData = () => {
  const [events, setEvents] = useState<SystemEventResponseApi[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchEvents = async () => {
      try {
        setLoading(true);
        setError(null);
        console.log('Fetching recent events...');
        const { data } = await api.events.eventsList({ page: 1, page_size: 10 });
        console.log('Events data received:', data);
        setEvents(data.items || []);
      } catch (err) {
        setError('Failed to fetch events');
        console.error('Error fetching events:', err);
      } finally {
        setLoading(false);
      }
    };

    fetchEvents();
  }, []);

  return { events, loading, error };
};
