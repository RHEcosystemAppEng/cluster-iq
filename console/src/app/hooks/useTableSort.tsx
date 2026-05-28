import React from 'react';
import { ThProps } from '@patternfly/react-table';

export function useTableSort<T>(
  filteredData: T[],
  getSortableRowValues: (item: T) => (string | number | null)[],
  defaultSortIndex: number = 0, // Default to column 0
  defaultSortDirection: 'asc' | 'desc' = 'asc' // Default to ascending order
) {
  const [activeSortIndex, setActiveSortIndex] = React.useState<number | undefined>(
    typeof defaultSortIndex === 'number' && defaultSortIndex !== null ? defaultSortIndex : 0
  );
  const [activeSortDirection, setActiveSortDirection] = React.useState<'asc' | 'desc' | undefined>(
    defaultSortDirection
  );

  let sortedData = filteredData;
  if (typeof activeSortIndex === 'number' && activeSortIndex !== null) {
    sortedData = filteredData.sort((a, b) => {
      const aValue = getSortableRowValues(a)[activeSortIndex];
      const bValue = getSortableRowValues(b)[activeSortIndex];

      if (typeof aValue === 'number') {
        return activeSortDirection === 'asc'
          ? (aValue as number) - (bValue as number)
          : (bValue as number) - (aValue as number);
      } else {
        return activeSortDirection === 'asc'
          ? (aValue as string).localeCompare(bValue as string)
          : (bValue as string).localeCompare(aValue as string);
      }
    });
  }

  const getSortParams = (columnIndex: number): ThProps['sort'] => ({
    sortBy: {
      index: activeSortIndex,
      direction: activeSortDirection,
      defaultDirection: 'asc' as const,
    },
    onSort: (_event: unknown, index: number, direction: 'asc' | 'desc') => {
      setActiveSortIndex(index);
      setActiveSortDirection(direction);
    },
    columnIndex,
  });

  return { sortedData, getSortParams };
}
