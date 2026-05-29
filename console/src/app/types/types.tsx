export enum ResultStatus {
  Pending = 'Pending',
  Running = 'Running',
  Success = 'Success',
  Failed = 'Failed',
  Warning = 'Warning',
  Unknown = 'Unknown',
}

export enum ActionStatus {
  Running = 'Running',
  Success = 'Success',
  Failed = 'Failed',
  Pending = 'Pending',
  Unknown = 'Unknown',
}

export enum ActionOperations {
  POWER_ON = 'PowerOn',
  POWER_OFF = 'PowerOff',
}

export enum ActionTypes {
  INSTANT_ACTION = 'instant_action',
  SCHEDULED_ACTION = 'scheduled_action',
  CRON_ACTION = 'cron_action',
}

export interface BaseAction {
  type: 'instant_action' | 'scheduled_action' | 'cron_action';
  operation: 'PowerOff' | 'PowerOn';
  target: {
    clusterID: string;
  };
  status: 'Pending';
  enabled: boolean;
}

export interface ScheduledAction extends BaseAction {
  type: 'scheduled_action';
  time: string;
}

export interface CronAction extends BaseAction {
  type: 'cron_action';
  cronExp: string;
}

export interface InstantAction {
  description?: string;
}
