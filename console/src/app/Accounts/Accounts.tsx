import { PageSection, Panel, Content } from '@patternfly/react-core';
import React from 'react';
import AccountsToolbar from './components/AccountsToolbar';
import AccountsTable from './components/AccountsTable';
import { parseAsArrayOf, parseAsString, parseAsStringEnum, useQueryStates } from 'nuqs';
import { ProviderApi } from '@api';

const filterParams = {
  provider: parseAsArrayOf(parseAsStringEnum<ProviderApi>(Object.values(ProviderApi))).withDefault([]),
  accountName: parseAsString.withDefault(''),
};

const Accounts: React.FunctionComponent = () => {
  const [{ provider, accountName }, setQuery] = useQueryStates(filterParams);

  return (
    <React.Fragment>
      <PageSection hasBodyWrapper={false}>
        <Content>
          <Content component="h1">Accounts</Content>
        </Content>
      </PageSection>
      <PageSection hasBodyWrapper={false} isFilled>
        <Panel>
          <AccountsToolbar
            searchValue={accountName}
            setSearchValue={value => setQuery({ accountName: value })}
            providerSelections={provider}
            setProviderSelections={value => setQuery({ provider: value || null })}
          />
          <AccountsTable searchValue={accountName} providerSelections={provider} />
        </Panel>
      </PageSection>
    </React.Fragment>
  );
};

export default Accounts;
