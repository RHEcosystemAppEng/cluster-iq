import { PageSection, Panel, Content } from '@patternfly/react-core';
import React from 'react';
import NodesTableToolbar from './components/NodesTableToolbar';
import NodesTable from './components/NodesTable';
import { parseAsArrayOf, parseAsString, parseAsStringEnum, parseAsBoolean, useQueryStates } from 'nuqs';
import { ResourceStatusApi, ProviderApi } from '@api';
import { useDocumentTitle } from '@app/utils/useDocumentTitle';

const filterParams = {
  status: {
    ...parseAsStringEnum<ResourceStatusApi>(Object.values(ResourceStatusApi)),
    defaultValue: null as ResourceStatusApi | null,
  },
  provider: parseAsArrayOf(parseAsStringEnum<ProviderApi>(Object.values(ProviderApi))).withDefault([]),
  serverName: parseAsString.withDefault(''),
  showTerminated: parseAsBoolean.withDefault(false),
};

const Nodes: React.FunctionComponent = () => {
  useDocumentTitle('Nodes — ClusterIQ');
  const [{ status, provider, serverName, showTerminated }, setQuery] = useQueryStates(filterParams);

  return (
    <React.Fragment>
      <PageSection hasBodyWrapper={false}>
        <Content>
          <Content component="h1">Nodes</Content>
        </Content>
      </PageSection>
      <PageSection hasBodyWrapper={false} isFilled>
        <Panel>
          <NodesTableToolbar
            searchValue={serverName}
            setSearchValue={value => setQuery({ serverName: value })}
            statusSelection={status}
            setStatusSelection={value => setQuery({ status: value })}
            providerSelections={provider}
            setProviderSelections={value => setQuery({ provider: value || [] })}
            showTerminated={showTerminated}
            setShowTerminated={value => setQuery({ showTerminated: value })}
          />
          <NodesTable
            searchValue={serverName}
            statusSelection={status}
            providerSelections={provider}
            showTerminated={showTerminated}
          />
        </Panel>
      </PageSection>
    </React.Fragment>
  );
};

export default Nodes;
