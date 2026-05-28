import React from 'react';
import { Button, Content } from '@patternfly/react-core';
import { Modal, ModalVariant } from '@patternfly/react-core/deprecated';
import { ActionOperations } from '@app/types/types';

interface ClusterActionConfirmProps {
  isOpen: boolean;
  clusterId: string;
  actionOperation: ActionOperations | null;
  onConfirm: () => void;
  onClose: () => void;
}

export const ClusterActionConfirm: React.FunctionComponent<ClusterActionConfirmProps> = ({
  isOpen,
  clusterId,
  actionOperation,
  onConfirm,
  onClose,
}) => {
  if (!isOpen || !actionOperation) {
    return null;
  }

  return (
    <Modal
      variant={ModalVariant.small}
      title="Confirm power action"
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
        Are you sure you want to <strong>{actionOperation}</strong> the cluster <strong>{clusterId}</strong>?
      </Content>
    </Modal>
  );
};
