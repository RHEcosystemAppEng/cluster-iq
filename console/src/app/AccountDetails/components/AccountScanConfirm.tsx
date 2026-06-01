import React from 'react';
import { Button, Content } from '@patternfly/react-core';
import { Modal, ModalVariant } from '@patternfly/react-core/deprecated';

interface AccountScanConfirmProps {
  isOpen: boolean;
  accountName: string;
  onConfirm: () => void;
  onClose: () => void;
}

export const AccountScanConfirm: React.FunctionComponent<AccountScanConfirmProps> = ({
  isOpen,
  accountName,
  onConfirm,
  onClose,
}) => {
  if (!isOpen) {
    return null;
  }

  return (
    <Modal
      variant={ModalVariant.small}
      title="Confirm scan"
      isOpen={isOpen}
      onClose={onClose}
      actions={[
        <Button key="confirm" variant="primary" onClick={onConfirm}>
          Confirm
        </Button>,
        <Button key="cancel" variant="link" onClick={onClose}>
          Cancel
        </Button>,
      ]}
    >
      <Content component="p">
        Are you sure you want to trigger a scan on the account <strong>{accountName}</strong>?
      </Content>
    </Modal>
  );
};
