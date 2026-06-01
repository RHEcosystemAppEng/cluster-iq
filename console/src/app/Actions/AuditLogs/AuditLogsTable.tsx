import { LoadingSpinner } from '@app/components/common/LoadingSpinner';
import { ActionOperations, ResultStatus } from '@app/types/types';
import { SystemEventResponseApi } from '@api';
import { Table, Tbody, Td, Th, Thead, Tr } from '@patternfly/react-table';
import React, { useMemo } from 'react';
import {
  renderOperationLabel,
  renderActionStatusLabel,
  renderResourceBadge,
  ResourceBadge,
} from '@app/utils/renderUtils';
import { parseScanTimestamp, resolveResourcePath } from '@app/utils/parseFuncs';
import { useTableSort } from '@app/hooks/useTableSort.tsx';
import { EmptyState } from '@patternfly/react-core';
import { TablePagination } from '@app/components/common/TablesPagination';
import { SearchIcon } from '@patternfly/react-icons';
import { AuditLogsTableProps } from './types';
import { Link } from 'react-router-dom';
import { useEvents } from '@app/hooks/useEvents';
import { useTablePagination } from '@app/hooks/useTablePagination';

const columnNames = {
  scheduledAction: 'Action',
  operation: 'Operation',
  resource: 'Resource',
  description: 'Description',
  status: 'Status',
  requester: 'Requester',
  date: 'Date',
};

const EmptyStateNoFound: React.FunctionComponent = () => (
  <EmptyState headingLevel="h4" icon={SearchIcon} titleText="No events"></EmptyState>
);

export const AuditLogsTable: React.FunctionComponent<AuditLogsTableProps> = ({
  accountName,
  action,
  provider,
  result,
  requester,
}) => {
  const { data: allEvents = [], isLoading } = useEvents();

  const filtered = useMemo(() => {
    let filteredResult = allEvents;

    if (accountName) {
      filteredResult = filteredResult.filter(event =>
        event.accountId?.toLowerCase().includes(accountName.toLowerCase())
      );
    }

    if (action?.length) {
      filteredResult = filteredResult.filter(event => action.includes(event.action as ActionOperations));
    }

    if (provider?.length) {
      filteredResult = filteredResult.filter(event => event.provider && provider.some(p => p === event.provider));
    }

    if (result?.length) {
      filteredResult = filteredResult.filter(event => result.includes(event.result as ResultStatus));
    }

    if (requester) {
      filteredResult = filteredResult.filter(event => event.requester?.toLowerCase().includes(requester.toLowerCase()));
    }

    return filteredResult;
  }, [allEvents, accountName, action, provider, result, requester]);

  const { page, perPage, setPage, setPerPage, paginatedData, totalItems } = useTablePagination({
    data: filtered,
    filterDeps: [accountName, action, provider, result, requester],
  });

  const getSortableRowValues = (event: SystemEventResponseApi): (string | number | null)[] => {
    const { timestamp, action, resourceId, result, requester, scheduleId, description } = event;
    return [
      timestamp ?? null,
      action ?? null,
      resourceId ?? null,
      result ?? null,
      requester ?? null,
      scheduleId ?? null,
      description ?? null,
    ];
  };

  const { sortedData, getSortParams } = useTableSort<SystemEventResponseApi>(
    paginatedData,
    getSortableRowValues,
    0,
    'desc'
  );

  if (isLoading) return <LoadingSpinner />;
  if (totalItems === 0) return <EmptyStateNoFound />;

  return (
    <React.Fragment>
      <Table aria-label="Events table">
        <Thead>
          <Tr>
            <Th sort={getSortParams(0)}>{columnNames.date}</Th>
            <Th sort={getSortParams(1)}>{columnNames.operation}</Th>
            <Th sort={getSortParams(2)}>{columnNames.resource}</Th>
            <Th sort={getSortParams(3)}>{columnNames.status}</Th>
            <Th sort={getSortParams(4)}>{columnNames.requester}</Th>
            <Th sort={getSortParams(5)}>{columnNames.scheduledAction}</Th>
            <Th>{columnNames.description}</Th>
          </Tr>
        </Thead>
        <Tbody>
          {sortedData.map(event => (
            <Tr key={event.id}>
              <Td dataLabel={columnNames.date}>{parseScanTimestamp(event.timestamp)}</Td>
              <Td dataLabel={columnNames.operation}>{renderOperationLabel(event.action)}</Td>
              <Td dataLabel={columnNames.resource}>
                {event.resourceId ? (
                  <>
                    {renderResourceBadge(event.resourceType)}{' '}
                    <Link to={resolveResourcePath(event.resourceType ?? '-', event.resourceId)}>
                      {event.resourceName || event.resourceId}
                    </Link>
                  </>
                ) : event.action === ActionOperations.SCAN ? (
                  <>
                    <ResourceBadge label="A" color="#c9190b" /> All Accounts
                  </>
                ) : (
                  '-'
                )}
              </Td>
              <Td dataLabel={columnNames.status}>{renderActionStatusLabel(event.result)}</Td>
              <Td dataLabel={columnNames.requester}>{event.requester}</Td>
              <Td dataLabel={columnNames.scheduledAction}>
                {event.scheduleId ? <Link to="/actions/scheduler">#{event.scheduleId}</Link> : '-'}
              </Td>
              <Td dataLabel={columnNames.description}>{event.description}</Td>
            </Tr>
          ))}
        </Tbody>
      </Table>
      <TablePagination
        itemCount={totalItems}
        page={page}
        perPage={perPage}
        onSetPage={setPage}
        onPerPageSelect={setPerPage}
      />
    </React.Fragment>
  );
};

export default AuditLogsTable;
