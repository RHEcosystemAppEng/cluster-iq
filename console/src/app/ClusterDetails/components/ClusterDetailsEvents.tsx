import { LoadingSpinner } from '@app/components/common/LoadingSpinner';
import { ResultStatus } from '@app/types/types';
import { api, SystemEventResponseApi } from '@api';
import { ThProps, Table, Thead, Tr, Th, Tbody, Td } from '@patternfly/react-table';
import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { getResultIcon } from '@app/utils/renderUtils';
import { useTableSort } from '@app/hooks/useTableSort.tsx';
import { EmptyState } from '@patternfly/react-core';
import { SearchIcon } from '@patternfly/react-icons';

interface TableEventsProps {
  data: SystemEventResponseApi[];
  getSortParams: (columnIndex: number) => ThProps['sort'];
}

const columnNames = {
  action: 'Action',
  result: 'Result',
  severity: 'Severity',
  loggedBy: 'Logged by',
  description: 'Description',
  date: 'Date',
};

export const EmptyStateNoFound: React.FunctionComponent = () => (
  <EmptyState headingLevel="h4" icon={SearchIcon} titleText="No events"></EmptyState>
);

const TableEvents: React.FunctionComponent<TableEventsProps> = ({ data, getSortParams }) => {
  return (
    <Table aria-label="Events table">
      <Thead>
        <Tr>
          <Th sort={getSortParams(0)}>{columnNames.action}</Th>
          <Th sort={getSortParams(1)}>{columnNames.result}</Th>
          <Th sort={getSortParams(2)}>{columnNames.severity}</Th>
          <Th sort={getSortParams(3)}>{columnNames.loggedBy}</Th>
          <Th>{columnNames.description}</Th>
          <Th sort={getSortParams(5)}>{columnNames.date}</Th>
        </Tr>
      </Thead>
      <Tbody>
        {data.map(event => (
          <Tr key={event.id}>
            <Td>{event.action}</Td>
            <Td>
              {getResultIcon(event.result as ResultStatus)} {event.result}
            </Td>
            <Td>{event.severity}</Td>
            <Td>{event.triggeredBy}</Td>
            <Td>{event.description}</Td>
            <Td>{event.timestamp}</Td>
          </Tr>
        ))}
      </Tbody>
    </Table>
  );
};

export const ClusterDetailsEvents: React.FunctionComponent = () => {
  const [data, setData] = useState<SystemEventResponseApi[] | []>([]);
  const [loading, setLoading] = useState(true);
  const { clusterID } = useParams();

  useEffect(() => {
    if (!clusterID) {
      setLoading(false);
      return;
    }

    let cancelled = false;
    const fetchData = async () => {
      try {
        const { data: clusterEvents } = await api.clusters.eventsList(clusterID);
        if (!cancelled) {
          // TODO. Move to debug
          console.log('Fetched events:', clusterEvents);
          setData(clusterEvents.items || []);
        }
      } catch (error) {
        if (!cancelled) {
          // TODO. Move to debug
          console.error('Error fetching events:', error);
        }
      } finally {
        if (!cancelled) {
          setLoading(false);
        }
      }
    };

    fetchData();

    return () => {
      cancelled = true;
    };
  }, [clusterID]);

  // TODO. Move to debug
  console.log('Rendered events data:', data);

  const getSortableRowValues = (event: SystemEventResponseApi): (string | number | null)[] => {
    const { action, result, severity, triggeredBy, description: description, timestamp } = event;
    return [action, result, severity, triggeredBy, description ?? null, timestamp];
  };

  const { sortedData, getSortParams } = useTableSort<SystemEventResponseApi>(data, getSortableRowValues, 5, 'desc');
  if (loading) return <LoadingSpinner />;
  if (sortedData.length === 0) return <EmptyStateNoFound />;
  return <TableEvents data={sortedData} getSortParams={getSortParams} />;
};
