import { TableSkeleton } from '@app/components/common/TableSkeleton';
import { ResultStatus } from '@app/types/types';
import { api, SystemEventResponseApi } from '@api';
import { ThProps, Table, Thead, Tr, Th, Tbody, Td } from '@patternfly/react-table';
import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { getResultIcon } from '@app/utils/renderUtils';
import { useTableSort } from '@app/hooks/useTableSort.tsx';
import { EmptyState } from '@patternfly/react-core';
import { SearchIcon } from '@patternfly/react-icons';
import { debug } from '@app/utils/debugLogs';

interface TableEventsProps {
  data: SystemEventResponseApi[];
  getSortParams: (columnIndex: number) => ThProps['sort'];
}

const columnNames = {
  action: 'Action',
  result: 'Result',
  severity: 'Severity',
  requester: 'Requester',
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
          <Th sort={getSortParams(3)}>{columnNames.requester}</Th>
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
            <Td>{event.requester}</Td>
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
          debug('Fetched events:', clusterEvents);
          setData(clusterEvents.items || []);
        }
      } catch (error) {
        if (!cancelled) {
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

  debug('Rendered events data:', data);

  const getSortableRowValues = (event: SystemEventResponseApi): (string | number | null)[] => {
    const { action, result, severity, requester, description: description, timestamp } = event;
    return [action, result, severity, requester, description ?? null, timestamp];
  };

  const { sortedData, getSortParams } = useTableSort<SystemEventResponseApi>(data, getSortableRowValues, 5, 'desc');
  if (loading) return <TableSkeleton columns={6} />;
  if (sortedData.length === 0) return <EmptyStateNoFound />;
  return <TableEvents data={sortedData} getSortParams={getSortParams} />;
};
