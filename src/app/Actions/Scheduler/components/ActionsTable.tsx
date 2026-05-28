import { renderActionTypeLabel, renderOperationLabel, renderActionStatusLabel } from '@app/utils/renderUtils';
import { Table, Thead, Tr, Th, Tbody, Td } from '@patternfly/react-table';
import { Label } from '@patternfly/react-core';
import React, { useEffect, useMemo } from 'react';
import { ActionStatus, ActionOperations, ActionTypes } from '@app/types/types';
import { Link } from 'react-router-dom';
import { LoadingSpinner } from '@app/components/common/LoadingSpinner';
import { TablePagination } from '@app/components/common/TablesPagination';
import { ActionsColumn } from '@patternfly/react-table';
import { rowActions } from './ActionsKebabMenu';
import { useScheduleActions, useInvalidateScheduleActions } from '@app/hooks/useScheduleActions';
import { useTablePagination } from '@app/hooks/useTablePagination';

export const ScheduleActionsTable: React.FunctionComponent<{
  actionType: ActionTypes | null;
  actionOperation: ActionOperations[] | null;
  actionStatus: ActionStatus | null;
  actionEnabled: boolean | null;
  accountId: string | null;
  reloadFlag: number;
}> = ({ actionType, actionOperation, actionStatus, actionEnabled, accountId, reloadFlag }) => {
  const { data: allActions = [], isLoading, refetch } = useScheduleActions();
  const invalidateScheduleActions = useInvalidateScheduleActions();

  useEffect(() => {
    if (reloadFlag > 0) {
      refetch();
    }
  }, [reloadFlag, refetch]);

  const filtered = useMemo(() => {
    let result = allActions;

    if (actionType) {
      result = result.filter(item => item.type === actionType);
    }

    if (accountId) {
      result = result.filter(item => item.accountId?.includes(accountId));
    }

    if (actionOperation?.length) {
      result = result.filter(item => {
        return actionOperation.includes(item.operation as never);
      });
    }

    if (actionStatus) {
      result = result.filter(item => item.status === actionStatus);
    }

    if (actionEnabled !== null) {
      result = result.filter(item => item.enabled === actionEnabled);
    }

    return result;
  }, [allActions, actionType, accountId, actionOperation, actionStatus, actionEnabled]);

  const { page, perPage, setPage, setPerPage, paginatedData, totalItems } = useTablePagination({
    data: filtered,
    filterDeps: [actionType, actionOperation, actionStatus, actionEnabled, accountId],
  });

  const columnNames = {
    id: 'ID',
    type: 'Action Type',
    time: 'Time',
    cronExpression: 'Cron Expression',
    operation: 'Operation',
    status: 'Status',
    clusterId: 'Cluster ID',
    accountId: 'Account ID',
    region: 'Region',
    enabled: 'Enabled',
  };

  return (
    <>
      {isLoading ? (
        <LoadingSpinner />
      ) : (
        <Table aria-label="ScheduleActions table">
          <Thead>
            <Tr>
              <Th>{columnNames.id}</Th>
              <Th>{columnNames.type}</Th>
              <Th>{columnNames.time}</Th>
              <Th>{columnNames.cronExpression}</Th>
              <Th>{columnNames.operation}</Th>
              <Th>{columnNames.status}</Th>
              <Th>{columnNames.clusterId}</Th>
              <Th>{columnNames.region}</Th>
              <Th>{columnNames.accountId}</Th>
              <Th>{columnNames.enabled}</Th>
            </Tr>
          </Thead>
          <Tbody>
            {paginatedData.map(action => (
              <Tr key={action.id}>
                <Td dataLabel={columnNames.id}>{action.id}</Td>
                <Td dataLabel={columnNames.type}>{renderActionTypeLabel(action.type)}</Td>
                <Td dataLabel={columnNames.time}>{action.type !== ActionTypes.CRON_ACTION ? action.time : '-'}</Td>
                <Td dataLabel={columnNames.cronExpression}>
                  {action.type === ActionTypes.CRON_ACTION ? action.cronExpression : '-'}
                </Td>
                <Td dataLabel={columnNames.operation}>{renderOperationLabel(action.operation)}</Td>
                <Td dataLabel={columnNames.status}>{renderActionStatusLabel(action.status)}</Td>
                <Td dataLabel={columnNames.clusterId}>
                  <Link to={`/clusters/${action.clusterId}`}>{action.clusterId}</Link>
                </Td>
                <Td dataLabel={columnNames.region}>{action.region}</Td>
                <Td dataLabel={columnNames.accountId}>
                  <Link to={`/accounts/${action.accountId}`}>{action.accountId}</Link>
                </Td>
                <Td dataLabel={columnNames.enabled}>
                  {action.enabled ? <Label color="green">Enabled</Label> : <Label color="red">Disabled</Label>}
                </Td>
                <Td isActionCell aria-label="Row actions">
                  <ActionsColumn items={rowActions(action, invalidateScheduleActions)} />
                </Td>
              </Tr>
            ))}
          </Tbody>
        </Table>
      )}
      <TablePagination
        itemCount={totalItems}
        page={page}
        perPage={perPage}
        onSetPage={setPage}
        onPerPageSelect={setPerPage}
      />
    </>
  );
};

export default ScheduleActionsTable;
