import { useQuery } from '@tanstack/react-query';
import { api, SystemEventResponseApi } from '@api';

export const useEventsData = () => {
  const { data, isLoading, error } = useQuery<SystemEventResponseApi[]>({
    queryKey: ['recentEvents'],
    queryFn: async ({ signal }) => {
      const { data } = await api.events.eventsList({ page: 1, page_size: 10 }, { signal });
      return data.items || [];
    },
    refetchInterval: 5_000,
  });

  return {
    events: data || [],
    loading: isLoading,
    error: error ? 'Failed to fetch events' : null,
  };
};
