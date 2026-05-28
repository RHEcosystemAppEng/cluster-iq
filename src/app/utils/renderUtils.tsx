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
      return <Label color="orange">Instant Action</Label>;
    case ActionTypes.SCHEDULED_ACTION:
      return <Label color="green">Scheduled Action</Label>;
    case ActionTypes.CRON_ACTION:
      return <Label color="blue">Cron Action</Label>;
    default:
      return <Label color="grey">{labelText}</Label>;
  }
}

export function renderOperationLabel(labelText: string | null | undefined) {
  switch (labelText) {
    case ActionOperations.POWER_ON:
      return <Label color="teal">{labelText}</Label>;
    case ActionOperations.POWER_OFF:
      return <Label color="purple">{labelText}</Label>;
    default:
      return <Label color="grey">{labelText}</Label>;
  }
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
