import { api, ActionRequestApi } from '@api';
import { Dropdown, DropdownItem, DropdownList, MenuToggle, MenuToggleElement } from '@patternfly/react-core';
import React from 'react';
import { ActionOperations, ActionTypes, ActionStatus } from '@app/types/types';
import { AccountScanConfirm } from './AccountScanConfirm';
import { useUser } from '@app/Contexts/UserContext';

interface AccountDetailsDropdownProps {
  accountId: string;
  accountName: string;
}

export const AccountDetailsDropdown: React.FunctionComponent<AccountDetailsDropdownProps> = ({
  accountId,
  accountName,
}) => {
  const { userEmail } = useUser();
  const [isOpen, setIsOpen] = React.useState(false);
  const [isModalOpen, setIsModalOpen] = React.useState(false);

  const onSelect = () => {
    setIsModalOpen(true);
    setIsOpen(false);
  };

  const actionCreate = () => {
    const actionRequest = {
      accountId,
      description: `Scan ${accountName} account`,
      enabled: true,
      operation: ActionOperations.SCAN,
      requester: userEmail || undefined,
      status: ActionStatus.Pending,
      type: ActionTypes.INSTANT_ACTION,
    } as ActionRequestApi;

    api.actions.actionsCreate([actionRequest]);
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
          <DropdownItem value={ActionOperations.SCAN} key="scan">
            {ActionOperations.SCAN}
          </DropdownItem>
        </DropdownList>
      </Dropdown>

      <AccountScanConfirm
        isOpen={isModalOpen}
        onConfirm={() => {
          actionCreate();
          setIsModalOpen(false);
        }}
        onClose={() => setIsModalOpen(false)}
        accountName={accountName}
      />
    </>
  );
};
