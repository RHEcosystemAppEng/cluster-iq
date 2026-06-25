import React from 'react';
import { ModalCreateAction } from './components/ModalCreateAction';
import { Flex, FlexItem, Button, PageSection, Panel, Content } from '@patternfly/react-core';
import ScheduleActionsTable from './components/ActionsTable';
import ScheduleActionsTableToolbar from './components/ActionsToolBar';
import { ActionOperations, ActionTypes, ActionStatus } from '@app/types/types';
import { parseAsArrayOf, parseAsString, parseAsStringEnum, useQueryStates } from 'nuqs';

import { parseAsBooleanNullable } from '@app/utils/parseFuncs';
import { useDocumentTitle } from '@app/utils/useDocumentTitle';

// Nullable boolean: "true" -> true, "false" -> false, missing/other -> null

const filterParams = {
  accountId: parseAsString.withDefault(''),
  action: parseAsArrayOf(parseAsStringEnum<ActionOperations>(Object.values(ActionOperations))).withDefault([]),
  type: parseAsStringEnum<ActionTypes>(Object.values(ActionTypes)),
  status: parseAsStringEnum<ActionStatus>(Object.values(ActionStatus)),
  enabled: parseAsBooleanNullable,
};

const Scheduler: React.FunctionComponent = () => {
  useDocumentTitle('Scheduler — ClusterIQ');
  const [{ accountId, action, type, status, enabled }, setQuery] = useQueryStates(filterParams);
  const [isModalOpen, setIsModalOpen] = React.useState(false);
  const [reloadFlag, setReloadFlag] = React.useState(0);

  const onClick = () => {
    setIsModalOpen(true);
  };

  const resetModalState = () => {
    setIsModalOpen(false);
  };

  return (
    <React.Fragment>
      <PageSection hasBodyWrapper={false}>
        <Flex justifyContent={{ default: 'justifyContentSpaceBetween' }} alignItems={{ default: 'alignItemsCenter' }}>
          <FlexItem>
            <Content>
              <Content component="h1">Scheduled Actions</Content>
            </Content>
          </FlexItem>
          <FlexItem>
            <Button variant="primary" onClick={onClick}>
              New Action
            </Button>
          </FlexItem>
        </Flex>
      </PageSection>
      <PageSection hasBodyWrapper={false} isFilled>
        <Panel>
          <ScheduleActionsTableToolbar
            searchValue={accountId}
            setSearchValue={value => setQuery({ accountId: value })}
            actionOperation={action}
            setOperation={value => setQuery({ action: value || [] })}
            actionType={type}
            setType={value => setQuery({ type: value })}
            actionStatus={status}
            setStatus={value => setQuery({ status: value })}
            actionEnabled={enabled}
            setEnabled={value => setQuery({ enabled: value })}
          />
          <ScheduleActionsTable
            actionType={type}
            actionOperation={action}
            actionStatus={status}
            actionEnabled={enabled}
            accountId={accountId}
            reloadFlag={reloadFlag}
          />
        </Panel>
        <ModalCreateAction isOpen={isModalOpen} onClose={resetModalState} onCreated={() => setReloadFlag(k => k + 1)} />
      </PageSection>
    </React.Fragment>
  );
};

export default Scheduler;
