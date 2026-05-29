import { LoadingSpinner } from '@app/components/common/LoadingSpinner';
import { ActionOperations, ResultStatus } from '@app/types/types';
import { SystemEventResponseApi } from '@api';
import { Table, Tbody, Td, Th, Thead, Tr } from '@patternfly/react-table';
import React, { useMemo } from 'react';
import { getResultIcon, renderOperationLabel } from '@app/utils/renderUtils';
import { useTableSort } from '@app/hooks/useTableSort.tsx';
import { EmptyState } from '@patternfly/react-core';
import { TablePagination } from '@app/components/common/TablesPagination';
import { SearchIcon } from '@patternfly/react-icons';
import { AuditLogsTableProps } from './types';
import { Link } from 'react-router-dom';
import { useEvents } from '@app/hooks/useEvents';
import { useTablePagination } from '@app/hooks/useTablePagination';

const columnNames = {
  action: 'Action',
  result: 'Result',
  resource: 'Resource',
  account: 'Account',
  provider: 'Provider',
  triggeredBy: 'Triggered By',
  description: 'Description',
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
  triggered_by,
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

    if (triggered_by) {
      filteredResult = filteredResult.filter(event =>
        event.triggeredBy?.toLowerCase().includes(triggered_by.toLowerCase())
      );
    }

    return filteredResult;
  }, [allEvents, accountName, action, provider, result, triggered_by]);

  const { page, perPage, setPage, setPerPage, paginatedData, totalItems } = useTablePagination({
    data: filtered,
    filterDeps: [accountName, action, provider, result, triggered_by],
  });

  const getSortableRowValues = (event: SystemEventResponseApi): (string | number | null)[] => {
    const { action, result, resourceId, accountId, provider, triggeredBy, description, timestamp } = event;
    return [
      action ?? null,
      result ?? null,
      resourceId ?? null,
      accountId ?? null,
      provider ?? null,
      triggeredBy ?? null,
      description ?? null,
      timestamp ?? null,
    ];
  };

  const { sortedData, getSortParams } = useTableSort<SystemEventResponseApi>(
    paginatedData,
    getSortableRowValues,
    7,
    'desc'
  );

  if (isLoading) return <LoadingSpinner />;
  if (totalItems === 0) return <EmptyStateNoFound />;

  return (
    <React.Fragment>
      <Table aria-label="Events table">
        <Thead>
          <Tr>
            <Th sort={getSortParams(2)}>{columnNames.resource}</Th>
            <Th sort={getSortParams(0)}>{columnNames.action}</Th>
            <Th sort={getSortParams(3)}>{columnNames.account}</Th>
            <Th sort={getSortParams(4)}>{columnNames.provider}</Th>
            <Th sort={getSortParams(5)}>{columnNames.triggeredBy}</Th>
            <Th>{columnNames.description}</Th>
            <Th sort={getSortParams(1)}>{columnNames.result}</Th>
            <Th sort={getSortParams(7)}>{columnNames.date}</Th>
          </Tr>
        </Thead>
        <Tbody>
          {sortedData.map(event => (
            <Tr key={event.id}>
              <Td dataLabel={event.resourceId}>
                <Link
                  to={
                    event.resourceType === 'instance'
                      ? `/instances/${event.resourceId}`
                      : `/clusters/${event.resourceId}`
                  }
                >
                  {event.resourceId}
                </Link>
              </Td>
              <Td>{renderOperationLabel(event.action)}</Td>
              <Td>
                <Link to={`/accounts/${event.accountId}`}>{event.accountId}</Link>
              </Td>
              <Td>{event.provider}</Td>
              <Td>{event.triggeredBy}</Td>
              <Td>{event.description}</Td>
              <Td>
                {getResultIcon(event.result as ResultStatus)} {event.result}
              </Td>
              <Td>{event.timestamp}</Td>
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
