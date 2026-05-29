import React, { CSSProperties } from 'react';
import { ActionTypes, ActionStatus, ActionOperations, ResultStatus } from '@app/types/types';
import { ResourceStatusApi } from '@api';
import { Label } from '@patternfly/react-core';
import {
  PendingIcon,
  OnRunningIcon,
  InfoCircleIcon,
  ExclamationTriangleIcon,
  ExclamationCircleIcon,
  UnknownIcon,
} from '@patternfly/react-icons';
import { Link } from 'react-router-dom';

const resourceBadgeStyle = (color: string): CSSProperties => ({
  display: 'inline-flex',
  alignItems: 'center',
  justifyContent: 'center',
  minWidth: '1.5em',
  height: '1.5em',
  padding: '0 0.35em',
  borderRadius: '50%',
  backgroundColor: color,
  color: '#fff',
  fontSize: '0.75rem',
  fontWeight: 700,
  lineHeight: 1,
  verticalAlign: 'middle',
});

export function ResourceBadge({ label, color }: { label: string; color: string }) {
  return <span style={resourceBadgeStyle(color)}>{label}</span>;
}

const resourceLabelStyle = (color: string): CSSProperties => ({
  display: 'inline-flex',
  alignItems: 'center',
  justifyContent: 'center',
  padding: '0.1em 0.5em',
  borderRadius: '0.75em',
  backgroundColor: color,
  color: '#fff',
  fontSize: '0.75em',
  fontWeight: 700,
  lineHeight: 1,
  verticalAlign: 'middle',
});

export function ResourceLabel({ label, color }: { label: string; color: string }) {
  return <span style={resourceLabelStyle(color)}>{label}</span>;
}

export function renderActionStatusLabel(labelText: string | null | undefined) {
  switch (labelText) {
    case ActionStatus.Running:
      return <Label color="purple">{labelText}</Label>;
    case ActionStatus.Success:
      return <Label color="green">{labelText}</Label>;
    case ActionStatus.Failed:
      return <Label color="red">{labelText}</Label>;
    case ActionStatus.Pending:
      return <Label color="yellow">{labelText}</Label>;
    default:
      return <Label color="grey">{labelText}</Label>;
  }
}

export function renderStatusLabel(labelText: string | null | undefined) {
  switch (labelText) {
    case ResourceStatusApi.Running:
      return <Label color="green">{labelText}</Label>;
    case ResourceStatusApi.Stopped:
      return <Label color="red">{labelText}</Label>;
    case ResourceStatusApi.Terminated:
      return <Label color="purple">{labelText}</Label>;
    default:
      return <Label color="grey">{labelText}</Label>;
  }
}

export function renderActionTypeLabel(labelText: string | null | undefined) {
  switch (labelText) {
    case ActionTypes.INSTANT_ACTION:
      return <Label color="orange">Instant</Label>;
    case ActionTypes.SCHEDULED_ACTION:
      return <Label color="green">Scheduled</Label>;
    case ActionTypes.CRON_ACTION:
      return <Label color="blue">Cron</Label>;
    default:
      return <Label color="grey">{labelText}</Label>;
  }
}

export function renderOperationLabel(labelText: string | null | undefined) {
  switch (labelText) {
    case ActionOperations.POWER_ON:
      return <Label color="green">{labelText}</Label>;
    case ActionOperations.POWER_OFF:
      return <Label color="red">{labelText}</Label>;
    case ActionOperations.SCAN:
      return <Label color="orange">{labelText}</Label>;
    default:
      return <Label color="grey">{labelText}</Label>;
  }
}

export function renderTargetLabel(
  clusterId: string | undefined,
  clusterName: string | undefined,
  targetAccountIds: string[] | undefined,
  targetAccountNames: string[] | undefined,
  selectAll: boolean | undefined
): React.ReactNode {
  if (clusterId) {
    return (
      <>
        <ResourceBadge label="C" color="#0066cc" />{' '}
        <Link to={`/clusters/${clusterId}`}>{clusterName || clusterId}</Link>
      </>
    );
  }
  if (!selectAll && targetAccountIds?.length) {
    const accId = targetAccountIds[0];
    const accName = targetAccountNames?.[0];
    return (
      <>
        <ResourceBadge label="A" color="#c9190b" /> <Link to={`/accounts/${accId}`}>{accName || accId}</Link>
      </>
    );
  }
  return (
    <>
      <ResourceBadge label="A" color="#c9190b" /> All Accounts
    </>
  );
}

export const getResultIcon = (result: ResultStatus) => {
  return (
    {
      [ResultStatus.Success]: (
        <InfoCircleIcon color="var(--pf-t--global--color--status--success--default)" title="Success" />
      ),
      [ResultStatus.Running]: (
        <OnRunningIcon color="var(--pf-t--global--color--status--info--default)" title="Running" />
      ),
      [ResultStatus.Pending]: (
        <PendingIcon color="var(--pf-t--global--color--status--warning--default)" title="Pending" />
      ),
      [ResultStatus.Failed]: (
        <ExclamationTriangleIcon color="var(--pf-t--global--color--status--danger--default)" title="Error" />
      ),
      [ResultStatus.Warning]: (
        <ExclamationCircleIcon color="var(--pf-t--global--color--status--warning--default)" title="Warning" />
      ),
      [ResultStatus.Unknown]: <UnknownIcon color="gray" title="Unknown" />,
    }[result] || <UnknownIcon color="gray" title="Unknown" />
  );
};
