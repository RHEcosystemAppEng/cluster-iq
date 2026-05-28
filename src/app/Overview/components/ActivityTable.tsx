import React from 'react';
import { Table, Thead, Tr, Th, Tbody, Td } from '@patternfly/react-table';
import { EmptyState } from '@patternfly/react-core';
import { SystemEventResponseApi } from '@api';
import { Link } from 'react-router-dom';
import { resolveResourcePath } from '@app/utils/parseFuncs';
import { InboxIcon } from '@patternfly/react-icons';
import { getResultIcon } from '@app/utils/renderUtils';
import { ResultStatus } from '@app/types/types';

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
          <Th>Time</Th>
          <Th>Action</Th>
          <Th>Result</Th>
          <Th>Resource</Th>
          <Th>Triggered By</Th>
        </Tr>
      </Thead>
      <Tbody>
        {events.map(event => (
          <Tr key={event.id}>
            <Td>{event.timestamp ? new Date(event.timestamp).toLocaleString('es-ES') : '-'}</Td>
            <Td>{event.action}</Td>
            <Td>
              {getResultIcon(event.result as ResultStatus)} {event.result}
            </Td>
            <Td>
              <Link to={resolveResourcePath(event.resourceType ?? '-', event.resourceId ?? '-')}>
                {event.resourceId}
              </Link>
            </Td>
            <Td>{event.triggeredBy}</Td>
          </Tr>
        ))}
      </Tbody>
    </Table>
  );
};
