import { api, ResourceStatusApi, ActionRequestApi } from '@api';
import { Dropdown, DropdownItem, DropdownList, MenuToggle, MenuToggleElement } from '@patternfly/react-core';
import React from 'react';
import { useParams } from 'react-router-dom';
import { ActionOperations, ActionTypes, ActionStatus } from '@app/types/types';
import { ClusterActionConfirm } from './ClusterActionConfirm';
import { useUser } from '@app/Contexts/UserContext';

interface ClusterDetailsDropdownProps {
  clusterStatus: ResourceStatusApi | null;
}

export const ClusterDetailsDropdown: React.FunctionComponent<ClusterDetailsDropdownProps> = () => {
  const { userEmail } = useUser();
  const [isOpen, setIsOpen] = React.useState(false);
  const [isModalOpen, setIsModalOpen] = React.useState(false);
  const [actionOperation, setActionOperation] = React.useState<ActionOperations | null>(null);

  const { clusterID } = useParams();

  const onSelect = (_event: React.MouseEvent<Element, MouseEvent> | undefined, value: string | number | undefined) => {
    const operation = value as ActionOperations;

    if (operation === ActionOperations.POWER_ON || operation === ActionOperations.POWER_OFF) {
      setActionOperation(operation);
      setIsModalOpen(true);
    }

    setIsOpen(false);
  };

  const actionCreate = async (clusterId: string, operation: string) => {
    const actionRequest = {
      clusterId,
      enabled: true,
      operation,
      requester: userEmail || undefined,
      status: ActionStatus.Pending,
      type: ActionTypes.INSTANT_ACTION,
    } as ActionRequestApi;

    try {
      await api.actions.actionsCreate([actionRequest]);
    } catch (error) {
      console.error('Failed to create action:', error);
    }
  };

  const resetModalState = () => {
    setIsModalOpen(false);
    setActionOperation(null);
  };

  return (
    <>
      <Dropdown
        isOpen={isOpen}
        onSelect={onSelect}
        onOpenChange={setIsOpen}
        popperProps={{ position: 'end' }}
        toggle={(toggleRef: React.Ref<MenuToggleElement>) => (
          <MenuToggle ref={toggleRef} onClick={() => setIsOpen(v => !v)} isExpanded={isOpen}>
            Actions
          </MenuToggle>
        )}
      >
        <DropdownList>
          <DropdownItem value={ActionOperations.POWER_ON} key="power-on">
            {ActionOperations.POWER_ON}
          </DropdownItem>
          <DropdownItem value={ActionOperations.POWER_OFF} key="power-off">
            {ActionOperations.POWER_OFF}
          </DropdownItem>
        </DropdownList>
      </Dropdown>

      <ClusterActionConfirm
        isOpen={isModalOpen}
        onConfirm={() => {
          if (!clusterID || !actionOperation) return;
          actionCreate(clusterID, actionOperation);
          resetModalState();
        }}
        onClose={resetModalState}
        actionOperation={actionOperation}
        clusterId={clusterID!}
      />
    </>
  );
};
