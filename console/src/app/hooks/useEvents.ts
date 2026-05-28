import { useQuery } from '@tanstack/react-query';
import { api, SystemEventResponseApi } from '@api';

export function useEvents() {
  return useQuery<SystemEventResponseApi[]>({
    queryKey: ['events'],
    queryFn: async ({ signal }) => {
      const { data } = await api.events.eventsList({}, { signal });
      return data.items || [];
    },
  });
}
