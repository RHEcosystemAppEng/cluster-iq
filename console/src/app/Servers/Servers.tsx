import { PageSection, Panel, Content } from '@patternfly/react-core';
import React from 'react';
import ServersTableToolbar from './components/ServersTableToolbar';
import ServersTable from './components/ServersTable';
import { parseAsArrayOf, parseAsString, parseAsStringEnum, parseAsBoolean, useQueryStates } from 'nuqs';
import { ResourceStatusApi, ProviderApi } from '@api';

const filterParams = {
  status: {
    ...parseAsStringEnum<ResourceStatusApi>(Object.values(ResourceStatusApi)),
    defaultValue: null as ResourceStatusApi | null,
  },
  provider: parseAsArrayOf(parseAsStringEnum<ProviderApi>(Object.values(ProviderApi))).withDefault([]),
  serverName: parseAsString.withDefault(''),
  showTerminated: parseAsBoolean.withDefault(false),
};

const Servers: React.FunctionComponent = () => {
  const [{ status, provider, serverName, showTerminated }, setQuery] = useQueryStates(filterParams);

  return (
    <React.Fragment>
      <PageSection hasBodyWrapper={false}>
        <Content>
          <Content component="h1">Servers</Content>
        </Content>
      </PageSection>
      <PageSection hasBodyWrapper={false} isFilled>
        <Panel>
          <ServersTableToolbar
            searchValue={serverName}
            setSearchValue={value => setQuery({ serverName: value })}
            statusSelection={status}
            setStatusSelection={value => setQuery({ status: value })}
            providerSelections={provider}
            setProviderSelections={value => setQuery({ provider: value || [] })}
            showTerminated={showTerminated}
            setShowTerminated={value => setQuery({ showTerminated: value })}
          />
          <ServersTable
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

export default Servers;
