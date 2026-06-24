import React from 'react';
import { PageSection, Panel, Content } from '@patternfly/react-core';
import AuditLogsTableToolbar from './AuditLogsTableToolbar';
import { parseAsArrayOf, parseAsString, parseAsStringEnum, useQueryStates } from 'nuqs';
import { ActionOperations, ResultStatus } from '@app/types/types.tsx';
import { ProviderApi } from '@api';
import { AuditLogsTable } from './AuditLogsTable.tsx';
import { useDocumentTitle } from '@app/utils/useDocumentTitle';

const filterParams = {
  accountName: parseAsString.withDefault(''),
  action: parseAsArrayOf(parseAsStringEnum<ActionOperations>(Object.values(ActionOperations))).withDefault([]),
  provider: parseAsArrayOf(parseAsStringEnum<ProviderApi>(Object.values(ProviderApi))).withDefault([]),
  result: parseAsArrayOf(parseAsStringEnum<ResultStatus>(Object.values(ResultStatus))).withDefault([]),
  requester: parseAsString.withDefault(''),
};

const AuditLogs: React.FunctionComponent = () => {
  useDocumentTitle('Audit Logs — ClusterIQ');
  const [{ accountName, action, provider, result, requester }, setQuery] = useQueryStates(filterParams);

  return (
    <React.Fragment>
      <PageSection hasBodyWrapper={false}>
        <Content>
          <Content component="h1">Audit logs</Content>
        </Content>
      </PageSection>
      <PageSection hasBodyWrapper={false} isFilled>
        <Panel>
          <AuditLogsTableToolbar
            searchValue={accountName}
            setSearchValue={value => setQuery({ accountName: value })}
            action={action}
            setAction={value => setQuery({ action: value || [] })}
            result={result}
            setResult={value => setQuery({ result: value })}
            requester={requester}
            setRequester={value => setQuery({ requester: value })}
            providerSelections={provider}
            setProviderSelections={value => setQuery({ provider: value || [] })}
          />
          <AuditLogsTable
            accountName={accountName}
            action={action}
            provider={provider}
            result={result}
            requester={requester}
          />
        </Panel>
      </PageSection>
    </React.Fragment>
  );
};

export default AuditLogs;
