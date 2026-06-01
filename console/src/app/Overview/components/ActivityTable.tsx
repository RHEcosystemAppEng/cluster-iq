import React from 'react';
import { Table, Thead, Tr, Th, Tbody, Td } from '@patternfly/react-table';
import { EmptyState } from '@patternfly/react-core';
import { SystemEventResponseApi } from '@api';
import { Link } from 'react-router-dom';
import { parseScanTimestamp, resolveResourcePath } from '@app/utils/parseFuncs';
import { InboxIcon } from '@patternfly/react-icons';
import {
  renderOperationLabel,
  renderActionStatusLabel,
  renderResourceBadge,
  ResourceBadge,
} from '@app/utils/renderUtils';
import { ActionOperations } from '@app/types/types';

interface ActivityTableProps {
  events: SystemEventResponseApi[];
}

export const ActivityTable: React.FunctionComponent<ActivityTableProps> = ({ events }) => {
  if (events.length === 0) {
    return <EmptyState headingLevel="h4" icon={InboxIcon} titleText="No recent events"></EmptyState>;
  }

  return (
    <Table aria-label="Recent events table" variant="compact">
      <Thead>
        <Tr>
          <Th>Date</Th>
          <Th>Operation</Th>
          <Th>Resource</Th>
          <Th>Status</Th>
          <Th>Requester</Th>
        </Tr>
      </Thead>
      <Tbody>
        {events.map(event => (
          <Tr key={event.id}>
            <Td>{parseScanTimestamp(event.timestamp)}</Td>
            <Td>{renderOperationLabel(event.action)}</Td>
            <Td>
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
            <Td>{renderActionStatusLabel(event.result)}</Td>
            <Td>{event.requester}</Td>
          </Tr>
        ))}
      </Tbody>
    </Table>
  );
};
