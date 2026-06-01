import {
  Button,
  Checkbox,
  FormHelperText,
  FormGroup,
  Form,
  TextInput,
  HelperText,
  HelperTextItem,
  Radio,
  ToggleGroup,
  ToggleGroupItem,
} from '@patternfly/react-core';
import { ExclamationCircleIcon, HelpIcon } from '@patternfly/react-icons';
import { Popover } from '@patternfly/react-core';
import { Modal, ModalVariant } from '@patternfly/react-core/deprecated';
import React from 'react';
import { ActionOperations, ActionTypes } from '@app/types/types';
import DateTimePicker from './DateTimePicker';
import { AccountTypeaheadSelect, ALL_ACCOUNTS_ID } from './AccountSelector';
import { ClusterTypeaheadSelect } from './ClusterSelector';
import { ActionStatus } from '@app/types/types';
import { debug } from '@app/utils/debugLogs';
import { api, AccountResponseApi, ClusterResponseApi, ActionRequestApi } from '@api';
import { useUser } from '@app/Contexts/UserContext';
import cronValidate from 'cron-validate';

interface ModalCreateActionProps {
  isOpen: boolean;
  onClose: () => void;
  onCreated: () => void;
}

export const ModalCreateAction: React.FunctionComponent<ModalCreateActionProps> = ({ isOpen, onClose, onCreated }) => {
  const { userEmail } = useUser();
  const [actionOperation, setActionOperation] = React.useState('');
  const [selectedAccount, setSelectedAccount] = React.useState<AccountResponseApi | null>(null);
  const [selectedCluster, setSelectedCluster] = React.useState<ClusterResponseApi | null>(null);
  const [scheduledDateTime, setScheduledDateTime] = React.useState('');
  const [showSchedule, setShowSchedule] = React.useState(false);
  const [cronExpression, setCronExpression] = React.useState('');
  const [cronTouched, setCronTouched] = React.useState(false);
  const [description, setDescription] = React.useState<string>('');
  const [actionType, setActionType] = React.useState<ActionTypes>(ActionTypes.INSTANT_ACTION);

  const [allAccounts, setAllAccounts] = React.useState<AccountResponseApi[]>([]);
  const [allClusters, setAllClusters] = React.useState<ClusterResponseApi[]>([]);

  const isScan = actionOperation === ActionOperations.SCAN;

  const isValidCronExpression = (expr: string): boolean => {
    if (!expr.trim()) return false;
    const result = cronValidate(expr, { preset: 'default' });
    return result.isValid();
  };

  const isTargetValid = isScan ? !!selectedAccount : !!selectedAccount && !!selectedCluster;

  const isExecutionValid =
    actionType === ActionTypes.INSTANT_ACTION ||
    (actionType === ActionTypes.SCHEDULED_ACTION && scheduledDateTime !== '') ||
    (actionType === ActionTypes.CRON_ACTION && cronExpression !== '' && isValidCronExpression(cronExpression));

  const isFormValid = actionOperation !== '' && isTargetValid && isExecutionValid;

  const handleOperationChange = (value: string) => {
    setActionOperation(value);
    setSelectedCluster(null);
    setShowSchedule(false);
    setActionType(ActionTypes.INSTANT_ACTION);
    setScheduledDateTime('');
    setCronExpression('');
    setCronTouched(false);
  };

  React.useEffect(() => {
    if (!isOpen) return;

    const controller = new AbortController();

    const fetchAccounts = async () => {
      try {
        const { data } = await api.accounts.accountsList({ page: 1, page_size: 10000 }, { signal: controller.signal });
        if (!controller.signal.aborted) {
          setAllAccounts(data.items || []);
        }
      } catch (error) {
        if (!controller.signal.aborted) {
          console.error('Error fetching accounts:', error);
          setAllAccounts([]);
        }
      }
    };

    fetchAccounts();
    return () => controller.abort();
  }, [isOpen]);

  React.useEffect(() => {
    if (!isOpen) return;

    setSelectedCluster(null);
    setAllClusters([]);

    const accountId = selectedAccount?.accountId;
    if (!accountId || accountId === ALL_ACCOUNTS_ID) return;

    const controller = new AbortController();

    const fetchClusters = async () => {
      try {
        const { data } = await api.accounts.clustersList(accountId, { signal: controller.signal });
        if (!controller.signal.aborted) {
          setAllClusters(data.items || []);
        }
      } catch (error) {
        if (!controller.signal.aborted) {
          console.error('Error fetching clusters:', error);
          setAllClusters([]);
        }
      }
    };

    fetchClusters();
    return () => controller.abort();
  }, [isOpen, selectedAccount?.accountId]);

  React.useEffect(() => {
    if (isOpen) return;

    setActionOperation('');
    setDescription('');
    setShowSchedule(false);
    setActionType(ActionTypes.INSTANT_ACTION);
    setScheduledDateTime('');
    setCronExpression('');
    setCronTouched(false);
    setSelectedAccount(null);
    setSelectedCluster(null);
    setAllClusters([]);
  }, [isOpen]);

  const handlerConfirmActionCreation = async () => {
    if (!isScan && !selectedCluster?.clusterId) {
      console.error('ClusterID is undefined. Cannot perform action');
      return;
    }

    const actionRequest = {
      accountId: selectedAccount?.accountId === ALL_ACCOUNTS_ID ? '' : selectedAccount?.accountId,
      clusterId: isScan ? undefined : selectedCluster?.clusterId,
      description: description || undefined,
      enabled: true,
      operation: actionOperation,
      region: isScan ? undefined : selectedCluster?.region,
      requester: userEmail || undefined,
      status: ActionStatus.Pending,
      type: actionType,
    } as ActionRequestApi;

    if (actionType === ActionTypes.SCHEDULED_ACTION) {
      actionRequest.time = scheduledDateTime;
    } else if (actionType === ActionTypes.CRON_ACTION) {
      actionRequest.cronExpression = cronExpression.trim();
    }

    debug('Creating action', actionRequest);
    await api.actions.actionsCreate([actionRequest]);
    onCreated();
    onClose();
  };

  if (!isOpen) {
    return null;
  }

  return (
    <Modal
      variant={ModalVariant.small}
      title="Create Action"
      isOpen={isOpen}
      onClose={onClose}
      actions={[
        <Button key="confirm" variant="primary" onClick={handlerConfirmActionCreation} isDisabled={!isFormValid}>
          Confirm
        </Button>,
        <Button key="cancel" variant="link" onClick={onClose}>
          Cancel
        </Button>,
      ]}
      appendTo={document.body}
    >
      <Form>
        {/* Operation selection */}
        <FormGroup
          label="Operation"
          isRequired
          fieldId="action-operation"
          labelHelp={
            <Popover
              headerContent="Operation"
              bodyContent="Scan discovers cloud resources across your accounts. Power On/Off starts or stops all instances in a cluster."
            >
              <Button variant="plain" aria-label="Operation help" icon={<HelpIcon />} />
            </Popover>
          }
        >
          <ToggleGroup aria-label="Select operation">
            <ToggleGroupItem
              text="Scan"
              buttonId="action-scan"
              isSelected={actionOperation === ActionOperations.SCAN}
              onChange={(_e, selected) => selected && handleOperationChange(ActionOperations.SCAN)}
            />
            <ToggleGroupItem
              text="Power On Cluster"
              buttonId="action-power-on"
              isSelected={actionOperation === ActionOperations.POWER_ON}
              onChange={(_e, selected) => selected && handleOperationChange(ActionOperations.POWER_ON)}
            />
            <ToggleGroupItem
              text="Power Off Cluster"
              buttonId="action-power-off"
              isSelected={actionOperation === ActionOperations.POWER_OFF}
              onChange={(_e, selected) => selected && handleOperationChange(ActionOperations.POWER_OFF)}
            />
          </ToggleGroup>
        </FormGroup>

        {/* Account selection */}
        <AccountTypeaheadSelect
          accounts={allAccounts}
          selectedAccount={selectedAccount}
          onSelectAccount={account => setSelectedAccount(account)}
          onClearAccount={() => setSelectedAccount(null)}
          showAllOption={isScan}
        />

        {/* Cluster selection (hidden for Scan operations) */}
        {!isScan && (
          <ClusterTypeaheadSelect
            accountId={selectedAccount?.accountId ?? null}
            clusters={allClusters}
            selectedCluster={selectedCluster}
            onSelectCluster={cluster => setSelectedCluster(cluster)}
            isDisabled={!selectedAccount}
            onClearCluster={() => setSelectedCluster(null)}
          />
        )}

        {/* Execution time */}
        <FormGroup
          label="Execution"
          fieldId="scheduled-action-type"
          labelHelp={
            <Popover
              headerContent="Execution"
              bodyContent="By default the action runs immediately. Check 'Schedule action' to run it at a specific time or on a recurring cron schedule."
            >
              <Button variant="plain" aria-label="Execution help" icon={<HelpIcon />} />
            </Popover>
          }
        >
          <Checkbox
            id="scheduled-action-type"
            name="scheduled-action-type"
            label="Schedule action"
            isChecked={showSchedule}
            onChange={(_event, checked) => {
              setShowSchedule(checked);
              setActionType(checked ? ActionTypes.SCHEDULED_ACTION : ActionTypes.INSTANT_ACTION);
            }}
          />
        </FormGroup>
        {showSchedule && (
          <>
            <FormGroup role="radiogroup" fieldId="schedule-type">
              <Radio
                id="specific-time"
                name="scheduleType"
                label="Specific time"
                isChecked={actionType === ActionTypes.SCHEDULED_ACTION}
                onChange={() => setActionType(ActionTypes.SCHEDULED_ACTION)}
              />
            </FormGroup>
            {actionType === ActionTypes.SCHEDULED_ACTION && (
              <FormGroup fieldId="datetime-picker" isRequired>
                <DateTimePicker onChange={setScheduledDateTime} />
              </FormGroup>
            )}
            <FormGroup fieldId="cron-expression-input" isRequired>
              <Radio
                id="cron-expression"
                name="scheduleType"
                label="Cron expression"
                isChecked={actionType === ActionTypes.CRON_ACTION}
                onChange={() => setActionType(ActionTypes.CRON_ACTION)}
              />
            </FormGroup>
            {actionType === ActionTypes.CRON_ACTION && (
              <FormGroup fieldId="cron-expression-input" isRequired>
                <TextInput
                  type="text"
                  id="cron-expression-input"
                  value={cronExpression}
                  onChange={(_event, value) => setCronExpression(value)}
                  onBlur={() => setCronTouched(true)}
                  validated={cronTouched && !isValidCronExpression(cronExpression) ? 'error' : 'default'}
                  placeholder="0 0 * * *"
                />
                <FormHelperText>
                  <HelperText>
                    <HelperTextItem
                      variant={cronTouched && !isValidCronExpression(cronExpression) ? 'error' : 'default'}
                      icon={
                        cronTouched && !isValidCronExpression(cronExpression) ? <ExclamationCircleIcon /> : undefined
                      }
                    >
                      {cronTouched && !isValidCronExpression(cronExpression)
                        ? 'Invalid cron expression'
                        : "Format: minute hour day-of-month month day-of-week (e.g., '0 0 * * *' for daily at midnight)"}
                    </HelperTextItem>
                  </HelperText>
                </FormHelperText>
              </FormGroup>
            )}
          </>
        )}

        {/* Description */}
        <FormGroup label="Description" fieldId="description-input">
          <TextInput
            type="text"
            id="description-input"
            value={description}
            onChange={(_event, value) => setDescription(value)}
            placeholder="Enter reason"
            aria-label="Reason for action"
          />
        </FormGroup>
      </Form>
    </Modal>
  );
};
