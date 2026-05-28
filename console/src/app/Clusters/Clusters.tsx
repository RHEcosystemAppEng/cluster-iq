import { PageSection, Panel, Content } from '@patternfly/react-core';
import React from 'react';
import ClustersTable from './components/ClustersTable';
import ClustersTableToolbar from './components/ClustersTableToolbar';
import { parseAsArrayOf, parseAsString, parseAsStringEnum, parseAsBoolean, useQueryStates } from 'nuqs';
import { ResourceStatusApi, ProviderApi } from '@api';

const filterParams = {
  status: {
    ...parseAsStringEnum<ResourceStatusApi>(Object.values(ResourceStatusApi)),
    defaultValue: null as ResourceStatusApi | null,
  },
  provider: parseAsArrayOf(parseAsStringEnum<ProviderApi>(Object.values(ProviderApi))).withDefault([]),
  clusterName: parseAsString.withDefault(''),
  accountName: parseAsString.withDefault(''),
  showTerminated: parseAsBoolean.withDefault(false),
};

const Clusters: React.FunctionComponent = () => {
  const [{ status, provider, clusterName, accountName, showTerminated }, setQuery] = useQueryStates(filterParams);

  return (
    <React.Fragment>
      <PageSection hasBodyWrapper={false}>
        <Content>
          <Content component="h1">Clusters</Content>
        </Content>
      </PageSection>
      <PageSection hasBodyWrapper={false} isFilled>
        <Panel>
          <ClustersTableToolbar
            clusterNameSearch={clusterName}
            setClusterNameSearch={value => setQuery({ clusterName: value })}
            accountNameSearch={accountName}
            setAccountNameSearch={value => setQuery({ accountName: value })}
            statusSelection={status}
            setStatusSelection={value => setQuery({ status: value })}
            providerSelections={provider}
            setProviderSelections={value => setQuery({ provider: value || [] })}
            showTerminated={showTerminated}
            setShowTerminated={value => setQuery({ showTerminated: value })}
          />
          <ClustersTable
            clusterNameSearch={clusterName}
            accountNameSearch={accountName}
            statusFilter={status}
            providerSelections={provider}
            showTerminated={showTerminated}
          />
        </Panel>
      </PageSection>
    </React.Fragment>
  );
};

export default Clusters;
